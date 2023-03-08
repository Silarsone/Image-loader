package repository

import (
	"context"
	"fmt"
	"github.com/Silarsone/image-loader/internal/config"
	"github.com/Silarsone/image-loader/internal/model"
	"github.com/jmoiron/sqlx"
)

type image struct {
	ID        int    `db:"id"`
	UserID    int    `db:"user_id"`
	Name      string `db:"name"`
	Extension string `db:"extension"`
}

type ImageRepo struct {
	db  *sqlx.DB
	cfg *config.DB
}

func NewImageRepo(db *sqlx.DB, cfg *config.DB) *ImageRepo {
	return &ImageRepo{
		db:  db,
		cfg: cfg,
	}
}

func (i *ImageRepo) AddImage(ctx context.Context, modelImage model.Image) error {
	query := `INSERT INTO images(user_id, name, extension) VALUES (:user_id, :name, :extension)`

	image := convertImage(modelImage)

	_, err := i.db.NamedExecContext(ctx, query, &image)
	if err != nil {
		return fmt.Errorf("failed to insert image: %w", err)
	}

	return nil
}

func (i *ImageRepo) GetImage(ctx context.Context, id int) (model.Image, error) {
	query := `SELECT * FROM images WHERE id = $1`

	var img image

	row := i.db.QueryRowxContext(ctx, query, id)

	err := row.StructScan(&img)
	if err != nil {
		return model.Image{}, fmt.Errorf("failed to scan struct image: %w", err)
	}

	return img.toModel(), nil
}

func (i *ImageRepo) GetImages(ctx context.Context, userID int) ([]model.Image, error) {
	query := `SELECT * FROM images WHERE user_id = $1`

	var images []model.Image

	rows, err := i.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query images: %w", err)
	}

	for rows.Next() {
		var img image

		err := rows.StructScan(&img)
		if err != nil {
			return nil, err
		}

		images = append(images, img.toModel())
	}

	return images, nil
}

func convertImage(modelImage model.Image) image {
	return image{
		ID:        modelImage.ID,
		Name:      modelImage.Name,
		Extension: modelImage.Extension,
		UserID:    modelImage.UserID,
	}
}

func (i image) toModel() model.Image {
	return model.Image{
		ID:        i.ID,
		Name:      i.Name,
		Extension: i.Extension,
		UserID:    i.UserID,
	}
}
