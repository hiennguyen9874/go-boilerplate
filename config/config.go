package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

var (
	cfg *Config
)

type Config struct {
	Server         ServerConfig
	Postgres       PostgresConfig
	Redis          RedisConfig
	Jwt            JwtConfig
	FirstSuperUser FirstSuperUserConfig
	Logger         Logger
	SmtpEmail      SmtpEmailConfig
	Email          EmailConfig
	TaskRedis      TaskRedisConfig
}

// Server config struct
type ServerConfig struct {
	AppVersion     string
	Port           string
	Mode           string
	ProcessTimeout int
	ReadTimeout    int
	WriteTimeout   int
	MigrateOnStart bool
}

// Logger config
type Logger struct {
	LoggerEncoding string
	LoggerLevel    string
}

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  string
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// Task redis config
type TaskRedisConfig struct {
	TaskRedisAddr string
	TaskRedisDb   int
}

// Jwt config
type JwtConfig struct {
	// JwtSecretKey                string
	JwtIssuer                     string
	JwtAccessTokenExpireDuration  int64
	JwtAccessTokenPrivateKey      string
	JwtAccessTokenPublicKey       string
	JwtRefreshTokenExpireDuration int64
	JwtRefreshTokenPrivateKey     string
	JwtRefreshTokenPublicKey      string
}

// First Super User
type FirstSuperUserConfig struct {
	FirstSuperUserEmail    string
	FirstSuperUserName     string
	FirstSuperUserPassword string
}

type EmailConfig struct {
	EmailFrom                string
	EmailName                string
	EmailLink                string
	EmailLogoLink            string
	EmailCopyright           string
	EmailVerificationSubject string
	EmailResetSubject        string
}

type SmtpEmailConfig struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string
	SmtpUseTls   bool
	SmtpUseSsl   bool
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	cfg = &c

	return &c, nil
}

func GetCfg() *Config {
	if cfg == nil {
		cfg = new(Config)
	}
	return cfg
}
