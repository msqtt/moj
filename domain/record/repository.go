package record

import "context"

type RecordRepository interface {
	FindRecordByID(ctx context.Context, recordID string) (*Record, error)
	FindBestRecord(ctx context.Context, uid, qid, gid string) (*Record, error)
	Save(context.Context, *Record) (string, error)
}
