package policy

import (
	"context"
	"moj/domain/judgement"
	"moj/domain/question"
)

type CaseFileService interface {
	Read(ctx context.Context, fileName string) (string, error)
	ReadAllCaseFile(ctx context.Context, caseFiles []question.Case) ([]judgement.Case, error)
}
