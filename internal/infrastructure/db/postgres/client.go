package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host            string        `env:"PG_HOST" envDefault:"localhost"`
	Port            string        `env:"PG_PORT" envDefault:"5432"`
	User            string        `env:"PG_USER" envDefault:"postgres"`
	Password        string        `env:"PG_PASSWORD" envDefault:"postgres"`
	DBName          string        `env:"PG_DB" envDefault:"shortened_urls"`
	SSLMode         string        `env:"PG_SSLMODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"5m"`
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME" envDefault:"1m"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewFromConfig(ctx context.Context) (*sqlx.DB, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening postgres: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}
	return db, nil
}
