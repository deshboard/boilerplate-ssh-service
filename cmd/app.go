package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/deshboard/boilerplate-ssh-service/app"
	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/fxt"
	"github.com/goph/healthz"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	gossh "golang.org/x/crypto/ssh"
)

// AuthorizedKeys holds the full list of authorized keys.
type AuthorizedKeys []ssh.PublicKey

// Err accepts an error which causes the application to stop.
type Err <-chan error

// SSHServerParams provides a set of dependencies for the service constructor.
type SSHServerParams struct {
	dig.In

	Config           *Config
	Signer           ssh.Signer
	PublicKeyHandler ssh.PublicKeyHandler
	Logger           log.Logger        `optional:"true"`
	HealthCollector  healthz.Collector `optional:"true"`
	Lifecycle        fxt.Lifecycle
}

// NewSSHServer returns a new SSH server.
func NewSSHServer(params SSHServerParams) Err {
	logger := params.Logger
	if logger == nil {
		logger = log.NewNopLogger()
	}

	logger = log.With(logger, "server", "ssh")

	if params.HealthCollector != nil {
		params.HealthCollector.RegisterChecker(healthz.ReadinessCheck, healthz.NewTCPChecker(params.Config.SSHAddr, healthz.WithTCPTimeout(2*time.Second)))
	}

	server := &ssh.Server{
		HostSigners: []ssh.Signer{params.Signer},
		Handler: func(s ssh.Session) {
			io.WriteString(s, fmt.Sprintf("Hello, %s!\n", s.User()))

			a := app.NewApplication(s)

			a.Run()
		},
		PublicKeyHandler: params.PublicKeyHandler,
	}

	errCh := make(chan error, 1)

	params.Lifecycle.Append(fxt.Hook{
		OnStart: func(ctx context.Context) error {
			network := "tcp"
			addr := params.Config.SSHAddr

			// Listen on loopback interface in development mode
			if params.Config.Environment == "development" && addr[0] == ':' {
				addr = "127.0.0.1" + addr
			}

			level.Info(logger).Log(
				"msg", "listening on address",
				"addr", addr,
				"network", network,
			)

			lis, err := net.Listen(network, addr)
			if err != nil {
				return errors.WithStack(err)
			}

			go func() {
				errCh <- server.Serve(lis)
			}()

			return nil
		},
		OnStop:  server.Shutdown,
		OnClose: server.Close,
	})

	return errCh
}

// NewSigner creates a host key signer.
func NewSigner(config *Config) (ssh.Signer, error) {
	if config.HostPrivateKey == "" && config.HostPrivateKeyFile != "" {
		file, err := os.Open(config.HostPrivateKeyFile)
		if err != nil {
			return nil, err
		}

		privateKey, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		config.HostPrivateKey = string(privateKey)
	}

	var privateKey *rsa.PrivateKey
	var err error

	// Generate host key if none configured
	if config.HostPrivateKey == "" {
		privateKey, err = rsa.GenerateKey(rand.Reader, 768)
	} else {
		block, _ := pem.Decode([]byte(config.HostPrivateKey))

		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	}

	if err != nil {
		return nil, err
	}

	return gossh.NewSignerFromKey(privateKey)
}

// LoadRootAuthorizedKeys loads authorized keys for the root user.
func LoadRootAuthorizedKeys(config *Config) (AuthorizedKeys, error) {
	if config.RootAuthorizedKeys == "" && config.RootAuthorizedKeysFile != "" {
		file, err := os.Open(config.RootAuthorizedKeysFile)
		if err != nil {
			return nil, err
		}

		authorizedKeys, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		config.RootAuthorizedKeys = string(authorizedKeys)
	}

	authorizedKeysBytes := []byte(config.RootAuthorizedKeys)

	authorizedKeys := []ssh.PublicKey{}
	authorizedKeysMap := map[string]bool{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			return nil, err
		}

		// Avoid duplicates
		if _, ok := authorizedKeysMap[string(pubKey.Marshal())]; !ok {
			authorizedKeysMap[string(pubKey.Marshal())] = true
			authorizedKeys = append(authorizedKeys, pubKey)
		}

		authorizedKeysBytes = rest
	}

	return authorizedKeys, nil
}

// NewPublicKeyHandler returns a new function which handles public key authentication.
func NewPublicKeyHandler(config *Config, keys AuthorizedKeys, logger log.Logger) ssh.PublicKeyHandler {
	return func(ctx ssh.Context, key ssh.PublicKey) bool {
		if ctx.User() == "root" {
			if !config.RootLoginAllowed {
				level.Info(logger).Log(
					"msg", "Root login attempt when disabled",
					"remote_addr", ctx.RemoteAddr(),
				)

				return false
			}

			for _, k := range keys {
				if ssh.KeysEqual(key, k) {
					return true
				}
			}
		} else {
			// Add user authentication here
		}

		return false
	}
}
