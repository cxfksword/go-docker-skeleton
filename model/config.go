package model

type Config struct {
	Log           string `json:"log"`
	AccessLog     string `json:"access_log"`
	Debug         bool   `json:"debug"`
	TimeFormat    string `json:"time_format"`
	LogTimeFormat string `json:"log_time_format"`
}
