package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
	"os"
	"sync"
)

const (
	appPort        = "port"
	guestRole      = "roles.guest"
	clientRole     = "roles.client"
	moderatorRole  = "roles.moderator"
	superAdminRole = "roles.superAdmin"
	dbPort         = "db.port"
	dbUsername     = "db.username"
	dbHost         = "db.host"
	dbName         = "db.name"
	dbSSLMode      = "db.sslMode"
	dbPassword     = "DB_PASSWORD"
	redisHost      = "redis.host"
	redisPort      = "redis.port"
	minioEndpoint  = "minio.endpoint"
	minioAccessKey = "minio.access_key"
	minioSecretKey = "minio.secret_key"
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

type MinIO struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

type Roles struct {
	Guest      string `yaml:"guest"`
	Client     string `yaml:"client"`
	Moderator  string `yaml:"moderator"`
	SuperAdmin string `yaml:"superAdmin"`
}

type Config struct {
	AppPort string `yaml:"port"`
	Roles
	Storage
	Redis
	MinIO
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
		minioInstance := MinIO{
			Endpoint:  viper.GetString(minioEndpoint),
			AccessKey: viper.GetString(minioAccessKey),
			SecretKey: viper.GetString(minioSecretKey),
		}

		rolesInstance := Roles{
			Guest:      viper.GetString(guestRole),
			Client:     viper.GetString(clientRole),
			Moderator:  viper.GetString(moderatorRole),
			SuperAdmin: viper.GetString(superAdminRole),
		}

		instance = &Config{
			AppPort: viper.GetString(appPort),
			Storage: storageInstance,
			Redis:   redisInstance,
			MinIO:   minioInstance,
			Roles:   rolesInstance,
		}
	})

	return instance
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
