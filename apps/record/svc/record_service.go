package svc

import (
	"context"
	"log/slog"
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

func fromRecordView(rv []*db.RecordViewModel) []*red_pb.RecordView {
	ret := make([]*red_pb.RecordView, len(rv))
	for i, v := range rv {
		ret[i] = &red_pb.RecordView{
			RecordID:         v.ID.Hex(),
			GameID:           &v.GameID,
			Language:         v.Language,
			CodeHash:         v.CodeHash,
			JudgeStatus:      v.JudgeStatus,
			NumberFinishedAt: int64(v.NumberFinisheAt),
			TotalQuestion:    int64(v.TotalQuestion),
			CreateTime:       v.CreateTime.Unix(),
			MemoryUsed:       int64(v.MemoryUsed),
			TimeUsed:         int64(v.TimeUsed),
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
		RecordView: fromRecordView(res),
		Total:      total,
	}
	return
}

// GetRecordInfo implements red_pb.RecordServiceServer.
func (s *Server) GetRecordInfo(ctx context.Context, req *red_pb.GetRecordInfoRequest) (
	resp *red_pb.GetRecordInfoResponse, err error) {
	slog.Debug("get record info request", "req", req)
	rec, err := s.recordRepository.FindRecordByID(ctx, req.RecordID)
	if err != nil {
		slog.Error("failed to get record info", "err", err)
		err = responseStatusError(err)
	}

	resp = &red_pb.GetRecordInfoResponse{
		RecordInfo: &red_pb.RecordInfo{
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

var _ red_pb.RecordServiceServer = (*Server)(nil)
