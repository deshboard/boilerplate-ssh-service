package app

// Configuration holds any kind of config that is necessary for running
type Configuration struct {
	Environment string `default:"production"`

	AirbrakeEnabled   bool   `split_words:"true"`
	AirbrakeHost      string `split_words:"true"`
	AirbrakeProjectID int64  `envconfig:"airbrake_project_id"`
	AirbrakeAPIKey    string `envconfig:"airbrake_api_key"`
}
