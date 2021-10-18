package coaster

import (
	"database/sql"
	"log"
)

// Die Implementierung eines Postgres Repository -- ohne OR-Mapper
type CoasterPostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(coasterDB *sql.DB) CoasterPostgresRepo {
	return CoasterPostgresRepo{db: coasterDB}
}

func (repo CoasterPostgresRepo) getCoasters() []Coaster {

	var ret []Coaster

	rows, err := repo.db.Query("SELECT id, name, manufacture, height FROM coaster")

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Daten aus dem JSON deserialisieren
	var coasterItem Coaster
	for rows.Next() {
		err = rows.Scan(&coasterItem.ID, &coasterItem.Name, &coasterItem.Manufacture, &coasterItem.Height)
		if err != nil {
			log.Println(err)
		} else {
			ret = append(ret, coasterItem)
		}
	}

	return ret
}

// Anlegen eines Datums
func (repo CoasterPostgresRepo) createCoaster(coaster Coaster) error {
	sqlStatement := `
		INSERT INTO coaster ( id, name, manufacture, height)
		VALUES ($1, $2, $3, $4)`

	_, err := repo.db.Exec(sqlStatement, coaster.ID, coaster.Name, coaster.Manufacture, coaster.Height)
	return err
}

// Ein Datensatz über ID lesen
func (repo CoasterPostgresRepo) getCoaster(id string) (Coaster, error) {
	ret := Coaster{}
	rows, err := repo.db.Query(`SELECT id, name, manufacture, height FROM coaster WHERE id = $1`, id)
	if err == nil {
		defer rows.Close()
		rows.Next()
		err = rows.Scan(&ret.ID, &ret.Name, &ret.Manufacture, &ret.Height)
		return ret, err
	} else {
		return Coaster{}, err
	}
}

// Ein Datensatz löschen
func (repo CoasterPostgresRepo) deleteCoaster(id string) error {
	sqlStatement := `
		DELETE from coaster where id = $1`

	_, err := repo.db.Exec(sqlStatement, id)
	return err
}
