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

	"github.com/deshboard/boilerplate-ssh-service/app"
	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz/aio"
	"github.com/goph/stdlib/net"
	gossh "golang.org/x/crypto/ssh"
)

// newSSHServer creates the main server instance for the service.
func newSSHServer(appCtx *application) *aio.Server {
	serviceChecker := healthz.NewTCPChecker(appCtx.config.SSHAddr, healthz.WithTCPTimeout(2*time.Second))
	appCtx.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	signer, err := createSigner(appCtx.config)
	if err != nil {
		panic(err)
	}

	publicKeys, err := loadRootAuthorizedKeys(appCtx.config)
	if err != nil {
		panic(err)
	}

	return &aio.Server{
		Server: &ssh.Server{
			HostSigners: []ssh.Signer{signer},
			Handler: func(s ssh.Session) {
				io.WriteString(s, fmt.Sprintf("Hello, %s!\n", s.User()))

				app := app.NewApplication(s)

				app.Run()
			},
			PublicKeyHandler: publicKeyHandler(appCtx.config, publicKeys, appCtx.logger),
		},
		Name: "ssh",
		Addr: net.ResolveVirtualAddr("tcp", appCtx.config.SSHAddr),
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
