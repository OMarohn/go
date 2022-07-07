package jwkstools

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/form3tech-oss/jwt-go"
)

// Struktur des JWKS
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type JWKManager struct {
	Url       string
	Filename  string
	certStore map[string]*rsa.PublicKey
}

// JWKS Ergbniss in eine Map id|rsaPublic Key wandeln
func (jwkm *JWKManager) mapJwks2Store(items Jwks) (int, error) {
	jwkm.certStore = make(map[string]*rsa.PublicKey)
	for k := range items.Keys {
		cert := "-----BEGIN CERTIFICATE-----\n" + items.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			log.Println("Public Key nicht gefunden")
		}
		jwkm.certStore[items.Keys[k].Kid] = pubKey
	}

	return len(jwkm.certStore), nil
}

// Remote lesen des JWKS
func (jwkm *JWKManager) getJWKS(url string) (int, error) {
	resp, err := http.Get(url)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return 0, err
	}

	return jwkm.mapJwks2Store(jwks)
}

// Public Key ermitteln
func (jwkm *JWKManager) GetPemCert(token *jwt.Token) (*rsa.PublicKey, error) {

	kid := fmt.Sprintf("%v", token.Header["kid"])
	pk, found := jwkm.certStore[kid]
	if !found {
		err := errors.New("konnte key für das zertifikat nicht finden")
		return nil, err
	}

	return pk, nil
}

func (jwkm *JWKManager) InitCertStore() {
	cnt, err := jwkm.getJWKS(jwkm.Url)
	if err != nil {
		log.Println("Keine Zerifikate im Store gefunden!", err)
		// Datei mit JWKS lesen - failover
		if len(jwkm.Filename) > 0 {
			jsonFile, err := os.Open(jwkm.Filename)

			if err != nil {
				panic("Auch keine JWKS-Datei gefunden!")
			} else {
				defer jsonFile.Close()
				log.Println("File OK")
				var jwks = Jwks{}
				err = json.NewDecoder(jsonFile).Decode(&jwks)
				if err != nil {
					panic(err)
				} else {
					cnt, _ = jwkm.mapJwks2Store(jwks)
				}
			}
		}
	}
	log.Println("Anzahl Zerifikate: ", cnt)
}
