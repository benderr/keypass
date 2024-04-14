package config

import (
	"errors"
	"flag"
	"regexp"

	"github.com/caarlos0/env/v6"
)

type ServerAddress string

func (address *ServerAddress) String() string {
	return string(*address)
}

func (address *ServerAddress) Set(flagValue string) error {
	if len(flagValue) == 0 {
		return errors.New("empty address not valid")
	}

	reg := regexp.MustCompile(`^([0-9A-Za-z\.]+)?(\:[0-9]+)?$`)

	if !reg.MatchString(flagValue) {
		return errors.New("invalid address and port")
	}

	*address = ServerAddress(flagValue)
	return nil
}

type Config struct {
	Server          ServerAddress `env:"RUN_ADDRESS"`
	DatabaseDsn     string        `env:"DATABASE_URI"`
	SecretKey       string        `env:"KEY"`
	RecordSecretKey string        `env:"RECORD_KEY"`
	PrivateKey      string        `env:"PRIVATE_KEY"`
	PublicKey       string        `env:"PUBLIC_KEY"`
}

var config = Config{
	Server:      ":8080",
	DatabaseDsn: "",
	SecretKey:   "",
}

func init() {
	flag.Var(&config.Server, "a", "address and port to run server")
	flag.StringVar(&config.DatabaseDsn, "d", "", "connection string for postgre")
	flag.StringVar(&config.SecretKey, "k", "", "sha256 based secret key")
	flag.StringVar(&config.RecordSecretKey, "r", "", "secret key for encrypt/decrypt records")
	flag.StringVar(&config.PrivateKey, "private-key", "", "private key file for TLS")
	flag.StringVar(&config.PublicKey, "public-key", "", "public cert file for TLS")
}

func MustLoad() *Config {
	flag.Parse()

	err := env.Parse(&config)

	if err != nil {
		panic(err)
	}

	return &config
}
