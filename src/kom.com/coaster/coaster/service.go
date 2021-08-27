package coaster

// Der Service
type CoasterService struct {
	repo CoasterRepo
}

func NewCoasterService(theRepo CoasterRepo) CoasterService {
	return CoasterService{repo: theRepo}
}

func (cs CoasterService) getCoasters() []Coaster {
	return cs.repo.getCoasters()
}

func (cs CoasterService) getCoaster(id string) (Coaster, error) {
	return cs.repo.getCoaster(id)
}

func (cs CoasterService) createCoaster(coaster Coaster) error {
	return cs.repo.createCoaster(coaster)
}

func (cs CoasterService) deleteCoaster(id string) error {
	return cs.repo.deleteCoaster(id)
}
