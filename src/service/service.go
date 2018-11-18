package service

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/andrecronje/evm/src/common"
	"github.com/andrecronje/evm/src/config"
	"github.com/andrecronje/evm/src/state"
)

var defaultGas = big.NewInt(90000)

type infoCallback func() (map[string]string, error)

type Service struct {
	sync.Mutex
	state       *state.State
	submitCh    chan []byte
	genesisFile string
	keystoreDir string
	apiAddr     string
	keyStore    *keystore.KeyStore
	pwdFile     string
	logger      *logrus.Logger

	rpcConfig *config.RpcConfig

	rpcAPIs       []rpc.API   // List of APIs currently provided by the node
	inprocHandler *rpc.Server // In-process RPC request handler to process the API requests

	ipcEndpoint string       // IPC endpoint to listen at (empty = IPC disabled)
	ipcListener net.Listener // IPC RPC listener socket to serve API requests
	ipcHandler  *rpc.Server  // IPC RPC request handler to process the API requests

	httpEndpoint  string       // HTTP endpoint (interface + port) to listen at (empty = HTTP disabled)
	httpWhitelist []string     // HTTP RPC modules to allow through this endpoint
	httpListener  net.Listener // HTTP RPC listener socket to server API requests
	httpHandler   *rpc.Server  // HTTP RPC request handler to process the API requests

	wsEndpoint string       // Websocket endpoint (interface + port) to listen at (empty = websocket disabled)
	wsListener net.Listener // Websocket RPC listener socket to server API requests
	wsHandler  *rpc.Server  // Websocket RPC request handler to process the API requests

	stop chan struct{} // Channel to wait for termination notifications
	lock sync.RWMutex

	//XXX
	getInfo infoCallback
}

func NewService(genesisFile, keystoreDir, apiAddr, pwdFile string,
	state *state.State,
	submitCh chan []byte,
	logger *logrus.Logger) *Service {
	rpcConfig := &config.DefaultRpcConfig

	return &Service{
		genesisFile: genesisFile,
		keystoreDir: keystoreDir,
		apiAddr:     apiAddr,
		pwdFile:     pwdFile,
		state:       state,
		submitCh:    submitCh,
		logger:      logger,
		// TODO: not default rpcConfig required
		rpcConfig:    rpcConfig,
		ipcEndpoint:  rpcConfig.IPCEndpoint(),
		httpEndpoint: rpcConfig.HTTPEndpoint(),
		wsEndpoint:   rpcConfig.WSEndpoint(),
	}
}

func (m *Service) Run() {
	m.checkErr(m.makeKeyStore())

	m.checkErr(m.unlockAccounts())

	m.checkErr(m.createGenesisAccounts())

	m.logger.Info("serving api...")

	err := m.StartRPC()
	if err != nil {
		panic(err)
	}

	m.serveAPI()
}

//XXX
func (m *Service) GetSubmitCh() chan []byte {
	return m.submitCh
}

//XXX
func (m *Service) SetInfoCallback(f infoCallback) {
	m.getInfo = f
}

func (m *Service) makeKeyStore() error {

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	if err := os.MkdirAll(m.keystoreDir, 0700); err != nil {
		return err
	}

	m.keyStore = keystore.NewKeyStore(m.keystoreDir, scryptN, scryptP)

	return nil
}

func (m *Service) unlockAccounts() error {

	if len(m.keyStore.Accounts()) == 0 {
		return nil
	}

	pwd, err := m.readPwd()
	if err != nil {
		m.logger.WithError(err).Error("Reading PwdFile")
		return err
	}

	for _, ac := range m.keyStore.Accounts() {
		if err := m.keyStore.Unlock(ac, string(pwd)); err != nil {
			return err
		}
		m.logger.WithField("address", ac.Address.Hex()).Debug("Unlocked account")
	}
	return nil
}

func (m *Service) createGenesisAccounts() error {
	if _, err := os.Stat(m.genesisFile); os.IsNotExist(err) {
		return nil
	}

	contents, err := ioutil.ReadFile(m.genesisFile)
	if err != nil {
		return err
	}

	var genesis struct {
		Alloc common.AccountMap
	}

	if err := json.Unmarshal(contents, &genesis); err != nil {
		return err
	}

	if err := m.state.CreateAccounts(genesis.Alloc); err != nil {
		return err
	}
	return nil
}

func (m *Service) serveAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/account/{address}", m.makeHandler(accountHandler)).Methods("GET")
	r.HandleFunc("/accounts", m.makeHandler(accountsHandler)).Methods("GET")
	r.HandleFunc("/block/{hash}", m.makeHandler(blockByHashHandler)).Methods("GET")
	r.HandleFunc("/blockById/{id}", m.makeHandler(blockByIdHandler)).Methods("GET")
	//r.HandleFunc("/blockIndex", m.makeHandler(blockIndexHandler)).Methods("GET")
	r.HandleFunc("/call", m.makeHandler(callHandler)).Methods("POST")
	r.HandleFunc("/tx", m.makeHandler(transactionHandler)).Methods("POST")
	r.HandleFunc("/transactions", m.makeHandler(transactionHandler)).Methods("POST")
	r.HandleFunc("/rawtx", m.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/sendRawTransaction", m.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/tx/{tx_hash}", m.makeHandler(txReceiptHandler)).Methods("GET")
	r.HandleFunc("/transaction/{tx_hash}", m.makeHandler(transactionReceiptHandler)).Methods("GET")
	r.HandleFunc("/info", m.makeHandler(infoHandler)).Methods("GET")
	r.HandleFunc("/html/info", m.makeHandler(htmlInfoHandler)).Methods("GET")
	http.Handle("/", &CORSServer{r})
	http.ListenAndServe(m.apiAddr, nil)
}

type CORSServer struct {
	r *mux.Router
}

func (s *CORSServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func (m *Service) makeHandler(fn func(http.ResponseWriter, *http.Request, *Service)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		fn(w, r, m)
		m.Unlock()
	}
}

func (m *Service) checkErr(err error) {
	if err != nil {
		m.logger.WithError(err).Error("ERROR")
		os.Exit(1)
	}
}

func (m *Service) readPwd() (pwd string, err error) {
	text, err := ioutil.ReadFile(m.pwdFile)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(text), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines[0], nil
}
