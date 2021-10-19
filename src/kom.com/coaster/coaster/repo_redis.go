package coaster

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
)

// Die Implementierung eines REDIS Repository
type CoasterRedisRepo struct {
	rclient *redis.Client
}

var ctx = context.Background()

func NewRedisRepo(redisClient *redis.Client) CoasterRedisRepo {
	return CoasterRedisRepo{rclient: redisClient}
}

// Ermitteln aller Daten hier auf 50 max. gedrosselt. Pageing müsste noch rein
func (repo CoasterRedisRepo) getCoasters() []Coaster {

	var ret []Coaster = make([]Coaster, 0)

	// schlüssel lesen (max. 50)
	keys, _, err := repo.rclient.Scan(ctx, 0, "coaster*", 50).Result()

	if err != nil {
		panic(err)
	}

	if len(keys) > 0 {

		// Daten lesen
		erg, err := repo.rclient.MGet(ctx, keys...).Result()

		if err != nil {
			panic(err)
		}

		// Daten aus dem JSON deserialisieren
		var coasterItem Coaster
		for _, item := range erg {

			err := json.Unmarshal([]byte(item.(string)), &coasterItem)
			if err != nil {
				panic(err)
			}

			ret = append(ret, coasterItem)

		}
	}

	return ret
}

// Anlegen eines Datums
func (repo CoasterRedisRepo) createCoaster(coaster Coaster) error {
	// Checken ob schon vorhanden
	erg, err := repo.rclient.Exists(ctx, "coaster."+coaster.ID).Result()
	if err != nil {
		panic(err)
	}
	if erg == 0 {
		jsonBytes, err := json.Marshal(coaster)
		if err != nil {
			panic(err)
		}
		repo.rclient.Set(ctx, "coaster."+coaster.ID, jsonBytes, 0)
		return nil
	} else {
		return errors.New("datensatz bereits existent")
	}
}

// Ein Datensatz über ID lesen
func (repo CoasterRedisRepo) getCoaster(id string) (Coaster, error) {

	coaster, err := repo.rclient.Get(ctx, "coaster."+id).Result()

	var coasterItem Coaster
	if err == nil {
		err = json.Unmarshal([]byte(coaster), &coasterItem)
		return coasterItem, err
	} else {
		return Coaster{}, errors.New("nicht gefunden")
	}
}

// Ein Datensatz löschen
func (repo CoasterRedisRepo) deleteCoaster(id string) error {

	cnt, err := repo.rclient.Del(ctx, "coaster."+id).Result()

	if err == nil && cnt == 1 {
		return nil
	} else {
		return errors.New("nicht gefunden")
	}
}
