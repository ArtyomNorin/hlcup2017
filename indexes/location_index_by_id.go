package indexes

import (
	"github.com/ArtyomNorin/hlcup2017/entities"
	"github.com/json-iterator/go"
	"sync"
)

type LocationIndexById struct {
	locations map[uint][]byte
	mutex     *sync.Mutex
}

func NewLocationIndexById() *LocationIndexById {
	return &LocationIndexById{locations: make(map[uint][]byte), mutex: new(sync.Mutex)}
}

func (locationIndexById *LocationIndexById) AddLocation(location *entities.Location) error {

	encodedLocation, err := jsoniter.ConfigFastest.Marshal(location)

	if err != nil {
		return err
	}

	locationIndexById.mutex.Lock()

	locationIndexById.locations[*location.Id] = encodedLocation

	locationIndexById.mutex.Unlock()

	return nil
}

func (locationIndexById *LocationIndexById) GetLocation(locationId uint) []byte {

	locationBytes, isIdExist := locationIndexById.locations[locationId]

	if !isIdExist {
		return nil
	}

	return locationBytes
}
