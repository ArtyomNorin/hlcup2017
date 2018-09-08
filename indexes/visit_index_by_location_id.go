package indexes

import (
		"hlcup/entities"
	"sync"
)

type VisitIndexByLocationId struct {
	visits map[uint][]uint
	mutex  *sync.Mutex
}

func NewVisitIndexByLocationId() *VisitIndexByLocationId{
	return &VisitIndexByLocationId{visits: make(map[uint][]uint), mutex: new(sync.Mutex)}
}

func (visitIndexByLocationId *VisitIndexByLocationId) AddVisit(visit *entities.Visit) {

	visitIndexByLocationId.mutex.Lock()

	visitIndexByLocationId.visits[visit.Location] = append(visitIndexByLocationId.visits[visit.Location], visit.Id)

	visitIndexByLocationId.mutex.Unlock()
}

func (visitIndexByLocationId *VisitIndexByLocationId) GetVisits(locationId uint) []uint {

	visits, isIdExist := visitIndexByLocationId.visits[locationId]

	if !isIdExist {
		return nil
	}

	return visits
}
