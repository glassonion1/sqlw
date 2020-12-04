package repository

import "github.com/glassonion1/sqlw/sample/domain/model"

// Item is interface for item repository
type Item interface {
	FindAll() ([]model.Item, error)
	FindByID(string) (*model.Item, error)
	Create(model.Item) (*model.Item, error)
}
