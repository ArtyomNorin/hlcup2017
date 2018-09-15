package indexes

import (
	"hlcup/entities"
	"sync"
	"sort"
)

type VisitIndexByLocationId struct {
	visits map[uint][]uint
	mutex  *sync.Mutex
}

func NewVisitIndexByLocationId() *VisitIndexByLocationId {
	return &VisitIndexByLocationId{visits: make(map[uint][]uint), mutex: new(sync.Mutex)}
}

func (visitIndexByLocationId *VisitIndexByLocationId) AddVisit(visit *entities.Visit) {

	visitIndexByLocationId.mutex.Lock()

	visitIndexByLocationId.visits[*visit.Location] = append(visitIndexByLocationId.visits[*visit.Location], *visit.Id)

	sort.Slice(visitIndexByLocationId.visits[*visit.Location], func(i, j int) bool {
		return visitIndexByLocationId.visits[*visit.Location][i] < visitIndexByLocationId.visits[*visit.Location][j]
	})

	visitIndexByLocationId.mutex.Unlock()
}

func (visitIndexByLocationId *VisitIndexByLocationId) GetVisits(locationId uint) []uint {

	visits, isIdExist := visitIndexByLocationId.visits[locationId]

	if !isIdExist {
		return nil
	}

	return visits
}

func (visitIndexByLocationId *VisitIndexByLocationId) DeleteVisit(locationId uint, visitId uint) {

	visitIndexByLocationId.mutex.Lock()

	visitsByLocationId, isLocationExist := visitIndexByLocationId.visits[locationId]

	if !isLocationExist || len(visitsByLocationId) == 0 {
		visitIndexByLocationId.mutex.Unlock()
		return
	}

	if len(visitsByLocationId) == 1 {
		delete(visitIndexByLocationId.visits, locationId)
		visitIndexByLocationId.mutex.Unlock()
		return
	}

	visitIndex := sort.Search(len(visitsByLocationId) - 1, func(index int) bool {
		return visitsByLocationId[index] >= visitId
	})

	if visitIndex == len(visitsByLocationId)-1 {
		visitIndexByLocationId.visits[locationId] = visitsByLocationId[:visitIndex]
	} else {
		visitIndexByLocationId.visits[locationId] = append(visitsByLocationId[:visitIndex], visitsByLocationId[visitIndex+1:]...)
	}

	visitIndexByLocationId.mutex.Unlock()
}
