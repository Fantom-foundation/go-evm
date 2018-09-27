package service

import (
	"errors"
	"github.com/andrecronje/lachesis/hashgraph"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/andrecronje/evm/state"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var defaultGas = big.NewInt(90000)

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
	defaultState *state.State
	states       map[*big.Int]*state.State
	chainIDs     []*big.Int
}

func NewService(genesisFile, keystoreDir, apiAddr, pwdFile string,
	state *state.State,
	submitCh chan []byte,
	logger *logrus.Logger) *Service {
	return &Service{
		genesisFile: genesisFile,
		keystoreDir: keystoreDir,
		apiAddr:     apiAddr,
		pwdFile:     pwdFile,
		state:       state,
		submitCh:    submitCh,
		logger:      logger}
}

func (m *Service) NewStates(chainIDs []*big.Int, dbFile string, dbCache int) error {
	if len(chainIDs) < 1 {
		return nil
	}
	ds, err := state.NewStateWithChainID(chainIDs[0], m.logger, dbFile, dbCache)
	if err != nil {
		return err
	}
	m.defaultState = ds
	m.states[chainIDs[0]] = ds

	m.states = make(map[*big.Int]*state.State)
	for _, id := range chainIDs[1:] {
		s, err := state.NewStateWithChainID(chainIDs[0], m.logger, dbFile, dbCache)
		if err != nil {
			return err
		}
		m.states[id] = s
	}
	return nil
}

func (m *Service) Run() {
	m.checkErr(m.makeKeyStore())

	m.checkErr(m.unlockAccounts())

	m.checkErr(m.createGenesisAccounts())

	m.sortChainIDs()

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
	genesisFilesDir := filepath.Join(m.genesisFile)
	if err := os.MkdirAll(genesisFilesDir, 0700); err != nil {
		return err
	}

	fileInfos, err := ioutil.ReadDir(genesisFilesDir)
	if err != nil {
		return err
	}

	for _, info := range fileInfos {
		if info.IsDir() {
			continue
		}

		fileName := filepath.Join(genesisFilesDir, info.Name())
		contents, err := ioutil.ReadFile(fileName)
		if err != nil {
			return err
		}

		genesis := &core.Genesis{}
		err = genesis.UnmarshalJSON(contents)
		if err != nil {
			return err
		}

		config := genesis.Config
		if config == nil {
			return errors.New("genesis.Config == nil")
		}
		chainID := config.ChainID
		if chainID == nil {
			return errors.New("genesis.Config.ChainID == nil")
		}
		s, ok := m.states[chainID]
		if !ok {
			s, err = state.CopyStateWithChainID(m.defaultState, chainID)
			if err != nil {
				return err
			}
			m.states[chainID] = s
		}
		if err := s.CreateAccounts(genesis.Alloc); err != nil {
			return err
		}
	}

	return nil
}

func (m *Service) sortChainIDs() {
	var ids []*big.Int
	for id := range m.states {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Cmp(ids[j]) < 0
	})
	m.chainIDs = ids
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

func (m *Service) GetBalance(addr common.Address) map[*big.Int]*big.Int {
	result := make(map[*big.Int]*big.Int)
	for key, value := range m.states {
		result[key] = value.GetBalance(addr)
	}
	return result
}

func (m *Service) GetNonce(addr common.Address) map[*big.Int]uint64 {
	result := make(map[*big.Int]uint64)
	for key, value := range m.states {
		result[key] = value.GetNonce(addr)
	}
	return result
}

func (m *Service) GetState(id string) *state.State {
	var i big.Int
	i.SetString(id, 10)
	s, ok := m.states[&i]
	if !ok {
		return m.defaultState
	}
	return s
}

func (m *Service) ProcessBlock(block hashgraph.Block) (hs common.Hash, err error) {
	m.logger.Debug("Process Block")

	blockHashBytes, _ := block.Hash()
	blockHash := common.BytesToHash(blockHashBytes)

	lazyCommit := make(map[*state.State]*BlockProcessResult)
	defer func() {
		for s, r := range lazyCommit {
			if r.Err != nil {
				continue
			}
			r.Hash, r.Err = s.Commit()
			s.GetCommitMutex().Unlock()
		}
	}()
	for txIndex, txBytes := range block.Transactions() {
		tx := &types.Transaction{}
		tx.UnmarshalJSON(txBytes)
		s, ok := m.states[tx.ChainId()]
		if !ok {
			m.logger.WithField("ChainID", tx.ChainId().String()).Debug("state not exists")
			continue
		}

		_, ok = lazyCommit[s]
		if !ok {
			lazyCommit[s] = &BlockProcessResult{}
			s.GetCommitMutex().Lock()
		}

		if err = s.ApplyTransaction(txBytes, txIndex, blockHash); err != nil {
			lazyCommit[s].Err = err
		}
	}

	return
}
