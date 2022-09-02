package category

type Service interface {
	GetAll() ([]Category, error)
	GetFiltration(fkName string, id int) ([]Filtration, error)
	Update(category Category) (int, error)
	Create(category Category) (int, error)
	Delete(id int) error
	AddFiltration(filtration Filtration) (int, error)
}

type service struct {
	repo Storage
}

func NewCategoryService(repo Storage) Service {
	return &service{repo: repo}
}

func (c *service) GetAll() ([]Category, error) {
	return c.repo.GetAll()
}

func (c *service) Update(category Category) (int, error) {
	return c.repo.Update(category)
}

func (c *service) Create(category Category) (int, error) {
	id, err := c.repo.Create(category)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (c *service) Delete(id int) error {
	return c.repo.Delete(id)
}

func (c *service) AddFiltration(filtration Filtration) (int, error) {
	return c.repo.AddFiltration(filtration)
}

func (c *service) GetFiltration(fkName string, id int) ([]Filtration, error) {
	return c.repo.GetFiltration(fkName, id)
}
