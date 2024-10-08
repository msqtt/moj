package question

import (
	"errors"
	"fmt"
	domain_err "moj/domain/pkg/error"
)

type QuestionLevel int

const (
	QuestionLevelEasy QuestionLevel = iota
	QuestionLevelNormal
	QuestionLevelHard
)

func (q QuestionLevel) IsValid() bool {
	return q >= QuestionLevelEasy && q <= QuestionLevelHard
}

func (q QuestionLevel) String() string {
	return []string{"Easy", "Normal", "Hard"}[q]
}

func FromStringLevel(level string) QuestionLevel {
	switch level {
	case "Easy":
		return QuestionLevelEasy
	case "Normal":
		return QuestionLevelNormal
	case "Hard":
		return QuestionLevelHard
	default:
		return QuestionLevelNormal
	}
}

type QuestionLanguage string

const (
	QuestionLangC      QuestionLanguage = "c"
	QuestionLangCpp    QuestionLanguage = "cpp"
	QuestionLangGo     QuestionLanguage = "golang"
	QuestionLangJava   QuestionLanguage = "java"
	QuestionLangPython QuestionLanguage = "python"
	QuestionLangRust   QuestionLanguage = "rust"
)

func (q QuestionLanguage) IsValid() bool {
	return q == QuestionLangC || q == QuestionLangCpp || q == QuestionLangGo ||
		q == QuestionLangJava || q == QuestionLangPython || q == QuestionLangRust
}

var (
	ErrQuestionNotFound        = errors.New("question not found")
	ErrInValidQuestionLanguage = errors.New("invalid language")
	ErrInvalidQuestionLevel    = errors.Join(domain_err.ErrInValided, errors.New("invalid level"))
	ErrEmpltyCases             = errors.Join(domain_err.ErrInValided, errors.New("empty cases"))
)

type Question struct {
	QuestionID       string
	AccountID        string
	Enabled          bool
	Title            string
	Content          string
	Level            QuestionLevel
	AllowedLanguages []QuestionLanguage
	Cases            []Case
	TimeLimit        int
	MemoryLimit      int
	Tags             []string
	CreateTime       int64
	ModifyTime       int64
}

func NewQuestion(questionID, accountID string, title, text string, level QuestionLevel,
	langs []QuestionLanguage, timeLimit, memoryLimit int,
	tags []string, createTime, modifyTime int64, cases []Case) (que *Question, err error) {

	if !level.IsValid() {
		err = ErrInvalidQuestionLevel
	}

	if err1 := ValidCases(cases); err1 != nil {
		err = errors.Join(err, err1)
	}

	for _, lang := range langs {
		if !lang.IsValid() {
			err = errors.Join(err, errors.Join(domain_err.ErrInValided,
				fmt.Errorf("%w: %s", ErrInValidQuestionLanguage, lang)))
		}
	}

	if err != nil {
		return nil, err
	}

	que = &Question{
		QuestionID:       questionID,
		AccountID:        accountID,
		Enabled:          false,
		Title:            title,
		Content:          text,
		Level:            level,
		AllowedLanguages: langs,
		TimeLimit:        timeLimit,
		MemoryLimit:      memoryLimit,
		Tags:             tags,
		CreateTime:       createTime,
		ModifyTime:       modifyTime,
		Cases:            cases,
	}
	return
}

func ValidCases(cases []Case) error {
	if len(cases) == 0 {
		return ErrEmpltyCases
	}
	return nil
}
