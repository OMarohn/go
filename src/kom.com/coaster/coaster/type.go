package coaster

// Datenstruktur
type Coaster struct {
	Name        string `json:"name"`
	Manufacture string `json:"manufacture,omitempty"`
	ID          string `json:"id"`
	Height      int    `json:"height"`
}
