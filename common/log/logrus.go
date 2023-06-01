package log

import (
	"github.com/sirupsen/logrus"
)

func Init(level logrus.Level) {
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		DisableQuote:    true,
		TimestampFormat: "15:04:05.000",
	})
}
