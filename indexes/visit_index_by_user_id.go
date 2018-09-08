package indexes

import (
		"hlcup/entities"
	"sync"
)

type VisitIndexByUserId struct {
	visits map[uint][]uint
	mutex  *sync.Mutex
}

func NewVisitIndexByUserId() *VisitIndexByUserId{
	return &VisitIndexByUserId{visits: make(map[uint][]uint), mutex: new(sync.Mutex)}
}

func (visitIndexByUserId *VisitIndexByUserId) AddVisit(visit *entities.Visit) {

	visitIndexByUserId.mutex.Lock()

	visitIndexByUserId.visits[visit.User] = append(visitIndexByUserId.visits[visit.User], visit.Id)

	visitIndexByUserId.mutex.Unlock()
}

func (visitIndexByUserId *VisitIndexByUserId) GetVisits(userId uint) []uint {

	visits, isIdExist := visitIndexByUserId.visits[userId]

	if !isIdExist {
		return nil
	}

	return visits
}
