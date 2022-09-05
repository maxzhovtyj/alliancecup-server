package review

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
)

type Storage interface {
	Create(dto CreateReviewDTO) (int, error)
	Delete(reviewId int) error
}

type storage struct {
	db *sqlx.DB
}

func NewReviewStorage(db *sqlx.DB) Storage {
	return &storage{db: db}
}

func (s *storage) Create(dto CreateReviewDTO) (int, error) {
	queryCreateReview := fmt.Sprintf(
		"INSERT INTO %s (user_id, user_name, mark, review_text) values ($1, $2, $3, $4) RETURNING id",
		postgres.ProductsReviewTable,
	)

	var reviewId int
	row := s.db.QueryRow(queryCreateReview, dto.UserId, dto.UserName, dto.Mark, dto.ReviewText)
	if err := row.Scan(&reviewId); err != nil {
		return 0, fmt.Errorf("failed to scan review id, due to: %v", err)
	}

	return reviewId, nil
}

func (s *storage) Delete(reviewId int) error {
	queryDeleteReview := fmt.Sprintf("DELETE FROM %s WHERE review_id=$1", postgres.ProductsReviewTable)

	_, err := s.db.Exec(queryDeleteReview, reviewId)
	if err != nil {
		return fmt.Errorf("failed to delete review due to: %v", err)
	}

	return nil
}
