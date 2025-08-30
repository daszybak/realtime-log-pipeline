package main

func main() {
	loggers := rest_api_log.NewJSONLoggerFactory(os.Stdout)
	setupLogger, err := loggers.New("setup")
	if err != nil {
		// TODO Format as JSON.
		fmt.Printf("couldn't create setup logger: %v", err)
		os.Exit(1)
	}
	_ = setupLogger.Logf(rest_api_log.LevelInfo, "starting")

	argv := os.Args
	if len(argv) != 3 {
		msg := "usage: %s <config-yaml> <listen-addr>"
		_ = setupLogger.Logf(rest_api_log.LevelError, msg, argv[0])
		os.Exit(1)
	}
	configYamlPath := argv[1]
	listenAddr := argv[2]

	err = run(loggers, setupLogger, configYamlPath, listenAddr)
	if err != nil {
		_ = setupLogger.Logf(rest_api_log.LevelError, "%v", err)
		os.Exit(1)
	}
}

func run(
	loggers *rest_api_log.JSONLoggerFactory,
	setupLogger rest_api_log.UtilStructuredLogger,
	configYamlPath string,
	listenAddr string,
) error {
	_ = setupLogger.Logf(rest_api_log.LevelInfo, "reading '%s'", configYamlPath)
	config, err := readConfig(configYamlPath)
	if err != nil {
		msg := "couldn't read configuration at '%s': %w"
		return fmt.Errorf(msg, configYamlPath, err)
	}

	mySQL := config.MySQL
	_ = setupLogger.Logf(rest_api_log.LevelInfo, "opening MySQL connection")
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			mySQL.User,
			mySQL.Pass,
			mySQL.Addr,
			mySQL.Port,
			mySQL.DB,
		),
	)
	if err != nil {
		return fmt.Errorf("couldn't connect to database: %w", err)
	}

