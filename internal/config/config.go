package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/zh0vtyj/alliancecup-server/pkg/logging"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	appPort            = "port"
	domain             = "domain"
	guestRole          = "roles.guest"
	clientRole         = "roles.client"
	moderatorRole      = "roles.moderator"
	superAdminRole     = "roles.superAdmin"
	dbPort             = "db.port"
	dbUsername         = "db.username"
	dbHost             = "db.host"
	dbName             = "db.name"
	dbSSLMode          = "db.sslMode"
	dbPassword         = "DB_PASSWORD"
	redisHost          = "redis.host"
	redisPort          = "redis.port"
	corsAllowedOrigins = "cors.allowedOrigins"
	minioEndpoint      = "minio.endpoint"
	minioAccessKey     = "MINIO_ACCESS_KEY"
	minioSecretKey     = "MINIO_SECRET_KEY"
	passwordSaltEnv    = "PASSWORD_SALT"
	jwtSigningKeyEnv   = "JWT_SIGNING_KEY"
)

type (
	Auth struct {
		JWT          JWT
		PasswordSalt string `env:"PASSWORD_SALT"`
	}

	JWT struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string        `env:"JWT_SIGNING_KEY"`
	}

	Redis struct {
		Host string `yml:"host"`
		Port string `yml:"port"`
	}

	Storage struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `env:"DB_PASSWORD"`
		DBName   string `yaml:"name"`
		SSLMode  string `yaml:"sslMode"`
	}

	MinIO struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `env:"MINIO_ACCESS_KEY"`
		SecretKey string `env:"MINIO_SECRET_KEY"`
	}

	Roles struct {
		Guest      string `yaml:"guest"`
		Client     string `yaml:"client"`
		Moderator  string `yaml:"moderator"`
		SuperAdmin string `yaml:"superAdmin"`
	}

	Cors struct {
		AllowedOrigins []string `yaml:"allowedOrigins"`
	}

	Config struct {
		Domain  string `yaml:"domain"`
		AppPort string `yaml:"port"`
		Auth
		Roles
		Cors
		Storage
		Redis
		MinIO
	}
)

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()

		logger.Info("initializing .yml file")
		if err := initConfig(); err != nil {
			logger.Fatalf("error while initializing .yml file, %v", err)
		}

		logger.Info("initializing .env file")
		if err := godotenv.Load(); err != nil {
			logger.Fatalf("error while initializing .env file, %v", err)
		}

		var authInstance Auth
		if err := viper.UnmarshalKey("auth", &authInstance.JWT); err != nil {
			return
		}

		authInstance.PasswordSalt = os.Getenv(passwordSaltEnv)
		authInstance.JWT.SigningKey = os.Getenv(jwtSigningKeyEnv)

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
			AccessKey: os.Getenv(minioAccessKey),
			SecretKey: os.Getenv(minioSecretKey),
		}

		rolesInstance := Roles{
			Guest:      viper.GetString(guestRole),
			Client:     viper.GetString(clientRole),
			Moderator:  viper.GetString(moderatorRole),
			SuperAdmin: viper.GetString(superAdminRole),
		}

		corsInstance := Cors{
			AllowedOrigins: strings.Split(viper.GetString(corsAllowedOrigins), ","),
		}

		instance = &Config{
			Domain:  viper.GetString(domain),
			AppPort: viper.GetString(appPort),
			Auth:    authInstance,
			Cors:    corsInstance,
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
