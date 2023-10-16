package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Saffica/image-storage/pkg/models"
	"github.com/gin-gonic/gin"
)

const (
	imgURL = "/img/:hash"
)

type imgServiceI interface {
	GetImgByHash(url string) ([]byte, error)
}

type handler struct {
	imgService imgServiceI
	router     *gin.Engine
	server     *http.Server
}

func New(imgService imgServiceI, router *gin.Engine) *handler {
	return &handler{
		imgService: imgService,
		router:     router,
	}
}

func (h *handler) Run(port int) error {
	h.router.GET(imgURL, h.getImg)
	h.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: h.router,
	}

	err := h.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (h *handler) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := h.server.Shutdown(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	<-ctx.Done()
}

func (h *handler) getImg(c *gin.Context) {
	// width := c.Query("w")
	// height := c.Query("h")
	byteFile, err := h.imgService.GetImgByHash(c.Param("hash"))
	if errors.Is(err, models.ErrBadHash) {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if errors.Is(err, models.ErrImageNotFound) {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println(byteFile)
	c.Header("Content-Disposition", "attachment; filename=output.webp")
	c.Data(http.StatusOK, "application/octet-stream", byteFile)
}
