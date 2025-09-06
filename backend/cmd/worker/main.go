package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adshao/go-binance/v2"
	"github.com/daszybak/realtime-log-pipeline/backend/internal/binance_rabbitmq"
	"github.com/daszybak/realtime-log-pipeline/backend/pkg/rabbitmq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := log.With().Str("service", "worker").Logger()
	setupLogger := log.With().Str("category", "setup").Logger()

	argv := os.Args
	if len(argv) != 3 {
		logger.
			Error().
			Str("expected", "<config-yaml> <listen-addr>").
			Str("got", argv[0]).
			Msg("Invalid arguments passed.")
		os.Exit(1)
	}
	configYamlPath := argv[1]
	listenAddr := argv[2]

	err := run(logger, setupLogger, configYamlPath, listenAddr)
	if err != nil {
		setupLogger.Error().Err(err).Msg("Startup failed.")
		os.Exit(1)
	}
}

func run(
	workerLogger zerolog.Logger,
	setupLogger zerolog.Logger,
	configYamlPath string,
	listenAddr string,
) error {
	setupLogger.Info().Str("path", configYamlPath).Msg("Reading config.")
	config, err := readConfig(configYamlPath)
	if err != nil {
		msg := "couldn't read configuration at '%s': %w"
		return fmt.Errorf(msg, configYamlPath, err)
	}

	pSQLLogger := setupLogger.With().Str("component", "postgresql").Logger()
	pSQLLogger.Info().Str("type", "postgresql").Msg("Opening database connection.")
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%d",
			config.PSQL.User,
			config.PSQL.Pass,
			config.PSQL.DB,
			config.PSQL.Addr,
			config.PSQL.Port,
		),
	)
	if err != nil {
		return fmt.Errorf("couldn't create database config: %w", err)
	}
	dbConfig.MaxConns = int32(config.PSQL.MaxConns)
	_, err = pgxpool.NewWithConfig(
		rootCtx,
		dbConfig,
	)
	if err != nil {
		return fmt.Errorf("couldn't connect to database: %w", err)
	}
	pSQLLogger.Info().Msg("Database connection established.")

	httpLogger := setupLogger.With().Str("component", "http").Logger()
	httpLogger.Info().Str("addr", listenAddr).Msg("Setting up HTTP router.")
	httpLogger = workerLogger.With().Str("component", "http").Logger()

	setupLogger.Info().Str("addr", listenAddr).Msg("worker service ready.")

	httpRouterSetup(workerLogger)

	rabbitMQClient, err := rabbitmq.New("worker", config.RabbitMQ.URL)
	if err != nil {
		return fmt.Errorf("couldn't set up RabbitMQ Client: %w", err)
	}
	setupLogger.Info().Str("component", "rabbitmq").Msg("RabbitMQ client instantiated")
	binanceRabbitMQClient, err := binance_rabbitmq.New[*binance.WsBookTickerEvent](rabbitMQClient)
	if err != nil {
		return fmt.Errorf("couldn't set up Binance RabbitMQ Client: %w", err)
	}
	binanceRabbitMQLogger := workerLogger.With().Str("component", "binance_rabbitmq").Logger()
	binanceRabbitMQLogger.Info().Msg("Binance RabbitMQ Client instantiated")

	messages, err := binanceRabbitMQClient.Consume(rootCtx)
	for m := range messages {
		// TODO Add OpenTelem to handle `trace_id`.
		binanceRabbitMQLogger.
			Info().
			Str("symbol", m.Data.Symbol).
			Str("best_bid_price", m.Data.BestBidPrice).
			Str("best_asl_price", m.Data.BestAskPrice).
			Msg("Received event from RabbitMQ")
	}

	setupLogger.Info().Msg("worker is running, press Ctrl+C to stop.")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	return nil
}

func httpRouterSetup(
	streamerLogger zerolog.Logger,
) {
	// TODO Setup router using `chi`.
	streamerLogger.Info().Msg("HTTP router is setting up...")
}

func handleBinanceTickerBook(
	streamerLogger zerolog.Logger,
	symbols []string,
	handler binance.WsBookTickerHandler,
) ([]chan struct{}, []chan struct{}) {
	binanceLogger := streamerLogger.With().Str("component", "binance").Logger()

	var stopChannels []chan struct{}
	var doneChannels []chan struct{}

	for _, symbol := range symbols {
		binanceLogger.Info().Str("symbol", symbol).Msg("Setting up websocket connection.")
		stopC, doneC, err := binance.WsBookTickerServe(
			symbol,
			func(event *binance.WsBookTickerEvent) {
				handler(event)
				// binanceLogger.
				// 	Info().
				// 	Str("symbol", symbol).
				// 	Str("best_ask_price", event.BestAskPrice).
				// 	Str("best_bid_price", event.BestBidPrice).
				// 	Msg("Received ticker update.")
			},
			func(err error) {
				binanceLogger.
					Err(err).
					Str("symbol", symbol).
					Msg("Failed to receive book ticker event.")
			},
		)
		if err != nil {
			// TODO Add exponential back-off reconnection.
			binanceLogger.
				Err(err).
				Str("symbol", symbol).
				Msg("Could not open binance book ticker websocket connection.")
			continue
		}

		binanceLogger.
			Info().
			Str("symbol", symbol).
			Msg("Websocket connection established.")
		stopChannels = append(stopChannels, stopC)
		doneChannels = append(doneChannels, doneC)

		// Monitor connection health in background.
		// TODO We should monitor with some tool.
		go func(symbol string, doneC chan struct{}) {
			<-doneC
			binanceLogger.Warn().Str("symbol", symbol).Msg("Websocket connection closed.")
		}(symbol, doneC)
	}

	return stopChannels, doneChannels
}
