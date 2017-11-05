package proxy

import (
	"time"

	"github.com/Sirupsen/logrus"
	bproxy "github.com/babbleio/babble/proxy/babble"
	"github.com/babbleio/evm-babble/service"
	"github.com/babbleio/evm-babble/state"
)

//------------------------------------------------------------------------------

type Config struct {
	proxyAddr  string //bind address of this app proxy
	babbleAddr string //address of babble node
	apiAddr    string //address of HTTP API service
	ethDir     string //directory containing eth config
	pwdFile    string //file containing password to unlock ethereum accounts
	timeout    time.Duration
}

func NewConfig(proxyAddr, babbleAddr, apiAddr, ethDir, pwdFile string, timeout time.Duration) Config {
	return Config{
		proxyAddr:  proxyAddr,
		babbleAddr: babbleAddr,
		apiAddr:    apiAddr,
		ethDir:     ethDir,
		pwdFile:    pwdFile,
		timeout:    timeout,
	}
}

//------------------------------------------------------------------------------

type Proxy struct {
	service     *service.Service
	state       *state.State
	babbleProxy *bproxy.SocketBabbleProxy
	submitCh    chan []byte
	logger      *logrus.Logger
}

func NewProxy(config Config, logger *logrus.Logger) (*Proxy, error) {
	submitCh := make(chan []byte)

	state, err := state.NewState(logger)
	if err != nil {
		return nil, err
	}

	service := service.NewService(config.ethDir,
		config.apiAddr,
		config.pwdFile,
		state,
		submitCh,
		logger)

	babbleProxy, err := bproxy.NewSocketBabbleProxy(config.babbleAddr,
		config.proxyAddr,
		config.timeout)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		service:     service,
		state:       state,
		babbleProxy: babbleProxy,
		submitCh:    submitCh,
		logger:      logger,
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
			p.logger.Debug("proxy about to submit tx")
			if err := p.babbleProxy.SubmitTx(tx); err != nil {
				p.logger.WithError(err).Error("SubmitTx")
			}
			p.logger.Debug("proxy submitted tx")
		case tx := <-p.babbleProxy.CommitCh():
			if err := p.state.AppendTx(tx); err != nil {
				p.logger.WithError(err).Error("AppendTx")
				break
			}
			if err := p.state.Commit(); err != nil {
				p.logger.WithError(err).Error("Commit")
				break
			}
		}
	}
}