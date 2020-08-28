package domain

import (
	"fmt"
	"time"
)

type Menu struct {
	Question  string
	Items     []MenuItem
	Closed    bool
	CloseDate time.Time
	Chooses   map[int]int
}

func NewMenuWithFlatPrice(items []string, flatPrice int64) *Menu {
	m := &Menu{
		Question: "Trưa nay ăn gì?",
		Closed:   false,
	}

	m.Items = make([]MenuItem, len(items))
	for idx, item := range items {
		m.Items[idx].Name = item
		m.Items[idx].Price = flatPrice
	}

	return m
}

func (m Menu) Chooser() []int {
	users := make([]int, 0)
	for uid := range m.Chooses {
		users = append(users, uid)
	}
	return users
}

func (m *Menu) Choose(userId int, optionIdx int) {
	if optionIdx >= 0 && optionIdx < len(m.Items) {
		m.Chooses[userId] = optionIdx
	}
}

func (m *Menu) RetractChoose(userId int) {
	delete(m.Chooses, userId)
}

type MenuItem struct {
	Name  string
	Price int64
}

func (m MenuItem) String() string {
	return fmt.Sprintf("%s - %d", m.Name, m.Price)
}
