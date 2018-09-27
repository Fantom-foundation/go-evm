package proxy

import (
	"io/ioutil"
	"math/big"
	"time"

	"github.com/andrecronje/evm/service"
	proxy "github.com/andrecronje/lachesis/proxy/lachesis"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log = logrus.New()

//------------------------------------------------------------------------------

type Config struct {
	proxyAddr    string //bind address of this app proxy
	lachesisAddr string //address of  node
	apiAddr      string //address of HTTP API service
	ethDir       string //directory containing eth config
	pwdFile      string //file containing password to unlock ethereum accounts
	databaseFile string //file containing LevelDB database
	cache        int    //Megabytes of memory allocated to internal caching (min 16MB / database forced)
	timeout      time.Duration
	configFile   string
	StateConfig  StateConfig `yaml:"state"`
}

type StateConfig struct {
	//ChainIDs for state config
	ChainIDs []*big.Int `yaml:"chainIDs"`
}

func NewConfig(proxyAddr,
lachesisAddr,
apiAddr,
ethDir,
pwdFile,
dbFile string,
	cache int,
	timeout time.Duration) Config {

	return Config{
		proxyAddr:    proxyAddr,
		lachesisAddr: lachesisAddr,
		apiAddr:      apiAddr,
		ethDir:       ethDir,
		pwdFile:      pwdFile,
		databaseFile: dbFile,
		cache:        cache,
		timeout:      timeout,
	}
}

func (c *Config) Load() error {
	if len(c.configFile) < 1 {
		return nil
	}

	b, err := ioutil.ReadFile(c.configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, c)
}

//------------------------------------------------------------------------------

type Proxy struct {
	service       *service.Service
	lachesisProxy *proxy.SocketLachesisProxy
	submitCh      chan []byte
	logger        *logrus.Logger
}

func NewProxy(config Config, logger *logrus.Logger) (*Proxy, error) {
	err := config.Load()
	if err != nil {
		return nil, err
	}

	submitCh := make(chan []byte)

	logger.Debug("service.NewService")
	service_ := service.NewService(config.ethDir,
		config.apiAddr,
		config.pwdFile,
		submitCh,
		logger)

	logger.Debug("state.NewState")
	service_.NewStates(config.StateConfig.ChainIDs, config.databaseFile, config.cache)

	logger.Debug("proxy.NewSocketLachesisProxy")
	lachesisProxy, err := proxy.NewSocketLachesisProxy(config.lachesisAddr,
		config.proxyAddr,
		config.timeout,
		logger)
	if err != nil {
		log.WithError(err).Error("error building socket proxy")
		return nil, err
	}

	logger.Debug("Return &Proxy")
	return &Proxy{
		service:       service_,
		lachesisProxy: lachesisProxy,
		submitCh:      submitCh,
		logger:        logger,
	}, nil
}

func (p *Proxy) Run() error {

	go p.service.Run()

	p.Serve()

	return nil
}

func (p *Proxy) Serve() {
	for {
		select {
		case tx := <-p.submitCh:
			p.logger.Debug("Proxy about to submit tx")
			if err := p.lachesisProxy.SubmitTx(tx); err != nil {
				p.logger.WithError(err).Error("SubmitTx")
			}
			p.logger.Debug("Proxy submitted tx")
		case commit := <-p.lachesisProxy.CommitCh():
			p.logger.Debug("CommitBlock")
			stateHash, err := p.service.ProcessBlock(commit.Block)
			commit.Respond(stateHash.Bytes(), err)
		}
	}
}
