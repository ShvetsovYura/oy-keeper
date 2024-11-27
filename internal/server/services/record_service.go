package services

type RecordService struct {
	recordStore any
}

func NewRecordService(store any) *RecordService {
	return &RecordService{recordStore: store}
}

func (r *RecordService) CreateRecord(data any) error {
	return nil
}
