package configs

import (
	"fmt"
	"strings"
)

type Database struct {
	Driver string `yaml:"driver" env:"DRIVER, overwrite"`

	Host     string `yaml:"host" env:"HOST, overwrite"`
	Port     int    `yaml:"port" env:"PORT, overwrite"`
	User     string `yaml:"user" env:"USER, overwrite"`
	Password string `yaml:"password" env:"PASSWORD, overwrite"`
	Name     string `yaml:"name" env:"NAME, overwrite"`

	SSLMode     string `yaml:"ssl_mode" env:"SSL_MODE, overwrite"`           // disable, require, verify-ca, verify-full
	SSLCert     string `yaml:"ssl_cert" env:"SSL_CERT, overwrite"`           // client certificate (client.crt)
	SSLKey      string `yaml:"ssl_key" env:"SSL_KEY, overwrite"`             // client key (client.key)
	SSLRootCert string `yaml:"ssl_root_cert" env:"SSL_ROOT_CERT, overwrite"` // root certificate (ca.crt)

	DSN string `yaml:"dsn" env:"DSN, overwrite"`
}

func (db *Database) GetConnectionString() string {
	switch db.Driver {
	case "sqlite", "sqlite3":
		if db.DSN == "" {
			return "audoctl.db"
		}
		return db.DSN
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			db.User, db.Password, db.Host, db.Port, db.Name)
	case "postgres":
		parts := []string{
			fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", db.Host, db.Port, db.User, db.Password, db.Name),
		}

		// SSLMode always should be present (default require)
		parts = append(parts, fmt.Sprintf("sslmode=%s", db.SSLMode))

		// only if defined
		if db.SSLCert != "" {
			parts = append(parts, fmt.Sprintf("sslcert=%s", db.SSLCert))
		}
		if db.SSLKey != "" {
			parts = append(parts, fmt.Sprintf("sslkey=%s", db.SSLKey))
		}
		if db.SSLRootCert != "" {
			parts = append(parts, fmt.Sprintf("sslrootcert=%s", db.SSLRootCert))
		}

		return strings.Join(parts, " ")
	default:
		return ""
	}
}
