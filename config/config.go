package config

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	Dsn string
}
type AuthConfig struct {
	Secret string
}

func newEnvType() *EnvType {
	return &EnvType{
		Test:  "test",
		Prod:  "prod",
		Local: "local",
	}
}

type EnvType struct {
	Test  string
	Prod  string
	Local string
}

func (c *Config) GetAuthSecret() string {
	return os.Getenv("SECRET")
}

func (c *Config) On(s string) {

}

func LoadConfig(env string) *Config {
	var envFilePath string
	if !slices.Contains([]string{"test", "prod", "local"}, env) {
		log.Fatal("Env variables not found")
	}
	envFilePath = fmt.Sprintf("%s.env", env)
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Println("Error loading .env file, using default config", err)
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
		Port: os.Getenv("PORT"),
	}
}
