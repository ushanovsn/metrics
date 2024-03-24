package logger

import (
	"io"
	"os"

	"github.com/ushanovsn/metrics/internal/options"

	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
)

func InitLogger(opt *options.LoggerOpt) {
	log := logrus.New()

	// init file
	f, opnFileErr := os.OpenFile(opt.FileName+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	// use multiwriter
	if opnFileErr != nil {
		log.Out = os.Stdout
	} else {
		opt.SetFile(f)
		log.Out = io.MultiWriter(f, os.Stdout)
	}

	// change formatter
	log.Formatter = &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%",
	}

	// Option of logger level
	lvl, levelErr := logrus.ParseLevel(opt.Level)
	if levelErr == nil {
		log.Level = lvl
	} else {
		log.Level = logrus.DebugLevel
	}

	// Option of logger report file entering
	log.ReportCaller = opt.Report

	// logging errors at init process
	if opnFileErr != nil {
		log.Errorf("Error when file open: %s\n", opnFileErr)
	}

	if levelErr != nil {
		log.Errorf("Error when parse level: %s\n", levelErr)
	}

	opt.SetLogger(log)
}
