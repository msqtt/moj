package record

type RecordRepository interface {
	FindRecordByID(recordID int) (*Record, error)
	Save(*Record) error
}
