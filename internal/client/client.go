package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	"github.com/ShvetsovYura/oykeeper/internal/utils"
	pb "github.com/ShvetsovYura/oykeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn          *grpc.ClientConn
	uploadService pb.FileServiceClient
	recordService pb.RecordServiceClient
}

func New(addr string) *Client {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &Client{
		conn:          conn,
		uploadService: pb.NewFileServiceClient(conn),
		recordService: pb.NewRecordServiceClient(conn),
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateRecord(ctx context.Context) error {

	recordUuid := "44eb365b-3a88-4bed-a59b-bce4ceffcde6"
	req := &pb.RecordReq{
		Uuid:        recordUuid,
		Name:        "",
		Username:    new(string),
		Password:    new(string),
		Url:         new(string),
		ExpiredAt:   new(string),
		UserUuid:    "21bd249e-77c7-4bf0-b78f-dd33d38ac54f",
		CardNum:     new(string),
		Description: new(string),
		Version:     2,
		Attributes:  []*pb.AttributeInfo{{Uuid: "6efe0930-d7aa-42b2-bc54-300655597fd5", Name: "mycoffe", Value: "espresso"}, {Uuid: "64a4b557-61f4-49a3-b4be-dc7ee5eda151", Name: "branch", Value: "toast"}},
		// Files:       [] //[]*pb.FileInfoReq{{Path: "/Users/21184534/Documents/practicum/oykeeper/cmd/client/mybin.zip"}},
	}
	paths := []string{"/Users/21184534/Documents/practicum/oykeeper/cmd/client/mybin.zip"}

	for _, p := range paths {
		ur, err := c.Upload(ctx, p, recordUuid)
		if err == nil {
			req.Files = append(req.Files, &pb.FileInfo{
				Uuid: "46628b5b-c2c7-4575-8d22-154c482e2369",
				Hash: ur.Hash,
				Path: &ur.Path,
				Size: uint64(ur.Size),
				Name: "hohoh",
			})
		}
	}

	resp, err := c.recordService.CreateRecord(ctx, req)
	if err != nil {
		fmt.Printf("error: %s \n", err.Error())
	}

	fmt.Printf("resp: %v\n", resp)
	return nil
}

func (c *Client) Upload(ctx context.Context, path string, recordUuid string) (*pb.FileUploadResponse, error) {
	logger.Log.Debug("start upload file", slog.String("path", path))
	hash, err := utils.MD5Sum(path)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	md := metadata.New(map[string]string{"hash": hash})
	outCtx := metadata.NewOutgoingContext(ctx, md)
	stream, err := c.uploadService.Upload(outCtx)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		file.Close()
	}()
	filename := filepath.Base(path)
	buf := make([]byte, 1024*1024)
	batchNumber := 1
	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		chunk := buf[:num]

		if err := stream.Send(&pb.FileUploadRequest{FileName: filename, RecordUuid: recordUuid, Chunk: chunk}); err != nil {
			return nil, err
		}
		logger.Log.Info("sent batch", slog.Int("batch_number", batchNumber), slog.Int("size", len(chunk)))
		batchNumber += 1
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
