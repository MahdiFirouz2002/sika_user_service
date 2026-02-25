package config

import "os"

type config struct {
	// database
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  bool

	// server
	HTTPServerHost string
	HTTPServerPort string
}

func Load() *config {
	return &config{
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseSSLMode:  false,

		HTTPServerHost: os.Getenv("HTTP_SERVER_HOST"),
		HTTPServerPort: os.Getenv("HTTP_SERVER_PORT"),
	}
}
