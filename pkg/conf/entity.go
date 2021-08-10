package conf

var App = new(configuration)

type configuration struct {
	Server struct {
		ListenAddr string `default:""`
		Port       int    `default:"9000"`
		TimeFormat string `json:"time_format"`
		Log        string `default:""`
		AccessLog  string
	}
}
