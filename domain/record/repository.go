package record

type RecordRepository interface {
	FindRecordByID(recordID string) (*Record, error)
	FindBestGameRecord(gameID, accountID string) (*Record, error)
	Save(*Record) error
}
