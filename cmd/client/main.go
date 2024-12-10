package main

import (
	"context"
	"fmt"

	"github.com/ShvetsovYura/oykeeper/internal/client"
	"github.com/ShvetsovYura/oykeeper/internal/logger"
	"github.com/fabiokaelin/terminalimage"
)

func main() {
	logger.Init()
	c := client.New(":9091")
	logger.Log.Info("startging up client app")
	c.Register(context.TODO(), "pipa")
	// f, err := filepath.Abs("mybin.zip")
	imageString, err := terminalimage.ImageToString("qr-code.png", 20, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(imageString)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// c.CreateRecord(context.Background())
	// c.Upload(context.Background(), f)
}
