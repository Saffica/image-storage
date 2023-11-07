package main

import (
	"github.com/Saffica/image-storage/pkg/db"
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
