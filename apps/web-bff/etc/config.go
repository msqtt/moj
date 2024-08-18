package etc

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort              int    `mapstructure:"APP_PORT"`
	KeyFile              string `mapstructure:"KEY_FILE"`
	CertFile             string `mapstructure:"CERT_FILE"`
	TLS                  bool
	Debug                bool
	AvatarFileSizeLimit  int64         `mapstructure:"AVATAR_FILE_SIZE_LIMIT"`
	AvatarFilePrefixPath string        `mapstructure:"AVATAR_FILE_PREFIX_PATH"`
	CaseFileSizeLimit    int64         `mapstructure:"CASE_FILE_SIZE_LIMIT"`
	CaseFilePrefixPath   string        `mapstructure:"CASE_FILE_PREFIX_PATH"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SymmetricKey         string        `mapstructure:"SYMMETRIC_KEY"`
	UserRPCAddr          string        `mapstructure:"USER_RPC_ADDR"`
	CaptchaRPCAddr       string        `mapstructure:"CAPTCHA_RPC_ADDR"`
	QuestionRPCAddr      string        `mapstructure:"QUESTION_RPC_ADDR"`
	GameRPCAddr          string        `mapstructure:"GAME_RPC_ADDR"`
	RecordRPCAddr        string        `mapstructure:"RECORD_RPC_ADDR"`
}

func NewAppConfig() *Config {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/moj")

	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	conf := &Config{}
	if err := viper.Unmarshal(conf); err != nil {
		panic(err)
	}
	return conf
}
