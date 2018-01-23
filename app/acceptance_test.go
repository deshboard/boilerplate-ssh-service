package app

import (
	"context"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/goph/fxt"
)

func AppContext(app *fxt.App) func(s *godog.Suite) {
	return func(s *godog.Suite) {
		s.BeforeScenario(func(scenario interface{}) {
			err := app.Start(context.Background())
			if err != nil {
				panic(err)
			}
		})

		s.AfterScenario(func(scenario interface{}, err error) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err = app.Stop(ctx)
			if err != nil {
				panic(err)
			}

			app.Close()
		})
	}
}
