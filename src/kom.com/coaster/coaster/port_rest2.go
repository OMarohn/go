// Das ist eine REST Implementierung mit dem Gorilla mux Package
package coaster

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Der Port als REST Implementierung
type CoasterRestPort2 struct {
	service CoasterService
}

func NewCoasterRestPort2(theService CoasterService) CoasterRestPort2 {
	return CoasterRestPort2{service: theService}
}

func (port CoasterRestPort2) HandleList(w http.ResponseWriter, r *http.Request) {
	coasters := port.service.getCoasters()
	result, err := json.Marshal(coasters)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (port CoasterRestPort2) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var cItem Coaster
	err := json.NewDecoder(r.Body).Decode(&cItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = port.service.createCoaster(cItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (port CoasterRestPort2) HandleGetOne(w http.ResponseWriter, r *http.Request) {
	theVars := mux.Vars(r)
	id := theVars["id"]
	log.Printf("parameter:= %v", id)

	coaster, err := port.service.getCoaster(id)
	log.Printf("Coaster:= %v", coaster)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := json.Marshal(coaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (port CoasterRestPort2) HandleDelete(w http.ResponseWriter, r *http.Request) {
	theVars := mux.Vars(r)
	id := theVars["id"]
	log.Printf("parameter:= %v", id)

	err := port.service.deleteCoaster(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
