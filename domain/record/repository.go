package record

type RecordRepository interface {
	findRecordByID(recordID int) (*Record, error)
	save(*Record) error
}
