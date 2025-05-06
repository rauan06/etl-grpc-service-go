package repository

type CategoryRepository struct {
	url string
}

func NewCategoryRepository(url string) *CategoryRepository {
	return &CategoryRepository{
		url,
	}
}
