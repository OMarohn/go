# Kleine Fingerübung um in das Thema golang reinzukommen

Angefangen bin ich mit [kubucation/go-rollercoaster-ap](https://github.com/kubucation/go-rollercoaster-api) mir mal das Thema REST-API in **golang** anzuschauen.

Die nächsten Schritte waren dann der Umbau auf die **hexagonale Architektur** also  Service - Repository (in Memmory und redis) und Port (grpc und zwei REST-Service Implementierungen)

Im REST Port (gorilla) habe ich mir auch das Thema AuthN/AuthZ mit **JWT** mal angetan, damit ich im K8s mit ISTIO ein wenig rumprobieren kann, was da geht. 

Die Anwendung lässt sich einfach containerisieren -> Dockerfile

Todo:
- [JWT]([RFC7519](https://datatracker.ietf.org/doc/html/rfc7519)) Routen noch AuthZ fähig machen - scope auswerten bzw. claim upn und group füllen
- OpenAPI / Swagger Erzeugen für die REST Routen
- graphQL Port bauen