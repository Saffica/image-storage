package models

type ImageRequest struct {
	Hash   string
	Width  uint16 `form:"w"`
	Height uint16 `form:"h"`
}
