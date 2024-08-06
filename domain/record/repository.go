package record

type RecordRepository interface {
	FindRecordByID(recordID string) (*Record, error)
	Save(*Record) (string, error)
}
