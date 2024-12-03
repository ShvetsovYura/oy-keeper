package service

import (
	"context"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	"github.com/ShvetsovYura/oykeeper/internal/server/store"
	pb "github.com/ShvetsovYura/oykeeper/proto"
)

type RecordService struct {
	pb.UnimplementedRecordServiceServer
	recordStore *store.RecordStore
}

func NewRecordService(store *store.RecordStore) *RecordService {
	if store == nil {
		panic("nil store")
	}
	return &RecordService{recordStore: store}
}

func (r *RecordService) CreateRecord(ctx context.Context, in *pb.RecordReq) (*pb.RecordResp, error) {
	// var resp *pb.RecordResp
	err := r.recordStore.NewRecord(ctx, in, in.Attributes, in.Files)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	return &pb.RecordResp{}, nil
}
