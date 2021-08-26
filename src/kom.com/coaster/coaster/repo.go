package coaster

// Repository Interface
type CoasterRepo interface {
	getCoasters() []Coaster
	getCoaster(id string) (Coaster, error)
	createCoaster(coaster Coaster) error
}
