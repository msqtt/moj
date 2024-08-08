package domain

import (
	"moj/domain/judgement"
	"moj/domain/policy"
	"moj/domain/question"
)

type MinioCaseReader struct {
}

// Read implements policy.CaseFileService.
func (m *MinioCaseReader) Read(fileName string) (string, error) {
	panic("unimplemented")
}

// ReadAllCaseFile implements policy.CaseFileService.
func (m *MinioCaseReader) ReadAllCaseFile(caseFiles []question.Case) ([]judgement.Case, error) {
	panic("unimplemented")
}

func NewMinioCaseReader() policy.CaseFileService {
	return &MinioCaseReader{}
}
