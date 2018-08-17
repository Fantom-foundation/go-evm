package main

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"time"

	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/andrecronje/evm/proxy"
	"github.com/andrecronje/lachesis/version"
)

var (
	DatadirFlag = cli.StringFlag{
		Name:  "datadir",
		Usage: "Directory for the databases and keystore",
		Value: defaultDataDir(),
	}
	LachesisAddressFlag = cli.StringFlag{
		Name:  "lachesis_addr",
		Usage: "IP:Port of Lachesis node",
		Value: "127.0.0.1:1338",
	}
	ProxyAddressFlag = cli.StringFlag{
		Name:  "proxy_addr",
		Usage: "IP:Port to bind Proxy server",
		Value: "127.0.0.1:1339",
	}
	APIAddrFlag = cli.StringFlag{
		Name:  "api_addr",
		Usage: "IP:Port to bind API server",
		Value: ":8080",
	}
	LogLevelFlag = cli.StringFlag{
		Name:  "log_level",
		Usage: "debug, info, warn, error, fatal, panic",
		Value: "debug",
	}
	PwdFlag = cli.StringFlag{
		Name:  "pwd",
		Usage: "Password file to unlock accounts",
		Value: fmt.Sprintf("%s/pwd.txt", defaultDataDir()),
	}
	DatabaseFlag = cli.StringFlag{
		Name:  "db",
		Usage: "Database file",
		Value: fmt.Sprintf("%s/chaindata", defaultDataDir()),
	}
	CacheFlag = cli.IntFlag{
		Name:  "cache",
		Usage: "Megabytes of memory allocated to internal caching (min 16MB / database forced)",
		Value: 128,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "evm"
	app.Usage = "Lightweight EVM"
	app.HideVersion = true
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Action: run,
			Flags: []cli.Flag{
				DatadirFlag,
				LachesisAddressFlag,
				ProxyAddressFlag,
				APIAddrFlag,
				LogLevelFlag,
				PwdFlag,
				DatabaseFlag,
				CacheFlag,
			},
		},
		{
			Name:   "version",
			Usage:  "Show version info",
			Action: printVersion,
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	logger := logrus.New()
	logger.Level = logLevel(c.String(LogLevelFlag.Name))

	datadir := c.String(DatadirFlag.Name)
	lachesisAddress := c.String(LachesisAddressFlag.Name)
	proxyAddress := c.String(ProxyAddressFlag.Name)
	apiAddress := c.String(APIAddrFlag.Name)
	pwdFile := c.String(PwdFlag.Name)
	databaseFile := c.String(DatabaseFlag.Name)
	dbCache := c.Int(CacheFlag.Name)

	logger.WithFields(logrus.Fields{
		"datadir":       datadir,
		"lachesis_addr": lachesisAddress,
		"proxy_addr":    proxyAddress,
		"api_addr":      apiAddress,
		"db":            databaseFile,
		"cache":         dbCache,
	}).Debug("RUN")

	config := proxy.NewConfig(
		proxyAddress,
		lachesisAddress,
		apiAddress,
		datadir,
		pwdFile,
		databaseFile,
		dbCache,
		1*time.Second)

	proxy_, err := proxy.NewProxy(config, logger)
	if err != nil {
		return fmt.Errorf("error building proxy: %s", err)
	}

	proxy_.Run()

	return nil
}

func printVersion(_ *cli.Context) error {
	fmt.Println(version.Version)
	return nil
}

func logLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}

func defaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "LACHESIS")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "EVM")
		} else {
			return filepath.Join(home, ".evm")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
