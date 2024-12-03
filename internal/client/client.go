package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	pb "github.com/ShvetsovYura/oykeeper/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	recordUuid, _ := uuid.NewRandom()
	userUuid, _ := uuid.NewRandom()
	fileUuid, _ := uuid.NewRandom()
	attr1Uuid, _ := uuid.NewRandom()
	attr2Uuid, _ := uuid.NewRandom()
	resp, err := c.recordService.CreateRecord(ctx, &pb.RecordReq{
		Uuid:     recordUuid.String(),
		UserUuid: userUuid.String(),
		Version:  "1",
		Attributes: []*pb.AttributeInfo{
			{
				Uuid:  attr1Uuid.String(),
				Name:  "attr1",
				Value: "value1",
			}, {
				Uuid:  attr2Uuid.String(),
				Name:  "attr2",
				Value: "value2",
			},
		},
		Files: []*pb.FileInfo{{
			Name: "file1",
			Uuid: fileUuid.String(),
			Path: "",
			Hash: "hash",
			Size: 1234,
		}},
	})
	if err != nil {
		fmt.Printf("error: %s \n", err.Error())
	}
	fmt.Printf("resp: %v\n", resp)
	return nil
}

func (c *Client) Upload(ctx context.Context) error {
	logger.Log.Debug("start upload file")
	stream, err := c.uploadService.Upload(ctx)
	if err != nil {
		return err
	}
	file, err := os.Open("cmds")
	if err != nil {
		return err
	}
	buf := make([]byte, 1024*1024)
	batchNumber := 1
	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunk := buf[:num]
		if err := stream.Send(&pb.FileUploadRequest{FileName: "cmds", Chunk: chunk}); err != nil {
			return err
		}
		logger.Log.Info("sent batch", slog.Int("batch_number", batchNumber), slog.Int("size", len(chunk)))
		batchNumber += 1
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	logger.Log.Info("sent done", slog.Int("size", int(res.GetSize())), slog.String("name", "hoho"))
	return nil
}
