package main

import (
	"fmt"
	"iis_server/config"
	"iis_server/scheduler"
	"os"

	"go.uber.org/zap"
)

func main() {
	err := setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot setup program environment\n")
		os.Exit(1)
	}

	scheduler.Start()
}

func setup() error {
	_ = os.Mkdir(config.UPLOAD_FOLDER, 0755)
	_ = os.Mkdir(config.TMP_FOLDER, 0755)
	_ = os.Mkdir(config.XML_SCHEMAS_FOLDER, 0755)

	logger, err := zap.NewDevelopment()
	if err != nil {
		zap.ReplaceGlobals(zap.NewNop())
	} else {
		zap.ReplaceGlobals(logger)
		zap.S().Infof("Console logger setup")
	}
	return nil
}
