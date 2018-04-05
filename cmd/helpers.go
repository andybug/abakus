package cmd

import (
	log "github.com/sirupsen/logrus"
)

func InitLogging() {
	log.SetLevel(log.DebugLevel)
}
