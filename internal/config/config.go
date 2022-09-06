package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
	"os"
	"sync"
)

const (
	appPort    = "port"
	dbPort     = "db.port"
	dbUsername = "db.username"
	dbHost     = "db.host"
	dbName     = "db.name"
	dbSSLMode  = "db.sslMode"
	dbPassword = "DB_PASSWORD"
	redisHost  = "redis.host"
	redisPort  = "redis.port"
)

type Redis struct {
	Host string `yml:"host"`
	Port string `yml:"port"`
}

type Storage struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `yaml:"name"`
	SSLMode  string `yaml:"sslMode"`
}

type Config struct {
	AppPort string `yaml:"port"`
	Storage
	Redis
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()

		logger.Info("initializing .yml file")
		if err := initConfig(); err != nil {
			logger.Panic("panic while initializing .yml file")
			panic(err)
		}

		logger.Info("initializing .env file")
		if err := godotenv.Load(); err != nil {
			logger.Panic("panic while initializing .env file")
			panic(err)
		}

		redisInstance := Redis{
			Host: viper.GetString(redisHost),
			Port: viper.GetString(redisPort),
		}
		storageInstance := Storage{
			Host:     viper.GetString(dbHost),
			Port:     viper.GetString(dbPort),
			Username: viper.GetString(dbUsername),
			Password: os.Getenv(dbPassword),
			DBName:   viper.GetString(dbName),
			SSLMode:  viper.GetString(dbSSLMode),
		}
		instance = &Config{
			AppPort: viper.GetString(appPort),
			Storage: storageInstance,
			Redis:   redisInstance,
		}
	})

	return instance
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
