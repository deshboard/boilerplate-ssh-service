package app

import "time"

// Configuration holds any kind of config that is necessary for running
type Configuration struct {
	Environment     string        `default:"production"`
	ShutdownTimeout time.Duration `default:"3s" split_words:"true"`

	AirbrakeEnabled   bool   `split_words:"true"`
	AirbrakeEndpoint  string `split_words:"true"`
	AirbrakeProjectID int64  `envconfig:"airbrake_project_id"`
	AirbrakeAPIKey    string `envconfig:"airbrake_api_key"`

	FluentdEnabled bool   `split_words:"true"`
	FluentdHost    string `split_words:"true"`
	FluentdPort    int    `split_words:"true"`
}
