package review

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/alliancecup-server/pkg/client/postgres"
)

type Storage interface {
	Get(createdAt string, productId int) ([]SelectReviewsDTO, error)
	Create(dto CreateReviewDTO) (int, error)
	Delete(reviewId int) error
}

type storage struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
}

func NewReviewStorage(db *sqlx.DB, psql sq.StatementBuilderType) Storage {
	return &storage{
		db: db,
		qb: psql,
	}
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
	queryDeleteReview := fmt.Sprintf("DELETE FROM %s WHERE id=$1", postgres.ProductsReviewTable)

	_, err := s.db.Exec(queryDeleteReview, reviewId)
	if err != nil {
		return fmt.Errorf("failed to delete review due to: %v", err)
	}

	return nil
}

func (s *storage) Get(createdAt string, productId int) ([]SelectReviewsDTO, error) {
	var reviews []SelectReviewsDTO

	querySelectReviews := s.qb.Select("*").From(postgres.ProductsReviewTable)

	if createdAt != "" {
		querySelectReviews = querySelectReviews.Where(sq.Lt{"created_at": createdAt})
	}

	if productId != 0 {
		querySelectReviews = querySelectReviews.Where(sq.Eq{"product_id": productId})
	}

	ordered := querySelectReviews.OrderBy("created_at DESC").Limit(12)

	querySelectReviewsSql, args, err := ordered.ToSql()
	err = s.db.Select(&reviews, querySelectReviewsSql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select reviews from db, %v", err)
	}

	return reviews, nil
}
