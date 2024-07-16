package record

type RecordRepository interface {
	FindRecordByID(recordID int) (*Record, error)
	FindBestGameRecord(gameID, accountID int) (*Record, error)
	Save(*Record) error
}
