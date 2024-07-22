package policy

import (
	"github.com/msqtt/moj/domain/judgement"
	"github.com/msqtt/moj/domain/question"
)

type CaseFileService interface {
	Read(fileName string) (string, error)
	ReadAllCaseFile(caseFiles []question.Case) ([]judgement.Case, error)
}
