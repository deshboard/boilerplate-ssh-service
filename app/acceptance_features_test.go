// +build acceptance

package app

import (
	"github.com/deshboard/boilerplate-service/pkg/app"
	"github.com/goph/fxt"
	"github.com/goph/fxt/test"
	"go.uber.org/fx"
)

func init() {
	acceptanceRunner = test.NewGodogRunner()

	var config Config

	a := fxt.New(
		fx.NopLogger,
		fx.Provide(newConfig, newApplicationInfo),
		Module,
		fx.Populate(&config),
	)

	acceptanceRunner.RegisterFeatureContext(AppContext(a))
	app.RegisterSuite(acceptanceRunner)
}
