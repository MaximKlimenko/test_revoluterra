package postgres

import (
	"fmt"

	"github.com/MaximKlimenko/scheduler/internal/config"
	"github.com/MaximKlimenko/scheduler/internal/storages"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connector struct {
	DB *gorm.DB
}

func NewConnector(cfg *config.Config) (*Connector, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	cnt := &Connector{
		DB: db,
	}

	db.AutoMigrate(&storages.Job{})

	return cnt, nil
}
