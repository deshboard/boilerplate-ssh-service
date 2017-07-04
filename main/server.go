package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/deshboard/boilerplate-ssh-service/app"
	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/goph/stdlib/ext"
	opentracing "github.com/opentracing/opentracing-go"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// newServer creates the main server instance for the service.
func newServer(config *configuration, logger log.Logger, errorHandler emperror.Handler, tracer opentracing.Tracer, healthCollector healthz.Collector, metricsReporter interface{}) (serverz.Server, ext.Closer) {
	serviceChecker := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	signer, err := createSigner(config)
	if err != nil {
		panic(err)
	}

	publicKeys, err := loadRootAuthorizedKeys(config)
	if err != nil {
		panic(err)
	}

	return &serverz.NamedServer{
		Server: &ssh.Server{
			HostSigners:      []ssh.Signer{signer},
			Handler:          handler,
			PublicKeyHandler: publicKeyHandler(config, publicKeys, logger),
		},
		Name: "ssh",
	}, ext.NoopCloser
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

	block, _ := pem.Decode([]byte(config.HostPrivateKey))

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
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

// handler is the SSH handler function.
func handler(s ssh.Session) {
	prompt := fmt.Sprintf("%s@deshboard:$ ", s.User())
	t := terminal.NewTerminal(s, prompt)

	io.WriteString(s, fmt.Sprintf("Hello, %s!\n", s.User()))

	app := app.NewApplication(s, t, prompt)

	for {
		line, err := t.ReadLine()

		// Ctrl+D received
		if err == io.EOF {
			io.WriteString(s, "\n")
			s.Exit(0)
		} else if err == nil {
			if line != "" {
				args := strings.Split(line, " ")
				app.Execute(args)
			}
		}
	}
}
