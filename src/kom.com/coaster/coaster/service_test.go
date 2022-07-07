package coaster

import (
	"reflect"
	"testing"
)

func TestCoasterService(t *testing.T) {
	var s CoasterService = NewCoasterService(NewCoasterMemmoryRepo())
	var testCoaster = Coaster{ID: "id123", Name: "TestCoaster", Manufacture: "TestManufature", Height: 123}

	t.Run("Initial leer", func(t *testing.T) {
		c := s.getCoasters()
		if len(c) == 0 {
			t.Log("OK -- keine Einträge vorhanden")
		}
	})

	t.Run("Anlegen eines leeren Coasters", func(t *testing.T) {
		err := s.createCoaster(Coaster{})
		if err != nil {
			t.Log("OK -- Fehler bei dem Versuch einen Leeren DS zu erstellen")
		} else {
			t.Error("Es wurde ein Fehler erwartet")
		}
	})

	t.Run("Anlegen eines vollständigen Coasters", func(t *testing.T) {
		err := s.createCoaster(testCoaster)
		if err == nil {
			t.Log("OK -- Anlage OK")
		} else {
			t.Errorf("Fehler bei der Anlage eines Coasters, %v", err)
		}
	})

	t.Run("Lesen des angelegten Coasters", func(t *testing.T) {
		c, err := s.getCoaster("id123")
		if err == nil && reflect.DeepEqual(c, testCoaster) {
			t.Log("OK -- Lesen des Coster")
		} else {
			t.Errorf("Fehler beim lesen des Coasters, %v", err)
		}
	})

	t.Run("Löschen eines unbekannten Coasters", func(t *testing.T) {
		err := s.deleteCoaster("id9999")
		if err != nil {
			t.Log("OK -- Löschen unbekannter Coster führt zum Fehler")
		} else {
			t.Error("Es wurde keine Fehlermeldung geworfen")
		}
	})

	t.Run("Löschen des angelegten Coasters", func(t *testing.T) {
		err := s.deleteCoaster("id123")
		if err == nil {
			t.Log("OK -- Löschen des Coster")
		} else {
			t.Errorf("Fehler beim Löschen des Coasters, %v", err)
		}
	})

}
