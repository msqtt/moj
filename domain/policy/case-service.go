package policy

import (
	"moj/domain/judgement"
	"moj/domain/question"
)

type CaseFileService interface {
	Read(fileName string) (string, error)
	ReadAllCaseFile(caseFiles []question.Case) ([]judgement.Case, error)
}
