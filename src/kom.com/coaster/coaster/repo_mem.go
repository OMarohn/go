package coaster

import (
	"errors"
)

type CoasterMemmoryRepo struct {
	store map[string]Coaster
}

func NewCoasterMemmoryRepo() CoasterMemmoryRepo {
	return CoasterMemmoryRepo{
		store: map[string]Coaster{
			"id1": {
				Name:        "Coaster 1",
				Manufacture: "Manufacture 1",
				ID:          "id1",
				Height:      100,
			},
			"id2": {
				Name:        "Coaster 2",
				Manufacture: "Manufacture 2",
				ID:          "id2",
				Height:      200,
			}},
	}
}

func (repo CoasterMemmoryRepo) getCoasters() []Coaster {

	// alloc für Kopie
	coasters := make([]Coaster, len(repo.store))
	i := 0
	// kopieren
	for _, coaster := range repo.store {
		coasters[i] = coaster
		i++
	}
	return coasters
}

func (repo CoasterMemmoryRepo) createCoaster(coaster Coaster) error {
	if (repo.store[coaster.ID] != Coaster{}) {
		return errors.New("datensatz bereits existent")
	}
	repo.store[coaster.ID] = coaster
	return nil
}

func (repo CoasterMemmoryRepo) getCoaster(id string) (Coaster, error) {
	c, ok := repo.store[id]
	if ok {
		return c, nil
	}

	return Coaster{}, errors.New("nicht implementiert")
}

func (repo CoasterMemmoryRepo) deleteCoaster(id string) error {
	return errors.New("nicht implementiert")
}
