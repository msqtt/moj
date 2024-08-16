package db

import (
	"moj/domain/question"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Case struct {
	Number         int
	InputFilePath  string
	OutputFilePath string
}

type QuestionModel struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	AccountID        string             `bson:"account_id"`
	Enabled          bool
	Title            string
	Content          string
	Level            string
	AllowedLanguages []string `bson:"allowed_languages"`
	Cases            []Case
	Tags             []string
	TimeLimit        int       `bson:"time_limit"`
	MemoryLimit      int       `bson:"memory_limit"`
	CreateTime       time.Time `bson:"create_time"`
	ModifyTime       time.Time `bson:"modify_time"`
}

func NewFromAggreation(q *question.Question) *QuestionModel {
	id, _ := primitive.ObjectIDFromHex(q.QuestionID)
	langs := FromQuestionLangs(q.AllowedLanguages)
	cases := FromQuestionCases(q.Cases)
	return &QuestionModel{
		ID:               id,
		AccountID:        q.AccountID,
		Enabled:          q.Enabled,
		Title:            q.Title,
		Level:            q.Level.String(),
		AllowedLanguages: langs,
		Cases:            cases,
		TimeLimit:        q.TimeLimit,
		MemoryLimit:      q.MemoryLimit,
		Tags:             q.Tags,
		CreateTime:       time.Unix(q.CreateTime, 0),
		ModifyTime:       time.Unix(q.ModifyTime, 0),
	}
}

func FromQuestionLangs(allowedLanguages []question.QuestionLanguage) []string {
	langs := make([]string, len(allowedLanguages))
	for idx, lang := range allowedLanguages {
		langs[idx] = string(lang)
	}
	return langs
}

func FromQuestionCases(cases []question.Case) []Case {
	ret := make([]Case, len(cases))
	for idx, cas := range cases {
		ret[idx] = Case{
			Number:         cas.Number,
			InputFilePath:  cas.InputFilePath,
			OutputFilePath: cas.OutputFilePath,
		}
	}
	return ret
}

func ToQuestionLangs(allowedLanguages []string) []question.QuestionLanguage {
	langs := make([]question.QuestionLanguage, len(allowedLanguages))
	for idx, lang := range allowedLanguages {
		langs[idx] = question.QuestionLanguage(lang)
	}
	return langs
}

func ToQuestionCases(cases []Case) []question.Case {
	ret := make([]question.Case, len(cases))
	for idx, cas := range cases {
		ret[idx] = question.Case{
			Number:         cas.Number,
			InputFilePath:  cas.InputFilePath,
			OutputFilePath: cas.OutputFilePath,
		}
	}
	return ret
}

func (q *QuestionModel) ToAggregate() *question.Question {
	langs := ToQuestionLangs(q.AllowedLanguages)
	cases := ToQuestionCases(q.Cases)
	return &question.Question{
		QuestionID:       q.ID.Hex(),
		AccountID:        q.AccountID,
		Enabled:          q.Enabled,
		Title:            q.Title,
		Content:          q.Content,
		Level:            question.FromStringLevel(q.Level),
		AllowedLanguages: langs,
		Cases:            cases,
		TimeLimit:        q.TimeLimit,
		MemoryLimit:      q.MemoryLimit,
		Tags:             q.Tags,
		CreateTime:       q.CreateTime.Unix(),
		ModifyTime:       q.ModifyTime.Unix(),
	}
}
