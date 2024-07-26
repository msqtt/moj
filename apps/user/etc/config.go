package etc

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppPort             int    `mapstructure:"APP_PORT"`
	MongoHost           string `mapstructure:"MONGO_HOST"`
	DatabaseName        string `mapstructure:"DATABASE_NAME"`
	KeyFile             string `mapstructure:"KEY_FILE"`
	CertFile            string `mapstructure:"CERT_FILE"`
	CaptchaLiveDuration int64  `mapstructure:"CAPTCHA_LIVE_DURATION"`
	TLS                 bool
	Debug               bool
	EmailFromAddr       string `mapstructure:"EMAIL_FROM_ADDR"`
	EmailSmtpHost       string `mapstructure:"EMAIL_SMTP_HOST"`
	EmailSmtpPort       int    `mapstructure:"EMAIL_SMTP_PORT"`
	EmailsmtpScrt       string `mapstructure:"EMAIL_SMTP_SCRT"`
}

func NewAppConfig() *Config {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/moj")

	viper.SetConfigType("env")
	viper.SetConfigName("user")
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
