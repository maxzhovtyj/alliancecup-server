package review

type Service interface {
	Get(createdAt string, productId int) ([]SelectReviewsDTO, error)
	Create(dto CreateReviewDTO) (int, error)
	Delete(reviewId int) error
}

type service struct {
	repo Storage
}

func NewReviewService(repo Storage) Service {
	return &service{repo: repo}
}

func (s *service) Create(dto CreateReviewDTO) (int, error) {
	return s.repo.Create(dto)
}

func (s *service) Delete(reviewId int) error {
	return s.repo.Delete(reviewId)
}

func (s *service) Get(createdAt string, productId int) ([]SelectReviewsDTO, error) {
	return s.repo.Get(createdAt, productId)
}
