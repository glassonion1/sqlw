package rest

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/glassonion1/sqlw/sample/domain/model"
	"github.com/glassonion1/sqlw/sample/usecase/interactor"
)

// ItemHandler is http handler for item resources
type ItemHandler struct {
	interactor *interactor.Item
}

func NewItemHandler(interactor *interactor.Item) *ItemHandler {
	return &ItemHandler{
		interactor: interactor,
	}
}

// List gets all items.
func (h *ItemHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		items, err := h.interactor.FindAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, items)
	}
}

// Get gets an item by specific item_id of url param.
func (h *ItemHandler) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("item_id")
		item, err := h.interactor.FindByID(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, item)
	}
}

// Create creates an item.
func (h *ItemHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		item := model.Item{}
		if err := c.Bind(&item); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		new, err := h.interactor.Create(item)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, new)
	}
}
