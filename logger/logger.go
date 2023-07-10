package logger

import (
	"github.com/sirupsen/logrus"
	"log-parser/config"
)

var Log = logrus.WithFields(logrus.Fields{
	"app": config.GetString("SERVICE_NAME"),
	"env": config.GetString("ENV"),
})

func init() {
	formatter := &logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    false,
	}

	logrus.SetFormatter(formatter)

	logrus.RegisterExitHandler(func() {
		logrus.Info("application will stop probably due to a os signal")
	})

	ll := config.GetString("LOG_LEVEL")
	l, err := logrus.ParseLevel(ll)
	if err != nil {
		Log.WithError(err).Errorf("error parsing log level %s", ll)
		return
	}

	logrus.SetLevel(l)
}
