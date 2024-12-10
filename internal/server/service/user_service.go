package service

import (
	"github.com/ShvetsovYura/oykeeper/internal/server/store"
	pb "github.com/ShvetsovYura/oykeeper/proto"
)

type UserService struct {
	pb.UnimplementedRecordServiceServer
	recordStore *store.UserStore
}

func NewUserService(store *store.UserStore) *UserService {
	return &UserService{
		recordStore: store,
	}
}
