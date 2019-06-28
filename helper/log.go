package helper

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	// Initialize the logging component
	logFile := fmt.Sprintf("%s/kucoin-go-level2-demo-%s.log", ".", time.Now().Format("2006-01-02"))
	logWriter, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Panicf("Open file failed: %s", err.Error())
	}
	mw := io.MultiWriter(os.Stdout, logWriter)

	logrus.SetOutput(mw)
	logrus.SetLevel(logrus.InfoLevel)
}

var logger = log.New(os.Stdout, "", log.LstdFlags)

func Debug(format string, v ...interface{}) {
	logrus.Debugf("[Debug] "+format, v...);
}

func Info(format string, v ...interface{}) {
	logrus.Infof("[Info] "+format, v...);
}

func Warn(format string, v ...interface{}) {
	logrus.Warnf("[Warn] "+format, v...);
}

func Error(format string, v ...interface{}) {
	logrus.Errorf("[Error] "+format, v...);
}

func Fatal(format string, v ...interface{}) {
	logrus.Fatalf("[Fatal] "+format, v...);
}
