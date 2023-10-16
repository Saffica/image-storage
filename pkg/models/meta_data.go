package models

import "time"

type MetaData struct {
	ID           int64
	DownloadLink string
	UpdatedAt    time.Time
	Downloaded   bool
}
