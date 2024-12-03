package service

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	pb "github.com/ShvetsovYura/oykeeper/proto"
)

type FileService struct {
	pb.UnimplementedFileServiceServer
}

func NewFileUploadHandler() *FileService {
	return &FileService{}
}

func (h *FileService) Upload(stream pb.FileService_UploadServer) error {
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
			file.SetFile(req.GetFileName(), "hoho")
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
	logger.Log.Info("file received", slog.String("name", fileName))
	return stream.SendAndClose(&pb.FileUploadResponse{Hash: "meme", Size: uint32(totalSize)})
}
