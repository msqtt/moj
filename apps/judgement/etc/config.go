package etc

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppPort                  int    `mapstructure:"APP_PORT"`
	MongoHost                string `mapstructure:"MONGO_HOST"`
	MongoDBName              string `mapstructure:"MonGO_DBNAME"`
	KeyFile                  string `mapstructure:"KEY_FILE"`
	CertFile                 string `mapstructure:"CERT_FILE"`
	TLS                      bool
	Debug                    bool
	SbJudgerRPCAddr          string `mapstructure:"SB_JUDGER_RPC_ADDR"`
	QuestionRPCAddr          string `mapstructure:"QUESTION_RPC_ADDR"`
	OutPutMsgLimit           int    `mapstructure:"OUTPUT_MSG_LIMIT"`
	NsqdAddr                 string `mapstructure:"NSQD_ADDR"`
	NsqLookUpAddr            string `mapstructure:"NSQ_LOOKUP_ADDR"`
	ExecuteJudgeTopic        string `mapstructure:"EXECUTE_JUDGE_TOPIC"`
	ExecuteJudgeTopicChannel string `mapstructure:"EXECUTE_JUDGE_TOPIC_CHANNEL"`
	FinishJudgementTopic     string `mapstructure:"FINISH_RECORD_TOPIC"`
}

func NewAppConfig() *Config {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/moj")

	viper.SetConfigType("env")
	viper.SetConfigName("judgement")
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
