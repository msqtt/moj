package domain

import (
	"context"
	"moj/domain/judgement"
	"moj/domain/policy"
	"moj/domain/question"
)

type MinioCaseReader struct {
}

// Read implements policy.CaseFileService.
func (m *MinioCaseReader) Read(ctx context.Context, fileName string) (string, error) {
	panic("unimplemented")
}

// ReadAllCaseFile implements policy.CaseFileService.
func (m *MinioCaseReader) ReadAllCaseFile(ctx context.Context, caseFiles []question.Case) ([]judgement.Case, error) {
	panic("unimplemented")
}

func NewMinioCaseReader() policy.CaseFileService {
	return &MinioCaseReader{}
}
