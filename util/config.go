package util

import (
	"time"

	"github.com/spf13/viper"
)

// / Config store all configuration of application
// / The values are read by viper from config file
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymetricKey     string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EMailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EMailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EMailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	CORSHTTP             string        `mapstructure:"CORSHTTP"`
	WebServer            string        `mapstructure:"WEB_SERVER"`
}

// Read configuration from env file
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
