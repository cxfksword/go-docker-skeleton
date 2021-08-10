package app

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cxfksword/go-docker-skeleton/pkg/conf"
	"github.com/cxfksword/go-docker-skeleton/pkg/log"
	l "github.com/cxfksword/go-docker-skeleton/pkg/log"
	"github.com/cxfksword/go-docker-skeleton/pkg/mode"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	Default_Time_Format     = "2006-01-02 15:04:05"
	Default_Log_Time_Format = "2006-01-02 15:04:05.000"
)

type application struct {
	AppName string
	AppDesc string
	Mode    string
	Version VersionInfo
	Port    int

	ConfigFilePath string

	DebugLogLevel   bool
	VerboseLogLevel bool
}

var currentApp *application

func New(name string, desc string, mode string, vInfo VersionInfo) *application {
	return &application{
		AppName: name,
		AppDesc: desc,
		Mode:    mode,
		Version: vInfo,
	}
}

func (a *application) Run(r *gin.Engine) {
	currentApp = a

	conf.Init(a.AppName, a.ConfigFilePath)
	l.AddFileOutput(conf.App.Server.Log)
	if a.DebugLogLevel || a.Mode == mode.Dev {
		l.SetDebugLevel()
	}
	if a.VerboseLogLevel {
		l.SetTraceLevel()
	}

	// set gin access log
	accessWriter := createAccessRollingLogFile(conf.App.Server.AccessLog)
	if accessWriter != nil {
		gin.DefaultWriter = io.MultiWriter(accessWriter, os.Stderr)
		log.Infof("Gin access log file path: %s", conf.App.Server.AccessLog)
	}

	// log.Info().Msgf("Config file path: %s", a.ConfigFilePath)
	// log.Info().Msgf("Log file path: %s", a.Config.Log)
	// log.Info().Msgf("Gin access log file path: %s", a.Config.AccessLog)

	// run web server
	if a.Port <= 0 {
		a.Port = conf.App.Server.Port
	}
	err := r.Run(fmt.Sprintf("%s:%d", conf.App.Server.ListenAddr, a.Port))
	if err != nil {
		log.Error(err.Error())
	}
}

func (a *application) ReloadConfig() {
	conf.Reload()
}

func DevMode() bool {
	return currentApp.Mode == mode.Dev
}

func createAccessRollingLogFile(logPath string) io.Writer {
	logDir := path.Dir(logPath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Errorf("Can't create log directory. path: %s, err: %s", logPath, err)
		return nil
	}

	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    500, // megabytes
		MaxBackups: 31,
		MaxAge:     7, //daysdays
	}
}
