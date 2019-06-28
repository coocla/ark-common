package logs

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestLogOutput(t *testing.T) {
	SetupLog("./ok.log", "./err.log", "debug")
	log.Warn("warning")
	log.Info("ok")
	log.Debug("debug")
	log.Error("error")
}
