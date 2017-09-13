package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	. "github.com/deshboard/boilerplate-ssh-service/app"
	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	gossh "golang.org/x/crypto/ssh"
)

// newSSHServer creates the main server instance for the service.
func newSSHServer(app *application) serverz.Server {
	serviceChecker := healthz.NewTCPChecker(app.config.SSHAddr, healthz.WithTCPTimeout(2*time.Second))
	app.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	signer, err := createSigner(app.config)
	if err != nil {
		panic(err)
	}

	publicKeys, err := loadRootAuthorizedKeys(app.config)
	if err != nil {
		panic(err)
	}

	return &serverz.AppServer{
		Server: &ssh.Server{
			HostSigners: []ssh.Signer{signer},
			Handler: func(s ssh.Session) {
				io.WriteString(s, fmt.Sprintf("Hello, %s!\n", s.User()))

				a := NewApplication(s)

				a.Run()
			},
			PublicKeyHandler: publicKeyHandler(app.config, publicKeys, app.logger),
		},
		Name:   "ssh",
		Addr:   serverz.NewAddr("tcp", app.config.SSHAddr),
		Logger: app.logger,
	}
}

// createSigner creates a host key signer.
func createSigner(config *configuration) (ssh.Signer, error) {
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

// loadRootAuthorizedKeys loads authorized keys for the root user.
func loadRootAuthorizedKeys(config *configuration) ([]ssh.PublicKey, error) {
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

// publicKeyHandler handles public key authentication.
func publicKeyHandler(config *configuration, keys []ssh.PublicKey, logger log.Logger) ssh.PublicKeyHandler {
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
