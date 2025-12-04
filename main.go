package main

import (
	"github.com/kiritosuki/sysgo/tools"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := tools.RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
