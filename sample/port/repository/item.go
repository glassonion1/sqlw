package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/glassonion1/sqlw"
	"github.com/glassonion1/sqlw/sample/domain/model"
)

// Item is repository implementation for item model.
type Item struct {
	db *sqlw.DB
}

func NewItem(db *sqlw.DB) *Item {
	return &Item{
		db: db,
	}
}

// FindAll finds all items
func (r *Item) FindAll() ([]model.Item, error) {
	rows, err := r.db.Query("SELECT * FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []model.Item{}
	for rows.Next() {
		item := model.Item{}
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// FindByID finds an items by specific id.
func (r *Item) FindByID(id string) (*model.Item, error) {
	row := r.db.QueryRow("SELECT * FROM items WHERE id=?", id)
	if row == nil {
		return nil, errors.New("data not found")
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	item := model.Item{}
	if err := row.Scan(&item.ID, &item.Name); err != nil {
		return nil, err
	}

	return &item, nil
}

// Create creates an item.
func (r *Item) Create(item model.Item) (*model.Item, error) {
	id, _ := uuid.NewUUID()
	_, err := r.db.Exec("INSERT INTO items(id, name) VALUES(?, ?)", id.String(), item.Name)
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRowForMaster("SELECT * FROM items WHERE id=?", id.String())
	if row == nil {
		return nil, fmt.Errorf("data not found")
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	new := &model.Item{}
	if err := row.Scan(&new.ID, &new.Name); err != nil {
		return nil, err
	}

	return new, nil
}
