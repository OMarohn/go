package coaster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoasterService(t *testing.T) {
	var s CoasterService = NewCoasterService(NewCoasterMemmoryRepo())
	var testCoaster = Coaster{ID: "id123", Name: "TestCoaster", Manufacture: "TestManufature", Height: 123}
	var assert = assert.New(t)

	t.Run("Initial leer", func(t *testing.T) {
		c := s.getCoasters()
		assert.Len(c, 0)
	})

	t.Run("Anlegen eines leeren Coasters", func(t *testing.T) {
		err := s.createCoaster(Coaster{})
		assert.Error(err)
		assert.EqualError(err, "id fehlt")
	})

	t.Run("Anlegen eines vollständigen Coasters", func(t *testing.T) {
		err := s.createCoaster(testCoaster)
		assert.NoError(err)
	})

	t.Run("Lesen des angelegten Coasters", func(t *testing.T) {
		c, err := s.getCoaster("id123")
		assert.NoError(err)
		assert.Equal(c, testCoaster)
	})

	t.Run("Löschen eines unbekannten Coasters", func(t *testing.T) {
		err := s.deleteCoaster("id9999")
		assert.Error(err)
		assert.EqualError(err, "datensatz nicht gefunden")
	})

	t.Run("Löschen des angelegten Coasters", func(t *testing.T) {
		err := s.deleteCoaster("id123")
		assert.NoError(err)
	})

}
