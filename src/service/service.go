package service

import (
	"encoding/json"
	"github.com/Fantom-foundation/go-evm/src/config"
	"github.com/Fantom-foundation/go-lachesis/src/common/hexutil"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Fantom-foundation/go-evm/src/common"
	"github.com/Fantom-foundation/go-evm/src/state"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var defaultGas = hexutil.Uint64(90000)

type infoCallback func() (map[string]string, error)

type Service struct {
	sync.Mutex
	chainConfig *params.ChainConfig
	state       *state.State
	submitCh    chan []byte
	genesisFile string
	keystoreDir string
	apiAddr     string
	keyStore    *keystore.KeyStore
	am          *accounts.Manager
	pwdFile     string
	logger      *logrus.Logger

	rpcConfig *node.Config
	rpcServer *RpcServer

	//XXX
	getInfo infoCallback
}

func NewService(genesisFile, keystoreDir, apiAddr, pwdFile string,
	state *state.State,
	submitCh chan []byte,
	logger *logrus.Logger) *Service {
	// TODO: replace DefaultRpcConfig with custom
	var rpcConfig = config.DefaultRpcConfig
	// TODO: replace ChainConfig with custom
	chainConfig := &params.ChainConfig{
		ChainID: big.NewInt(666),
	}

	s := &Service{
		chainConfig: chainConfig,
		genesisFile: genesisFile,
		keystoreDir: keystoreDir,
		apiAddr:     apiAddr,
		pwdFile:     pwdFile,
		state:       state,
		submitCh:    submitCh,
		logger:      logger,
		// TODO: no-default rpcConfig required
		rpcConfig: rpcConfig,
	}
	var err error
	s.rpcServer, err = NewRpcServer(rpcConfig, s)
	if err != nil {
		panic(err)
	}
	err = s.rpcServer.Register(NewWeb3AccountServiceConstructor(s))
	if err != nil {
		panic(err)
	}
	err = s.rpcServer.Register(NewWeb3ChainServiceConstructor(s))
	if err != nil {
		panic(err)
	}

	return s
}

func (service *Service) Run() {
	service.checkErr(service.makeKeyStore())

	service.checkErr(service.unlockAccounts())

	service.checkErr(service.createGenesisAccounts())

	service.logger.Info("serving web3-api ...")
	if err := service.rpcServer.Start(); err != nil {
		panic(err)
	}
	defer func() {
		if err := service.rpcServer.Stop(); err != nil {
			panic(err)
		}
	}()

	service.serveAPI()
}

//XXX
func (service *Service) GetSubmitCh() chan []byte {
	return service.submitCh
}

//XXX
func (service *Service) SetInfoCallback(f infoCallback) {
	service.getInfo = f
}

func (service *Service) makeKeyStore() error {

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	if err := os.MkdirAll(service.keystoreDir, 0700); err != nil {
		return err
	}

	service.keyStore = keystore.NewKeyStore(service.keystoreDir, scryptN, scryptP)

	service.am = accounts.NewManager(service.keyStore)

	return nil
}

func (service *Service) unlockAccounts() error {

	if len(service.keyStore.Accounts()) == 0 {
		return nil
	}

	pwd, err := service.readPwd()
	if err != nil {
		service.logger.WithError(err).Error("Reading PwdFile")
		return err
	}

	for _, ac := range service.keyStore.Accounts() {
		if err := service.keyStore.Unlock(ac, string(pwd)); err != nil {
			return err
		}
		service.logger.WithField("address", ac.Address.Hex()).Debug("Unlocked account")
	}
	return nil
}

func (service *Service) createGenesisAccounts() error {
	if _, err := os.Stat(service.genesisFile); os.IsNotExist(err) {
		return nil
	}

	contents, err := ioutil.ReadFile(service.genesisFile)
	if err != nil {
		return err
	}

	var genesis struct {
		Alloc common.AccountMap
	}

	if err := json.Unmarshal(contents, &genesis); err != nil {
		return err
	}

	if err := service.state.CreateAccounts(genesis.Alloc); err != nil {
		return err
	}
	return nil
}

func (service *Service) serveAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/account/{address}", service.makeHandler(accountHandler)).Methods("GET")
	r.HandleFunc("/accounts", service.makeHandler(accountsHandler)).Methods("GET")
	r.HandleFunc("/block/{hash}", service.makeHandler(blockByHashHandler)).Methods("GET")
	r.HandleFunc("/blockById/{id}", service.makeHandler(blockByIdHandler)).Methods("GET")
	//r.HandleFunc("/blockIndex", service.makeHandler(blockIndexHandler)).Methods("GET")
	r.HandleFunc("/call", service.makeHandler(callHandler)).Methods("POST")
	r.HandleFunc("/tx", service.makeHandler(transactionHandler)).Methods("POST")
	r.HandleFunc("/transactions", service.makeHandler(transactionHandler)).Methods("POST")
	r.HandleFunc("/rawtx", service.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/sendRawTransaction", service.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/tx/{tx_hash}", service.makeHandler(transactionReceiptHandler)).Methods("GET")
	r.HandleFunc("/transaction/{tx_hash}", service.makeHandler(transactionReceiptHandler)).Methods("GET")
	r.HandleFunc("/info", service.makeHandler(infoHandler)).Methods("GET")
	r.HandleFunc("/html/info", service.makeHandler(htmlInfoHandler)).Methods("GET")
	http.Handle("/", &CORSServer{r})
	if err := http.ListenAndServe(service.apiAddr, nil); err != nil {
		panic(err)
	}
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

func (service *Service) makeHandler(fn func(http.ResponseWriter, *http.Request, *Service)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.Lock()
		fn(w, r, service)
		service.Unlock()
	}
}

func (service *Service) checkErr(err error) {
	if err != nil {
		service.logger.WithError(err).Error("ERROR")
		os.Exit(1)
	}
}

func (service *Service) readPwd() (pwd string, err error) {
	text, err := ioutil.ReadFile(service.pwdFile)
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

func (service *Service) AccountManager() *accounts.Manager {
	return service.am
}

func (service *Service) ChainConfig() *params.ChainConfig {
	return service.chainConfig
}
