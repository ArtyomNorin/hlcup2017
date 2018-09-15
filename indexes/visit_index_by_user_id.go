package indexes

import (
	"hlcup/entities"
	"sync"
	"sort"
	)

type VisitIndexByUserId struct {
	visits map[uint][]uint
	mutex  *sync.Mutex
}

func NewVisitIndexByUserId() *VisitIndexByUserId {
	return &VisitIndexByUserId{visits: make(map[uint][]uint), mutex: new(sync.Mutex)}
}

func (visitIndexByUserId *VisitIndexByUserId) AddVisit(visit *entities.Visit) {

	visitIndexByUserId.mutex.Lock()

	visitIndexByUserId.visits[*visit.User] = append(visitIndexByUserId.visits[*visit.User], *visit.Id)

	sort.Slice(visitIndexByUserId.visits[*visit.User], func(i, j int) bool {
		return visitIndexByUserId.visits[*visit.User][i] < visitIndexByUserId.visits[*visit.User][j]
	})

	visitIndexByUserId.mutex.Unlock()
}

func (visitIndexByUserId *VisitIndexByUserId) GetVisits(userId uint) []uint {

	visits, isIdExist := visitIndexByUserId.visits[userId]

	if !isIdExist {
		return nil
	}

	return visits
}

func (visitIndexByUserId *VisitIndexByUserId) DeleteVisit(userId uint, visitId uint) {

	visitIndexByUserId.mutex.Lock()

	visitsByUserId, isUserExist := visitIndexByUserId.visits[userId]

	if !isUserExist || len(visitsByUserId) == 0 {
		visitIndexByUserId.mutex.Unlock()
		return
	}

	if len(visitsByUserId) == 1 {
		delete(visitIndexByUserId.visits, userId)
		visitIndexByUserId.mutex.Unlock()
		return
	}

	visitIndex := sort.Search(len(visitsByUserId) - 1, func(index int) bool {
		return visitsByUserId[index] >= visitId
	})

	if visitIndex == len(visitsByUserId)-1 {
		visitIndexByUserId.visits[userId] = visitsByUserId[:visitIndex]
	} else {
		visitIndexByUserId.visits[userId] = append(visitsByUserId[:visitIndex], visitsByUserId[visitIndex+1:]...)
	}

	visitIndexByUserId.mutex.Unlock()
}
