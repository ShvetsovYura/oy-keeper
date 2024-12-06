package service

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	"github.com/ShvetsovYura/oykeeper/internal/utils"
	pb "github.com/ShvetsovYura/oykeeper/proto"
	"google.golang.org/grpc/metadata"
)

type FileService struct {
	pb.UnimplementedFileServiceServer
	path_dir string
}

func NewFileService(path string) *FileService {
	return &FileService{path_dir: path}
}

func (s *FileService) Upload(stream pb.FileService_UploadServer) error {
	var inHash string

	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		values := md.Get("hash")
		if len(values) > 0 {
			inHash = values[0]
		}
	}
	file := NewFile()
	var totalSize uint32
	defer func() {
		if err := file.OutputFile.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
	}()
	for {
		req, err := stream.Recv()
		if file.FilePath == "" {
			file.SetFile(req.GetFileName(), s.path_dir)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error %w", err)
		}

		chunk := req.GetChunk()
		totalSize += uint32(len(chunk))
		logger.Log.Info("received a chunk", slog.Int("size", int(totalSize)))
		if err := file.Write(chunk); err != nil {
			return err
		}
	}
	fileName := filepath.Base(file.FilePath)
	hash, err := utils.MD5Sum(file.FilePath)

	if hash != inHash || err != nil {
		s.removeFile(file.FilePath)
		return errors.New("Hashes is not equals or error hash")
	}
	logger.Log.Info("file received", slog.String("name", fileName), slog.String("hash", hash))
	return stream.SendAndClose(&pb.FileUploadResponse{})
}

func (s *FileService) removeFile(path string) error {
	return os.Remove(path)
}
