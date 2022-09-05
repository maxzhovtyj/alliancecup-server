package review

type Service interface {
	Create(dto CreateReviewDTO) (int, error)
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
