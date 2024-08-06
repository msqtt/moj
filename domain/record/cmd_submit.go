package record

import "moj/domain/pkg/queue"

type SubmitRecordCmd struct {
	AccountID  string
	GameID     string
	QuestionID string
	Language   string
	Code       string
	Time       int64
}

type SubmitRecordCmdHandler struct {
	repo RecordRepository
}

func NewSubmitRecordCmdHandler(repo RecordRepository) *SubmitRecordCmdHandler {
	return &SubmitRecordCmdHandler{
		repo: repo,
	}
}

func (s *SubmitRecordCmdHandler) Handle(queue queue.EventQueue, cmd SubmitRecordCmd) (string, error) {
	rec := NewRecord(cmd.AccountID, cmd.GameID, cmd.QuestionID, cmd.Language, cmd.Code, cmd.Time)
	id, err := s.repo.Save(rec)
	if err != nil {
		return "", err
	}
	return id, rec.submit(queue)
}
