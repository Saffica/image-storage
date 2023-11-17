package metadata

import (
	"context"
	"testing"
	"time"

	"github.com/Saffica/image-storage/pkg/db"
	"github.com/Saffica/image-storage/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	const downloadLink = "test"
	const downloadLinkAfterUpdate = "test update ok"
	var metadataID int64

	testDB, err := db.New(
		"postgres",
		"postgres",
		"localhost",
		"img_db",
		"../../../migrations/db",
		6433,
	)

	require.NoError(t, err)

	defer testDB.Close()
	defer func() {
		_, err := testDB.Exec(
			context.Background(),
			"DELETE FROM metadata WHERE download_link = ANY($1)",
			[]string{downloadLink, downloadLinkAfterUpdate},
		)
		require.NoError(t, err)
	}()

	repo := New(testDB)

	t.Run("Insert OK", func(t *testing.T) {
		inputMetadata := &models.MetaData{
			DownloadLink: downloadLink,
			Downloaded:   true,
		}
		resultMetadata, err := repo.Insert(context.Background(), inputMetadata)
		require.NoError(t, err)
		require.NotEmpty(t, resultMetadata.ID)
		require.NotEmpty(t, resultMetadata.UpdatedAt)
		require.Equal(t, inputMetadata.DownloadLink, resultMetadata.DownloadLink)
		require.Equal(t, inputMetadata.Downloaded, resultMetadata.Downloaded)

		row := testDB.QueryRow(
			context.Background(),
			"SELECT metadata_id FROM metadata WHERE metadata_id = $1",
			resultMetadata.ID,
		)

		err = row.Scan(&metadataID)
		require.NoError(t, err)
		require.Equal(t, resultMetadata.ID, metadataID)
	})

	t.Run("Insert duplicate download_link error", func(t *testing.T) {
		inputMetadata := &models.MetaData{
			DownloadLink: downloadLink,
			Downloaded:   true,
		}

		_, err := repo.Insert(context.Background(), inputMetadata)
		require.Error(t, err)
	})

	t.Run("Get OK", func(t *testing.T) {
		resultMetadata, err := repo.Get(context.Background(), downloadLink)
		require.NoError(t, err)
		require.NotEmpty(t, resultMetadata.ID)
		require.NotEmpty(t, resultMetadata.DownloadLink)
		require.NotEmpty(t, resultMetadata.Downloaded)
		require.NotEmpty(t, resultMetadata.UpdatedAt)
	})

	t.Run("Get with error", func(t *testing.T) {
		_, err := repo.Get(context.Background(), "")
		require.Error(t, err)
		require.ErrorIs(t, err, models.ErrMetaDataNotFound)
	})

	t.Run("Update OK", func(t *testing.T) {
		now := time.Now()
		inputMetadata := &models.MetaData{
			ID:           metadataID,
			DownloadLink: downloadLinkAfterUpdate,
			Downloaded:   false,
		}

		updatedMetadata, err := repo.Update(
			context.Background(),
			inputMetadata,
		)

		require.NoError(t, err)
		require.Equal(t, inputMetadata.Downloaded, updatedMetadata.Downloaded)
		require.Equal(t, inputMetadata.DownloadLink, updatedMetadata.DownloadLink)
		require.GreaterOrEqual(t, updatedMetadata.UpdatedAt, now)
	})

	t.Run("Update with error", func(t *testing.T) {
		_, err := repo.Update(
			context.Background(),
			&models.MetaData{},
		)

		require.Error(t, err)
	})
}
