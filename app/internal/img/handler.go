package img

import (
	"fmt"
	"net/http"

	"github.com/Saffica/image-storage/app/internal/handlers"
	"github.com/gin-gonic/gin"
)

const (
	imgURL = "/img/:hash"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *gin.Engine) {
	router.GET(imgURL, h.GetImg)
}

func (h *handler) GetImg(c *gin.Context) {
	hash := c.Param("hash")
	width := c.Query("w")
	height := c.Query("h")

	// c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s"))
	c.Header("Conten-type", c.ContentType())
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("hash: %s, w: %s, h: %s, contentType: %s", hash, width, height, c.ContentType())})

}
