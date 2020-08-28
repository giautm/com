package lunch

import (
	"context"
	"sync"

	"giautm.dev/com/internal/lunch/domain"
	"github.com/google/uuid"
)

type MemoryGroupRepo struct {
	lock   *sync.RWMutex
	groups map[int64]domain.Group
	polls  map[uuid.UUID]domain.Poll
}

func NewMemoryGroupRepo() *MemoryGroupRepo {
	return &MemoryGroupRepo{
		lock:   &sync.RWMutex{},
		groups: make(map[int64]domain.Group),
		polls:  make(map[uuid.UUID]domain.Poll),
	}
}

func (m MemoryGroupRepo) getOrCreateGroup(chatID int64) (*domain.Group, error) {
	group, ok := m.groups[chatID]
	if !ok {
		return &domain.Group{}, nil
	}

	return &group, nil
}

func (m *MemoryGroupRepo) UpdateGroup(
	ctx context.Context,
	chatID int64,
	updateFn func(*domain.Group) (*domain.Group, error),
) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	group, err := m.getOrCreateGroup(chatID)
	if err != nil {
		return err
	}

	updatedGroup, err := updateFn(group)
	if err != nil {
		return err
	}

	m.groups[chatID] = *updatedGroup

	return nil
}

func (m *MemoryGroupRepo) UpdateGroupAndPoll(
	ctx context.Context,
	chatID int64,
	updateFn func(*domain.Group) (*domain.Group, *domain.Poll, error),
) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	group, err := m.getOrCreateGroup(chatID)
	if err != nil {
		return err
	}

	updatedGroup, newPoll, err := updateFn(group)
	if err != nil {
		return err
	}

	m.groups[chatID] = *updatedGroup
	if newPoll != nil {
		m.polls[newPoll.ID] = *newPoll
	}

	return nil
}
