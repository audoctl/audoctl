package gormdb

import (
	"fmt"
	"strings"

	"github.com/audoctl/audoctl/configs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormDB struct {
	Db *gorm.DB
}

// Connect opens a Gorm DB connection based on Database config
func Connect(cfg *configs.Database, logger *GormLogger) (*GormDB, error) {
	dsn := cfg.GetConnectionString()
	var dialector gorm.Dialector

	switch strings.ToLower(cfg.Driver) {
	case "postgres", "pg":
		dialector = postgres.Open(dsn)
	case "mysql":
		dialector = mysql.Open(dsn)
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	gormCfg := &gorm.Config{}
	if logger != nil {
		gormCfg.Logger = logger
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, err
	}

	return &GormDB{db}, nil
}

// Disconnect closes the underlying SQL connection
func (g *GormDB) Disconnect() error {
	sqlDB, err := g.Db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
