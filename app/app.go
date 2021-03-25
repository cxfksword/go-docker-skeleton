package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	l "github.com/cxfksword/go-docker-skeleton/log"
	"github.com/cxfksword/go-docker-skeleton/mode"
	"github.com/cxfksword/go-docker-skeleton/model"
	"github.com/cxfksword/go-docker-skeleton/utils"
	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

const (
	Default_Time_Format     = "2006-01-02 15:04:05"
	Default_Log_Time_Format = "2006-01-02 15:04:05.000"
)

type application struct {
	AppName string
	AppDesc string
	Mode    string
	Version model.VersionInfo
	Port    int

	ConfigFilePath string
	Config         *model.Config

	DebugLogLevel   bool
	VerboseLogLevel bool
}

var currentApp *application

func New(mode string, vInfo model.VersionInfo) *application {
	return &application{
		Mode:    mode,
		Version: vInfo,
	}
}

func (a *application) Run(cmd *cli.App, r *gin.Engine) {
	currentApp = a
	a.AppName = cmd.Name
	a.AppDesc = cmd.Description

	cmd.Action = func(c *cli.Context) error {
		// here run after command line argument parse
		a.initialize()

		// run web server
		return r.Run(fmt.Sprintf(":%d", a.Port))
	}

	err := cmd.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func (a *application) initialize() {
	// load file config
	a.Config = a.loadConfig()

	// set log output format
	a.initLogger()

	log.Info().Msgf("Config file path: %s", a.ConfigFilePath)
	log.Info().Msgf("Log file path: %s", a.Config.Log)
	log.Info().Msgf("Gin access log file path: %s", a.Config.AccessLog)
}

func (a *application) loadConfig() *model.Config {
	conf := &model.Config{}

	if a.ConfigFilePath == "" {
		a.ConfigFilePath = filepath.Join(utils.DefaultSavePath(), a.AppName+".yaml")
	}
	viper.SetConfigFile(a.ConfigFilePath)
	err := viper.ReadInConfig() // Find and read the config file
	if err == nil {
		err = viper.Unmarshal(conf)
		if err != nil {
			fmt.Printf("Unable to decode into struct, %v\n", err)
		}
	} else {
		fmt.Printf("Load config file failed: %s, will ignore.\n", err)
	}

	if conf.LogTimeFormat == "" {
		conf.LogTimeFormat = Default_Log_Time_Format
	}
	if conf.TimeFormat == "" {
		conf.TimeFormat = Default_Time_Format
	}
	if conf.Log == "" {
		conf.Log = filepath.Join(utils.DefaultSavePath(), a.AppName+".log")
	}
	if conf.AccessLog == "" {
		conf.AccessLog = filepath.Join(utils.DefaultSavePath(), a.AppName+"_access.log")
	}

	return conf
}

func (a *application) SaveConfig(newConf model.Config) (*model.Config, error) {
	conf := a.loadConfig()

	if err := mergo.Merge(conf, newConf, mergo.WithOverwriteWithEmptyValue); err != nil {
		fmt.Println(err)
		return nil, err
	}

	data, err := yaml.Marshal(&conf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	f, err := os.OpenFile(a.ConfigFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer f.Close()

	err = ioutil.WriteFile(a.ConfigFilePath, data, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	a.ReloadConfig()
	return conf, nil
}

func (a *application) ReloadConfig() {
	a.Config = a.loadConfig()
}

func (a *application) initLogger() {
	// set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if a.DebugLogLevel || a.Mode == mode.Dev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if a.VerboseLogLevel {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	// set log time format
	zerolog.TimeFieldFormat = a.Config.LogTimeFormat

	// add multi log output
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	fileWriter := l.NewFileFormatWriter(a.createRollingLogFile())
	multiWriter := zerolog.MultiLevelWriter(consoleWriter, fileWriter)
	log.Logger = log.With().Caller().Logger().Output(multiWriter)

	// set gin access log
	gin.DefaultWriter = io.MultiWriter(a.createGinRollingLogFile(), os.Stderr)
}

func (a *application) createRollingLogFile() io.Writer {
	logDir := path.Dir(a.Config.Log)
	if err := os.MkdirAll(logDir, 0744); err != nil {
		panic(fmt.Sprintf("can't create log directory: %s", logDir))
	}

	return &lumberjack.Logger{
		Filename:   a.Config.Log,
		MaxSize:    500, // megabytes
		MaxBackups: 31,
		MaxAge:     7, //daysdays
	}
}

func (a *application) createGinRollingLogFile() io.Writer {
	logDir := path.Dir(a.Config.AccessLog)
	if err := os.MkdirAll(logDir, 0744); err != nil {
		panic(fmt.Sprintf("can't create access log directory: %s", logDir))
	}

	return &lumberjack.Logger{
		Filename:   a.Config.Log,
		MaxSize:    500, // megabytes
		MaxBackups: 31,
		MaxAge:     7, //daysdays
	}
}

func Current() *application {
	return currentApp
}

func Config() *model.Config {
	return currentApp.Config
}

func DevMode() bool {
	return currentApp.Mode == mode.Dev
}
