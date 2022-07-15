package coaster

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type CoasterMemmoryRepo struct {
	store map[string]Coaster
}

func initFromFile() map[string]Coaster {
	jsonFile, err := os.Open("/app/coasters.json")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully Opened coasters.json")
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var coasters []Coaster

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &coasters)

	store := make(map[string]Coaster, len(coasters))

	for i := 0; i < len(coasters); i++ {
		fmt.Println("ID: " + coasters[i].ID)
		store[coasters[i].ID] = coasters[i]
	}

	return store

}

func NewCoasterMemmoryRepo() CoasterMemmoryRepo {
	return CoasterMemmoryRepo{
		store: initFromFile(),
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
	if len(coaster.ID) == 0 {
		return errors.New("id fehlt")
	}
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

	return Coaster{}, errors.New("datensatz nicht gefunden")
}

func (repo CoasterMemmoryRepo) deleteCoaster(id string) error {
	_, ok := repo.store[id]
	if ok {
		delete(repo.store, id)
		return nil
	}

	return errors.New("datensatz nicht gefunden")
}
