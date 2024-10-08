package etc

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort                    int    `mapstructure:"APP_PORT"`
	MongoHost                  string `mapstructure:"MONGO_HOST"`
	MongoDBName                string `mapstructure:"MONGO_DBNAME"`
	KeyFile                    string `mapstructure:"KEY_FILE"`
	CertFile                   string `mapstructure:"CERT_FILE"`
	TLS                        bool
	Debug                      bool
	NsqdAddr                   string        `mapstructure:"NSQD_ADDR"`
	NsqLookUpAddr              string        `mapstructure:"NSQ_LOOKUP_ADDR"`
	QuestionRPCAddr            string        `mapstructure:"QUESTION_RPC_ADDR"`
	ExecuteJudgeTopic          string        `mapstructure:"EXECUTE_JUDGE_TOPIC"`
	FinishRecordTopic          string        `mapstructure:"FINISH_RECORD_TOPIC"`
	FinishRecordTopicChannel   string        `mapstructure:"FINISH_RECORD_TOPIC_CHANNEL"`
	CalculateScoreTopic        string        `mapstructure:"CALCULATE_SCORE_TOPIC"`
	ScheduleRedoSubmitDuration time.Duration `mapstructure:"SCHEDULE_REDO_SUBMIT_DURATION"`
}

func NewAppConfig() *Config {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/moj")

	viper.SetConfigType("env")
	viper.SetConfigName("record")
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
