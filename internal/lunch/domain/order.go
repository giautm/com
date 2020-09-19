package domain

import (
	"time"
)

type Order struct {
	ID        string
	CloseDate time.Time

	items  []OrderItem
	closed bool
}

func (o Order) CanAddMore() bool {
	return !o.closed
}

func (o *Order) Close() {
	o.closed = true
}

type menuI interface {
	ID() string
	GetItem(idx int) *MenuItem
}

func (o *Order) PickMenu(menu menuI, items []int) {
	menuID := menu.ID()
	for _, idx := range items {
		if i := menu.GetItem(idx); i != nil {
			o.items = append(o.items, OrderItem{
				menuID: menuID,
				itemID: idx,
				name:   i.Name,
				qty:    1,
				price:  i.Price,
			})
		}
	}
}

type OrderItem struct {
	menuID string
	itemID int
	name   string
	qty    int
	price  int
}
