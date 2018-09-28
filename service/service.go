package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/andrecronje/evm/common"
	"github.com/andrecronje/evm/state"
	"github.com/andrecronje/lachesis/poset"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var defaultGas = big.NewInt(90000)

type Service struct {
	sync.Mutex
	state        *state.State
	submitCh     chan []byte
	statesFile   string
	keystoreDir  string
	apiAddr      string
	keyStore     *keystore.KeyStore
	pwdFile      string
	logger       *logrus.Logger
	dbFile       string
	dbCache      int
	defaultState *state.State
	states       map[string]*state.State
	chainIDs     []*big.Int
}

func NewService(statesFile, keystoreDir, apiAddr, pwdFile string,
	dbFile string, dbCache int,
	submitCh chan []byte,
	logger *logrus.Logger) *Service {
	return &Service{
		statesFile:  statesFile,
		keystoreDir: keystoreDir,
		apiAddr:     apiAddr,
		pwdFile:     pwdFile,
		dbFile:      dbFile,
		dbCache:     dbCache,
		submitCh:    submitCh,
		logger:      logger}
}

func (m *Service) Run() {
	m.checkErr(m.makeKeyStore())

	m.checkErr(m.unlockAccounts())

	m.checkErr(m.createGenesisAccounts())

	m.logger.Info("serving api...")
	m.serveAPI()
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
	if _, err := os.Stat(m.statesFile); os.IsNotExist(err) {
		return nil
	}

	contents, err := ioutil.ReadFile(m.statesFile)
	if err != nil {
		return err
	}

	c := &States{}
	err = yaml.Unmarshal(contents, c)
	if err != nil {
		return err
	}

	m.states = make(map[string]*state.State)

	for _, info := range c.StateConfigs {
		s, err := state.NewStateWithChainID(info.ChainID, m.logger, m.dbFile, m.dbCache)
		if err != nil {
			return err
		}
		m.states[info.ChainID.String()] = s

		if m.defaultState == nil {
			m.defaultState = s
		}

		if len(info.GenesisFile) < 1 {
			continue
		}

		if _, err := os.Stat(info.GenesisFile); os.IsNotExist(err) {
			return err
		}

		contents, err := ioutil.ReadFile(info.GenesisFile)
		if err != nil {
			return err
		}

		var genesis struct {
			Alloc common.AccountMap
		}

		if err := json.Unmarshal(contents, &genesis); err != nil {
			return err
		}

		if err := s.CreateAccounts(genesis.Alloc); err != nil {
			return err
		}
	}

	return nil
}

func (m *Service) serveAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/account/{address}", m.makeHandler(accountHandler)).Methods("GET")
	r.HandleFunc("/accounts", m.makeHandler(accountsHandler)).Methods("GET")
	r.HandleFunc("/call", m.makeHandler(callHandler)).Methods("POST")
	r.HandleFunc("/tx", m.makeHandler(transactionHandler)).Methods("POST")
	r.HandleFunc("/transactions", m.makeHandler(transactionHandler)).Methods("POST")
	r.HandleFunc("/rawtx", m.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/sendRawTransaction", m.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/tx/{tx_hash}", m.makeHandler(txReceiptHandler)).Methods("GET")
	r.HandleFunc("/transaction/{tx_hash}", m.makeHandler(transactionReceiptHandler)).Methods("GET")
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

func (m *Service) GetBalance(addr ethcommon.Address) map[string]*big.Int {
	result := make(map[string]*big.Int)
	for key, value := range m.states {
		result[key] = value.GetBalance(addr)
	}
	return result
}

func (m *Service) GetNonce(addr ethcommon.Address) map[string]uint64 {
	result := make(map[string]uint64)
	for key, value := range m.states {
		result[key] = value.GetNonce(addr)
	}
	return result
}

func (m *Service) GetState(id string) *state.State {
	s, ok := m.states[id]
	if !ok {
		return m.defaultState
	}
	return s
}

func (m *Service) ProcessBlock(block poset.Block) (hs []byte, err error) {
	m.logger.Debug("Process Block")

	blockHashBytes, _ := block.Hash()
	blockHash := ethcommon.BytesToHash(blockHashBytes)

	var fifo []*state.State
	lazyCommit := make(map[*state.State]*BlockProcessResult)
	defer func() {
		for s, r := range lazyCommit {
			s.GetCommitMutex().Unlock()
			if r.Err != nil {
				continue
			}
			r.Hash, r.Err = s.Commit()
		}
		hs = make([]byte, len(fifo)*ethcommon.HashLength)
		var errStr bytes.Buffer
		hasErr := false
		errStr.WriteString("Process Block Error:\n")
		for i, s := range fifo {
			result := lazyCommit[s]
			if result.Err != nil {
				hasErr = true
			} else {
				copy(hs[i*ethcommon.HashLength:], result.Hash[:])
			}
			errStr.WriteString(fmt.Sprintf("chain:%s err:%v\n", s.GetChainID().String(), result.Err))
		}
		if hasErr {
			err = errors.New(errStr.String())
		}
	}()
	for txIndex, txBytes := range block.Transactions() {
		tx := &types.Transaction{}
		tx.UnmarshalJSON(txBytes)
		s, ok := m.states[tx.ChainId().String()]
		if !ok {
			m.logger.WithField("ChainID", tx.ChainId().String()).Debug("state not exists")
			continue
		}

		_, ok = lazyCommit[s]
		if !ok {
			fifo = append(fifo, s)
			lazyCommit[s] = &BlockProcessResult{}
			s.GetCommitMutex().Lock()
		}

		if lazyCommit[s].Err != nil {
			continue
		}

		if err = s.ApplyTransaction(txBytes, txIndex, blockHash); err != nil {
			lazyCommit[s].Err = err
		}
	}

	return
}
