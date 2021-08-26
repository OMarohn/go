package coaster

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Das scheint mir irgendwie unschön!
var ctx = context.Background()

// Der Port als REST Implementierung
type CoasterRestPort struct {
	service CoasterService
}

func NewCoasterRestPort(theService CoasterService) CoasterRestPort {
	return CoasterRestPort{service: theService}
}

// Auch interessant, der fn Name muss Grossgeschriebenw werden sonst wird er nicht exportiert, obwohl ja der CoasterRestPort exportiert wird !
func (port CoasterRestPort) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		coasters := port.service.getCoasters()
		result, err := json.Marshal(coasters)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)

	case http.MethodPost:
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
}

func (port CoasterRestPort) HandleGet(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	log.Printf("Path:= %v", parts)

	if len(parts) == 3 {
		coaster, err := port.service.getCoaster(parts[2])
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
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
