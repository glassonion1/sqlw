package interactor

import (
	"fmt"

	"github.com/glassonion1/sqlw/sample/domain/model"
	"github.com/glassonion1/sqlw/sample/usecase/repository"
)

// Item is interactor for item model.
type Item struct {
	repo repository.Item
}

func NewItem(repo repository.Item) *Item {
	return &Item{
		repo: repo,
	}
}

// FindAll finds all items.
func (ii *Item) FindAll() ([]model.Item, error) {
	return ii.repo.FindAll()
}

// FindByID finds an item by specific id.
func (ii *Item) FindByID(id string) (*model.Item, error) {
	return ii.repo.FindByID(id)
}

// Create creates an item.
func (ii *Item) Create(item model.Item) (*model.Item, error) {
	if item.Name == "" {
		return nil, fmt.Errorf("name is a required field")
	}
	return ii.repo.Create(item)
}
