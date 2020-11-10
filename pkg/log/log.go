package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Config .
type Config struct {
	LoggerFile string
}

// NewLog log module
func NewLog(c *Config) error {
	// log.Out= os.Stdout
	// log out to file
	Infoln("🚀 setup logrus out put: ", c.LoggerFile)
	file, err := os.OpenFile(c.LoggerFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Error("\t->Failed to log to file, using default stderr")
		return err
	}
	SetOutput(file)

	// set log formatter like text
	SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	// add call log method but had Additional performance overhead
	// log.SetReportCaller(true)
	SetLevel(logrus.TraceLevel)

	// send errors to an exception tracking service on , and , info to StatsD or log to multiple places simultaneously
	// hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")
	// if err != nil {
	// 	log.Error("Unable to connect to local syslog daemon")
	// } else {
	// 	log.AddHook(hook)
	// }
	return nil
}
