package client

import (
	"context"
	"io"
	"log/slog"
	"os"

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

	// resp, err := c.recordService.CreateRecord(ctx, &pb.RecordReq{
	// 	Uuid:     "d9651da1-6945-43ba-a730-544c5d17ad4f",
	// 	UserUuid: "21bd249e-77c7-4bf0-b78f-dd33d38ac54f",
	// 	Version:  2,
	// 	Attributes: []*pb.AttributeInfo{
	// 		{
	// 			Uuid:  "59fb3cd6-64c3-4f9c-ba74-b8bc24b8a57e",
	// 			Name:  "attr1",
	// 			Value: "value1",
	// 		}, {
	// 			Uuid:  "b7c72694-acd0-43cb-b57a-f43c147e6620",
	// 			Name:  "attr2",
	// 			Value: "value2",
	// 		},
	// 	},
	// 	Files: []*pb.FileInfo{{
	// 		Name: "file1",
	// 		Uuid: "ef4ad74b-c7ad-4268-a6e8-6086467721bf",
	// 		Path: *"fs",
	// 		Hash: "hash",
	// 		Size: 1234,
	// 	}},
	// })
	// if err != nil {
	// 	fmt.Printf("error: %s \n", err.Error())
	// }
	// fmt.Printf("resp: %v\n", resp)
	return nil
}

func (c *Client) Upload(ctx context.Context) error {
	logger.Log.Debug("start upload file")
	hash, err := utils.MD5Sum("cmds")
	if err != nil {
		return err
	}
	md := metadata.New(map[string]string{"hash": hash})
	outCtx := metadata.NewOutgoingContext(ctx, md)
	stream, err := c.uploadService.Upload(outCtx)
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
	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}
	return nil
}
