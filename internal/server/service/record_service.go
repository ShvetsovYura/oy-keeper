package service

import (
	"context"
	"log/slog"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	"github.com/ShvetsovYura/oykeeper/internal/server/store"
	"github.com/ShvetsovYura/oykeeper/internal/utils"
	pb "github.com/ShvetsovYura/oykeeper/proto"
	"github.com/jackc/pgerrcode"
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

func (r *RecordService) GetUserRecords(ctx context.Context, in *pb.UserRecordsReq) (*pb.UserRecordsResp, error) {
	rr, _ := r.recordStore.FetchUserRecords(ctx, in.Uuid)
	aa, _ := r.recordStore.FetchUserAttributes(ctx, in.Uuid)
	ff, _ := r.recordStore.FetchUserFilesInfos(ctx, in.Uuid)
	var response = pb.UserRecordsResp{Uuid: in.Uuid, Records: []*pb.RecordInfo{}}

	for _, r := range rr {
		filteredAttrs := utils.Filter(aa, func(item store.AttributeDB) bool {
			return item.ItemUuid == r.Uuid
		})
		attrs := utils.Map(filteredAttrs, func(item store.AttributeDB) *pb.AttributeInfo {
			return &pb.AttributeInfo{
				Uuid:  item.Uuid,
				Name:  item.Name,
				Value: item.Value,
			}
		})
		filteredFilesInfo := utils.Filter(ff, func(item store.FileDB) bool {
			return item.ItemUuid == r.Uuid
		})
		files := utils.Map(filteredFilesInfo, func(f store.FileDB) *pb.FileInfo {
			return &pb.FileInfo{
				Uuid: f.Uuid,
				Path: f.Path,
				Hash: f.Hash,
				Size: uint64(f.Size),
				Meta: f.Meta,
			}
		})

		ri := &pb.RecordInfo{
			Uuid:     r.Uuid,
			Name:     r.Name,
			Username: r.Username,
			Password: r.Password,
			Url:      r.Url,
			// ExpiredAt:   &expiredAt,
			Description: r.Description,
			Version:     int32(r.Version),
			Attributes:  attrs,
			Files:       files,
		}
		response.Records = append(response.Records, ri)
	}

	logger.Log.Info("res", slog.Any("rec", rr), slog.Any("attrs", aa), slog.Any("files", ff))
	return &response, nil
}

// создание записи
func (r *RecordService) CreateRecord(ctx context.Context, in *pb.RecordReq) (*pb.RecordResp, error) {
	_, err := r.recordStore.GetRecordVersion(ctx, in.Uuid, uint32(in.Version))
	if err != nil {
		logger.Log.Debug("error on get record version", slog.String("msg", err.Error()))
		if err.Error() != "sql: no rows in result set" {
			return nil, err
		}

	}
	err = r.recordStore.UpsertRecord(ctx, in, in.Attributes, in.Files)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	return &pb.RecordResp{}, nil
}

func (r *RecordService) ChangeRecord(ctx context.Context, in *pb.RecordReq) (*pb.RecordResp, error) {
	_, err := r.recordStore.GetRecordVersion(ctx, in.Uuid, uint32(in.Version))
	if err != nil {
		if !pgerrcode.IsNoData(err.Error()) {
			return nil, err
		}

	}
	err = r.recordStore.UpsertRecord(ctx, in, in.Attributes, in.Files)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	return &pb.RecordResp{}, nil
}

func (r *RecordService) RemoveRecord(ctx context.Context, in *pb.RecordRemoveReq) (*pb.RecordRemoveResp, error) {
	return nil, nil
}
