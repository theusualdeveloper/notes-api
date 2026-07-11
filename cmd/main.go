package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/theusualdeveloper/notes-api/db"
)

type Config struct {
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_HOST_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
}

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:         ":8080",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	logger := InitSlog()
	var config Config
	err := loadConfig(&config, ".")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	var dsn string
	if config.DBUser == "" ||
		config.DBPassword == "" ||
		config.DBHost == "" ||
		config.DBPort == "" ||
		config.DBName == "" {
		dsn = "postgres://admin:admin@localhost:5332/notes"
	} else {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBName,
		)
	}
	pgxpool, err := db.NewDB(context.Background(), dsn)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("database connected")
	// running migrations
	err = db.Run(context.Background(), pgxpool, "./migrations")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	logger.Info("server is starting on http://localhost:8080")
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("starting server failed", "Err", err.Error())
		return
	}
}

func loadConfig(config *Config, path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("reading config file failed: %w", err)
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return fmt.Errorf("unmarshaling config values failed: %w", err)
	}
	return nil
}

func InitSlog() *slog.Logger {
	options := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, options)
	return slog.New(handler)
}
