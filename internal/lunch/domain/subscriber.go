package domain

type Subscriber struct {
	subscribeGroups map[int]bool
}

func (s *Subscriber) SubscribeForChat(groupId int, subscribe bool) {
	if subscribe {
		s.subscribeGroups[groupId] = subscribe
	} else {
		delete(s.subscribeGroups, groupId)
	}
}
