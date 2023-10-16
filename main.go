// Сервис кеширования картинок
// - на вход приходит запрос вида /img/{hash}?w=100&h=50
// -в hash зашифрованный адрес изображения на внешнем ресурсе
// -сервис проверяет присутствует ли изображение в s3 если нет то скачивает с внешнего ресурса и кладет в s3 предварительно пережав его в webp в случае ошибки скачивания помечает что картинка отсутствует и не пытается скачать в течении суток после этой попытки
// -сервис берет изображение из s3 масштабирует его в зависимости от значений параметров w и h и отдает картинку в ответ на запрос
// -в случае отсутствия картики и не возможности скачать ее с внешнего ресурса отдает 404

package main

// "log"
// "os"

// "github.com/Saffica/image-storage/internal/server"
// "github.com/Saffica/image-storage/internal/usecase"
// "github.com/Saffica/image-storage/pkg/restclient"
// "github.com/gin-gonic/gin"
// "github.com/minio/minio-go"

func main() {
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
