package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/app"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

var (
	versionFlag = flag.Bool("v", false, "Print version and exit")
	configFlag  = flag.String("c", "", "Path to config file")
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Brevity Server v%s\n", configs.BrevityApp.App.Version)
		return
	}

	// Initialize basic logger with default configuration
	err := logger.Init(logger.Config{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		fmt.Printf("Failed to initialize basic logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	log := logger.Get()

	// Load configuration
	configPath := configs.GetConfigPath()
	if *configFlag != "" {
		configPath = *configFlag
	}

	cfg, err := configs.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Failed to load config",
			logger.NamedError("error", err),
			logger.String("hint", "Ensure app.yaml exists in configs/ directory or provide path with -c flag"))
	}

	// Reinitialize logger with config settings
	loggerCfg := logger.Config{
		Level:    cfg.Logger.Level,
		Format:   cfg.Logger.Format,
		FilePath: cfg.Logger.FilePath,
	}

	err = logger.Init(loggerCfg)
	if err != nil {
		log.Fatal("Failed to initialize logger", logger.ErrorField(err))
	}
	defer logger.Sync()

	log = logger.Get()

	// Create server
	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatal("Failed to create server", logger.ErrorField(err))
	}

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal("Server exited with error", logger.ErrorField(err))
	}
}
