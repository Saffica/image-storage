package main

import (
	"context"
	"fmt"

	"github.com/Saffica/image-storage/pkg/db"
	"github.com/Saffica/image-storage/pkg/models"
	"github.com/Saffica/image-storage/pkg/repository/metadata"
)

// "log"
// "os"

// "github.com/Saffica/image-storage/internal/server"
// "github.com/Saffica/image-storage/internal/usecase"
// "github.com/Saffica/image-storage/pkg/restclient"
// "github.com/gin-gonic/gin"
// "github.com/minio/minio-go"

func main() {
	db, err := db.New(
		"postgres",
		"postgres",
		"localhost",
		"img_db",
		"./migrations/db",
		6433,
	)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	metadataRepository := metadata.New(db)
	metadata, err := metadataRepository.Get(ctx, "test")
	if err != nil {
		panic(err)
	}

	fmt.Println(*metadata)

	sMetadata := &models.MetaData{
		DownloadLink: "s",
		Downloaded:   true,
	}

	md, err := metadataRepository.Insert(ctx, sMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Println(md)

	md.DownloadLink = "updated_link"
	md.Downloaded = false
	md2, err := metadataRepository.Update(ctx, md)
	fmt.Println(md2)
	db.Close()
	// minioEndpoint := "localhost:9000"
	// accessKeyID := os.Getenv("MINIO_SERVER_ACCESS_KEY")
	// secretAccessKey := os.Getenv("MINIO_SERVER_SECRET_KEY")
	// useSSL := false

	// minioClient, err := minio.New(minioEndpoint, accessKeyID, secretAccessKey, useSSL)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// minioClient.PutObject()
	// log.Printf("%#v\n", minioClient)
	// client := restclient.New()
	// imgService := usecase.New(client)
	// router := gin.Default()

	// s := server.New(imgService, router)
	// err := s.Run(8080)
	// if err != nil {
	// 	log.Panic(err)
	// }
}
