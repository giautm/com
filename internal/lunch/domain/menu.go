package domain

type Menu struct {
}

type MenuService struct {
}

type Item struct {
	Name  string
	Price int32
}

func (s *MenuService) Create(question string, items []Item) error {
	return nil
}
