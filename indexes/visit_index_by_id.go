package indexes

import (
	"github.com/ArtyomNorin/hlcup2017/entities"
	"github.com/json-iterator/go"
	"sync"
)

type VisitIndexById struct {
	visits map[uint][]byte
	mutex  *sync.Mutex
}

func NewVisitIndexById() *VisitIndexById {
	return &VisitIndexById{visits: make(map[uint][]byte), mutex: new(sync.Mutex)}
}

func (visitIndexById *VisitIndexById) AddVisit(visit *entities.Visit) error {

	encodedVisit, err := jsoniter.ConfigFastest.Marshal(visit)

	if err != nil {
		return err
	}

	visitIndexById.mutex.Lock()

	visitIndexById.visits[*visit.Id] = encodedVisit

	visitIndexById.mutex.Unlock()

	return nil
}

func (visitIndexById *VisitIndexById) GetVisit(visitId uint) []byte {

	visitBytes, isIdExist := visitIndexById.visits[visitId]

	if !isIdExist {
		return nil
	}

	return visitBytes
}
