package domain

import (
	"context"
	"fmt"
)

type Menu struct {
	items []MenuItem
}

type MenuRepo interface {
	CreateMenu(ctx context.Context, m *Menu, createFn func() error) error
	CreateMenuWithOrder(ctx context.Context, m *Menu, createFn func() (*Order, error)) error
}

type MenuSender interface {
	SendMenu(ctx context.Context, items []MenuItem) (string, error)
}

func NewMenu() *Menu {
	return &Menu{}
}

func NewMenuFlatPrice(itemsName []string, flatPrice int) *Menu {
	items := make([]MenuItem, len(itemsName))
	for idx, o := range itemsName {
		items[idx] = MenuItem{
			Name:  o,
			Price: flatPrice,
		}
	}

	return &Menu{
		items: items,
	}
}

func (m *Menu) AddItem(item MenuItem) {
	m.items = append(m.items, item)
}

func (m Menu) Send(ctx context.Context, sender MenuSender) (*Order, error) {
	id, err := sender.SendMenu(ctx, m.items)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID: id,
	}, nil
}

func (m Menu) Options() []string {
	opts := make([]string, len(m.items))
	for idx, i := range m.items {
		opts[idx] = i.String()
	}
	return opts
}

type MenuItem struct {
	Name  string
	Price int
}

func (m MenuItem) String() string {
	return fmt.Sprintf("%s - %d", m.Name, m.Price)
}
