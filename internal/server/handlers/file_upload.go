package handlers

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	pb "github.com/ShvetsovYura/oy-keeper/proto"
)

type FileUploadHandler struct {
	pb.UnimplementedFileServiceServer
}

func NewFileUploadHandler() *FileUploadHandler {
	return &FileUploadHandler{}
}

func (h *FileUploadHandler) Upload(stream pb.FileService_UploadServer) error {
	file := NewFile()
	var fileSize uint32 = 0
	var buf *bytes.Buffer = &bytes.Buffer{}

	for {
		req, err := stream.Recv()
		if file.FilePath == "" {
			file.SetFile(req.GetFileName(), g.cfg.FilesStorage.Location)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error %w", err)
		}

		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
	}
	fileName := filepath.Base(file.FilePath)
	return stream.SendAndClose(&pb.FileUploadResponse{Hassh: "meme", Size: uint32(fileSize)})
}
