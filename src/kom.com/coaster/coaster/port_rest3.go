// Das ist eine REST Implementierung mit dem Gorilla mux Package
package coaster

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Der Port als REST Implementierung
type CoasterRestPort3 struct {
	service CoasterService
}

func NewCoasterRestPort3(theService CoasterService) CoasterRestPort3 {
	return CoasterRestPort3{service: theService}
}

func (port CoasterRestPort3) HandleList(c echo.Context) error {
	coasters := port.service.getCoasters()
	c.JSON(http.StatusOK, coasters)
	return nil
}

func (port CoasterRestPort3) HandleCreate(c echo.Context) error {
	var cItem Coaster

	if err := c.Bind(&cItem); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := port.service.createCoaster(cItem)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusCreated, cItem)
	return nil
}

func (port CoasterRestPort3) HandleGetOne(c echo.Context) error {

	id := c.Param("id")
	log.Printf("parameter:= %v", id)

	coaster, err := port.service.getCoaster(id)
	log.Printf("Coaster:= %v", coaster)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	c.JSON(http.StatusOK, coaster)
	return nil
}

func (port CoasterRestPort3) HandleDelete(c echo.Context) error {
	id := c.Param("id")
	err := port.service.deleteCoaster(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	c.NoContent(http.StatusOK)
	return nil
}
