package service_discovery

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/exp/slog"
)

type Config struct {
	LogFormat string     `json:"log_format,omitempty"`
	LogLevel  slog.Level `json:"log_level,omitempty"`
	HTTPPort  int        `json:"http_port,omitempty"`
}

var (
	defaultConfig = Config{
		LogFormat: "json",
		HTTPPort:  8080,
		LogLevel:  slog.LevelInfo,
	}
)

func LoadConfig(args []string) Config {
	var config = defaultConfig

	loadEnv(&config)
	loadFlags(&config, args)

	return config
}

func loadEnv(config *Config) {
	if v, ok := getStringEnv("SD_LOG_FORMAT"); ok {
		config.LogFormat = v
	}

	if v, ok := getStringEnv("SD_LOG_LEVEL"); ok {
		config.LogLevel = parseLevel(v, config.LogLevel)
	}

	if v, ok := getIntEnv("SD_HTTP_PORT"); ok {
		config.HTTPPort = v
	}
}

func loadFlags(config *Config, args []string) {
	var (
		logFormatName = "log-format"
		logLevelName  = "log-level"
		httpPortName  = "http-port"
	)

	fs := flag.NewFlagSet("service_discovery", flag.ExitOnError)

	fs.String(logFormatName, config.LogFormat, `Log messages format ("text", "json")`)
	fs.String(
		logLevelName, config.LogLevel.String(),
		fmt.Sprintf(
			"Log level (%q, %q, %q, %q)",
			slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError,
		),
	)
	fs.Int(httpPortName, config.HTTPPort, "HTTP Port where the service listen on")

	_ = fs.Parse(args)

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case logFormatName:
			config.LogFormat = f.Value.String()

		case logLevelName:
			config.LogLevel = parseLevel(f.Value.String(), config.LogLevel)

		case httpPortName:
			config.HTTPPort = parserInt(f.Value.String(), config.HTTPPort)
		}
	})
}

func getStringEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func getIntEnv(key string) (int, bool) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return 0, false
	}

	if intVal, err := strconv.Atoi(v); err == nil {
		return intVal, true
	}
	return 0, false
}

func parseLevel(lvl string, defaultLevel slog.Level) slog.Level {
	var result slog.Level
	if err := result.UnmarshalText([]byte(lvl)); err != nil {
		return defaultLevel
	}
	return result
}

func parserInt(val string, defaultValue int) int {
	if v, err := strconv.Atoi(val); err == nil {
		return v
	}
	return defaultValue
}
