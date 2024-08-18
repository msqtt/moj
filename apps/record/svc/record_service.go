package svc

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/game/pkg/app_err"
	"moj/apps/record/db"
	red_pb "moj/apps/record/rpc"
	"moj/domain/pkg/queue"
	"moj/domain/record"
	"time"
)

// ModifyRecord implements red_pb.RecordServiceServer.
func (s *Server) ModifyRecord(ctx context.Context, req *red_pb.ModifyRecordRequest) (
	resp *red_pb.ModifyRecordResponse, err error) {
	slog.Debug("modify record request", "req", req)

	cmd := record.ModifyRecordCmd{
		RecordID:       req.RecordID,
		JudgeStatus:    req.JudgeStatus,
		FailedReason:   req.GetFailedReason(),
		MemoryUsed:     int(req.MemoryUsed),
		TimeUsed:       int(req.TimeUsed),
		CPUTimeUsed:    int(req.CpuTimeUsed),
		NumberFinishAt: int(req.NumberFinishAt),
		TotalQuestion:  int(req.TotalQuestion),
		Time:           time.Now().Unix(),
	}

	slog.Info("start to invoke modify record command", "cmd", cmd)
	err = s.commandInvoker.Invoke(ctx, func(ctx1 context.Context, eq queue.EventQueue) (any, error) {
		return s.modifyRecordCmdHandler.Handle(ctx1, eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke modify record command", "err", err)
		err = responseStatusError(err)
	}
	resp = &red_pb.ModifyRecordResponse{
		Time: time.Now().Unix(),
	}
	return
}

// SubmitRecord implements red_pb.RecordServiceServer.
func (s *Server) SubmitRecord(ctx context.Context, req *red_pb.SubmitRecordRequest) (
	resp *red_pb.SubmitRecordResponse, err error) {
	slog.Debug("submit record request", "req", req)

	cmd := record.SubmitRecordCmd{
		AccountID:  req.AccountID,
		GameID:     req.GetGameID(),
		QuestionID: req.QuestionID,
		Language:   req.Language,
		Code:       req.GetCode(),
		Time:       time.Now().Unix(),
	}

	slog.Info("start to invoke submit record command", "cmd", cmd)
	var id string
	err = s.commandInvoker.Invoke(ctx, func(ctx1 context.Context, eq queue.EventQueue) (any, error) {
		id1, err := s.submitRecordCmdHandler.Handle(ctx1, eq, cmd)
		id = id1
		return id1, err
	})
	if err != nil {
		slog.Error("failed to invoke submit record command", "err", err)
		err = responseStatusError(err)
	}

	resp = &red_pb.SubmitRecordResponse{
		RecordID: id,
		Time:     time.Now().Unix(),
	}
	return
}

func fromRecord(rv []*db.RecordModel) []*red_pb.Record {
	ret := make([]*red_pb.Record, len(rv))
	for i, v := range rv {
		ret[i] = &red_pb.Record{
			RecordID:         v.ID.Hex(),
			AccountID:        v.AccountID,
			GameID:           &v.GameID,
			QuestionID:       v.QuestionID,
			Language:         v.Language,
			Code:             v.Code,
			CodeHash:         v.CodeHash,
			JudgeStatus:      v.JudgeStatus,
			FailedReason:     &v.FailedReason,
			NumberFinishedAt: int64(v.NumberFinisheAt),
			TotalQuestion:    int64(v.TotalQuestion),
			CreateTime:       v.CreateTime.Unix(),
			FinishTime:       v.FinishTime.Unix(),
			MemoryUsed:       int64(v.MemoryUsed),
			TimeUsed:         int64(v.TimeUsed),
			CpuTimeUsed:      int64(v.CPUTimeUsed),
		}
	}
	return ret
}

// GetRecordPage implements red_pb.RecordServiceServer.
func (s *Server) GetRecordPage(ctx context.Context, req *red_pb.GetRecordPageRequest) (
	resp *red_pb.GetRecordPageResponse, err error) {
	slog.Debug("get record page request", "req", req)

	m := make(map[string]any)

	res, total, err := s.recordViewDao.FindPage(ctx, req.QuestionID,
		req.AccountID, int(req.Page), int(req.PageSize), m)
	if err != nil {
		slog.Error("failed to get record page", "err", err)
		err = responseStatusError(err)
	}

	resp = &red_pb.GetRecordPageResponse{
		Records: fromRecord(res),
		Total:   total,
	}
	return
}

// GetRecord implements red_pb.RecordServiceServer.
func (s *Server) GetRecord(ctx context.Context, req *red_pb.GetRecordRequest) (
	resp *red_pb.GetRecordResponse, err error) {
	slog.Debug("get record info request", "req", req)
	rec, err := s.recordRepository.FindRecordByID(ctx, req.RecordID)
	if err != nil {
		slog.Error("failed to get record info", "err", err)
		err = responseStatusError(err)
	}

	resp = &red_pb.GetRecordResponse{
		Record: &red_pb.Record{
			RecordID:         rec.RecordID,
			AccountID:        rec.AccountID,
			GameID:           &rec.GameID,
			QuestionID:       rec.QuestionID,
			Language:         rec.Language,
			Code:             rec.Code,
			CodeHash:         rec.CodeHash,
			JudgeStatus:      rec.JudgeStatus,
			FailedReason:     &rec.FailedReason,
			NumberFinishedAt: int64(rec.NumberFinishedAt),
			TotalQuestion:    int64(rec.TotalQuestion),
			CreateTime:       rec.CreateTime,
			FinishTime:       rec.FinishTime,
			MemoryUsed:       int64(rec.MemoryUsed),
			TimeUsed:         int64(rec.TimeUsed),
			CpuTimeUsed:      int64(rec.CPUTimeUsed),
		},
	}
	return
}

// CheckAccountPass implements red_pb.RecordServiceServer.
func (s *Server) CheckAccountPass(ctx context.Context, req *red_pb.CheckAccountPassRequest) (
	resp *red_pb.CheckAccountPassResponse, err error) {
	slog.Debug("check account pass request", "req", req)
	view, err := s.passQuestionViewDao.FindByAccountIDAndQuestionID(ctx, req.AccountID, req.QuestionID)

	passStatus := red_pb.CheckAccountPassResponse_Working
	if err != nil {
		if !errors.Is(err, app_err.ErrModelNotFound) {
			slog.Error("failed to check account pass", "err", err)
			err = responseStatusError(err)
			return
		}
		passStatus = red_pb.CheckAccountPassResponse_Undo
	} else if view.Status == db.PassStatusPass {
		passStatus = red_pb.CheckAccountPassResponse_Pass
	}

	resp = &red_pb.CheckAccountPassResponse{
		Status: passStatus,
	}
	return
}

// GetAccountPassedCount implements red_pb.RecordServiceServer.
func (s *Server) GetAccountPassedCount(ctx context.Context, req *red_pb.GetAccountPassedCountRequest) (
	resp *red_pb.GetAccountPassedCountResponse, err error) {
	slog.Debug("get account passed count request", "req", req)

	eazy, normal, hard, err := s.passQuestionViewDao.CountLevelByAccountID(ctx, req.AccountID)
	if err != nil {
		slog.Error("failed to get account passed count", "err", err)
		err = responseStatusError(err)
	}

	resp = &red_pb.GetAccountPassedCountResponse{
		Eazy:   eazy,
		Normal: normal,
		Hard:   hard,
	}
	return
}

// GetDailyTaskView implements red_pb.RecordServiceServer.
func (s *Server) GetDailyTaskView(ctx context.Context, req *red_pb.GetDailyTaskViewRequest) (
	resp *red_pb.GetDailyTaskViewResponse, err error) {
	slog.Debug("get daily task view request", "req", req)
	date := time.Unix(req.Time, 0)
	taskView, err := s.dayTaskViewDao.FindByDate(ctx, date)
	if err != nil {
		slog.Error("failed to get daily task view", "err", err)
		err = responseStatusError(err)
	}
	resp = &red_pb.GetDailyTaskViewResponse{
		SubmitNumber: int64(taskView.SubmitNumber),
		FinishNumber: int64(taskView.FinishNumber),
		Time:         taskView.Time.Unix(),
	}
	return
}

// GetQuestionRecordCount implements red_pb.RecordServiceServer.
func (s *Server) GetQuestionRecordCount(ctx context.Context, req *red_pb.GetQuestionRecordCountRequest) (
	resp *red_pb.GetQuestionRecordCountResponse, err error) {
	// slog.Debug("get question record count request", "req", req)

	total, err1 := s.recordViewDao.CountAllByID(ctx, req.QuestionID, req.GetGameID())
	if err1 != nil {
		slog.Error("failed to get question record count", "err", err1)
		err = responseStatusError(err1)
	}

	passCount, err1 := s.passQuestionViewDao.CountByQuestionID(ctx, req.QuestionID)
	if err1 != nil {
		slog.Error("failed to get question record count", "err", err1)
		err = responseStatusError(err1)
	}

	resp = &red_pb.GetQuestionRecordCountResponse{
		PassedCount: passCount,
		SubmitTotal: total,
	}
	return
}

var _ red_pb.RecordServiceServer = (*Server)(nil)
