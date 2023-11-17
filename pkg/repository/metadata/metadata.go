package metadata

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Saffica/image-storage/pkg/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type dataBase interface {
	Close()
	BeginTx(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, query string, arg ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, arg ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, arg ...interface{}) pgx.Row
}

type metaDataRepository struct {
	db dataBase
}

func New(db dataBase) *metaDataRepository {
	return &metaDataRepository{
		db: db,
	}
}

func (r *metaDataRepository) Get(ctx context.Context, url string) (*models.MetaData, error) {
	metadata := &models.MetaData{}
	row := r.db.QueryRow(
		ctx,
		"SELECT metadata_id, download_link, downloaded, updated_at FROM metadata WHERE download_link = $1",
		url,
	)
	err := row.Scan(&metadata.ID, &metadata.DownloadLink, &metadata.Downloaded, &metadata.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", models.ErrMetaDataNotFound, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (r *metaDataRepository) Insert(ctx context.Context, metadata *models.MetaData) (*models.MetaData, error) {
	row := r.db.QueryRow(
		ctx,
		"INSERT INTO metadata (download_link, downloaded) VALUES ($1, $2) RETURNING metadata_id, updated_at",
		metadata.DownloadLink,
		metadata.Downloaded,
	)

	err := row.Scan(&metadata.ID, &metadata.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (r *metaDataRepository) Update(ctx context.Context, metadata *models.MetaData) (*models.MetaData, error) {
	now := time.Now()
	c, err := r.db.Exec(
		ctx,
		`UPDATE metadata
		SET download_link=$1, downloaded=$2, updated_at=$3
		WHERE metadata_id=$4`,
		metadata.DownloadLink,
		metadata.Downloaded,
		now,
		metadata.ID,
	)

	if err != nil {
		return nil, err
	}

	if c.RowsAffected() == 0 {
		return nil, models.ErrMetaDataNotFound
	}

	metadata.UpdatedAt = now
	return metadata, nil
}
