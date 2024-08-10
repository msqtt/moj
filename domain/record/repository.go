package record

import "context"

type RecordRepository interface {
	FindRecordByID(ctx context.Context, recordID string) (*Record, error)
	Save(context.Context, *Record) (string, error)
}
