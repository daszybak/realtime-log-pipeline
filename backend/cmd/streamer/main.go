package main

import (
	"fmt"
	"os"

	"github.com/adshao/go-binance/v2"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := log.With().Str("role", "streamer").Logger()
	setupLogger := log.With().Str("category", "setup").Logger()
	
	argv := os.Args
	if len(argv) != 3 {
	    logger.
			Error().
			Str("expected", "<config-yaml> <listen-addr>").
			Str("got", argv[0]).
			Msg("invalid arguments")
		os.Exit(1)
	}
	configYamlPath := argv[1]
	listenAddr := argv[2]

	err := run(logger, setupLogger, configYamlPath, listenAddr)
	if err != nil {
	    setupLogger.Error().Err(err).Msg("startup failed")
		os.Exit(1)
	}
}

func run(
	streamerLogger zerolog.Logger,
	setupLogger zerolog.Logger,
	configYamlPath string,
	listenAddr string,
) error {
	setupLogger.Info().Str("path", configYamlPath).Msg("reading config")
	config, err := readConfig(configYamlPath)
	if err != nil {
		msg := "couldn't read configuration at '%s': %w"
		return fmt.Errorf(msg, configYamlPath, err)
	}

	pSQL := config.PSQL
	gormLogger := setupLogger.With().Str("component", "gorm").Logger()
	// TODO Use TimescaleDB.
	gormLogger.Info().Str("type", "postgresql").Msg("opening database connection")
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			pSQL.User,
			pSQL.Pass,
			pSQL.Addr,
			pSQL.Port,
			pSQL.DB,
		),
	)
	if err != nil {
		return fmt.Errorf("couldn't connect to database: %w", err)
	}
	db.DB().SetMaxOpenConns(pSQL.MaxConns)
	defer db.Close()
	db.LogMode(true)
	gormLogger.Info().Msg("database connection established")

	httpLogger := setupLogger.With().Str("component", "http").Logger()
	httpLogger.Info().Str("addr", listenAddr).Msg("setting up HTTP router")
	httpLogger = streamerLogger.With().Str("component", "http").Logger()
	// TODO Add logger to middleware for logging requests, responses, websocket.
	// TODO Open connection with best practices to Binance.
	handleBinanceTickerBook(streamerLogger, config.Binance.Symbols)
	// TODO Publish trades to RabbitMQ with telemetry.

	setupLogger.Info().Str("addr", listenAddr).Msg("streamer service ready")

	return nil
}

func handleBinanceTickerBook(
	streamerLogger zerolog.Logger,
	symbols []string,
) {
	binanceLogger := streamerLogger.With().Str("component", "binance").Logger()
	for _, symbol := range symbols {
		_, _, err := binance.WsBookTickerServe(
			symbol, 
			func(event *binance.WsBookTickerEvent) {
				binanceLogger.
					Info().
					Str("best ask price", event.BestAskPrice)
			}, 
			func(err error) {},
		)
		if err != nil {
			binanceLogger.
				Err(err).
				Str("symbol", symbol).
				Msg("could not open binance book ticker websocket connection")
		}
	}
	
}
