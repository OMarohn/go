# Anbinden eines externen AuthZ-Providers

Ein externer AutZ Provider wird von ISTIO mit Hilfe eines Envoy-Filters unterstützt. Da das schreiben eines Envoy-Filters relativ komplex ist, hat sich ISTIO für diese Lösung entschieden. In ISTIO werden ein - oder auch mehrere - Callback-Services für den Filter registeriert. [Envoy](https://www.envoyproxy.io/) ist der in ISTIO verwendete Proxy.

Hier wollte ich untersuchen, ob das ein Lösungsansatz für ein einfachesres **siltend login** sein könnte. Der KeyCloak bietet - zwar aktuell nicht offiziell unterstützt - einen [Token Exchange](https://www.keycloak.org/docs/latest/securing_apps/#_token-exchange) Endpunkt an, der von so einem AuthZ-Provider  verwendet werden könnte.

## Zuerst nochmal das Procedere um ISTIO zu installieren

Download von ISTIO -- aktuelles Release (stand jetzt 1.14.1)

```bash
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.14.1
export PATH=$PWD/bin:$PATH
```

>Installation mit Egress - nicht unbedignt notwendig

```bash
istioctl install --set profile=demo -y --set meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY`
```

## Optional die schicken Tools

Zum Auswerten und analysieren machen die Tools auf jeden Fall Sinn 😃

### Prometheus

```bash
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.14/samples/addons/prometheus.yaml
```

### Grafana

```bash
kubectl apply -f <https://raw.githubusercontent.com/istio/istio/release-1.14/samples/addons/grafana.yaml>
```

### Jaeger

```bash
kubectl apply -f <https://raw.githubusercontent.com/istio/istio/release-1.14/samples/addons/jaeger.yaml>
```

### Kiali

```bash
kubectl apply -f <https://raw.githubusercontent.com/istio/istio/release-1.14/samples/addons/kiali.yaml>
```

## Konfiguration des Externen AuthZ Providers

Der AuthZ Provider kann eine REST und/oder ein GRPC - Schnittstelle haben. In der ISTIO Doku ist ein [Beispielprovider](https://github.com/istio/istio/blob/master/samples/extauthz/cmd/extauthz/main.go) in golang implementiert der beide Schnittstellen exemplarisch implementiert. Diesen habe ich in der REST-Implementierung leicht angepasst um das Austauschen einen Tokens zu simulieren. Den Rest der Implementierung habe ich nicht geändert, da er auch noch andere Aspekte wie das Auswerten von Headern adressiert.

[Weitere Info zum Thema ExtAuthZ von ISTIO](https://istio.io/latest/docs/tasks/security/authorization/authz-custom/)

Mehr oder weniger die ganze *Magie* findet sich in der folgenden überschaubaren Methode.

```golang
func (s *ExtAuthzServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
 body, err := io.ReadAll(request.Body) // lesen des Request-Bodies 

 if err != nil {
  log.Printf("[HTTP] read body failed: %v", err)
 }

 l := fmt.Sprintf("%s %s%s, headers: %v, body: [%s]\n", request.Method, request.Host, request.URL, request.Header, body)
 if allowedValue == request.Header.Get(checkHeader) { // Checken ob der x-ext-authz Header === allow ist
  response.Header().Set(resultHeader, resultAllowed)
  response.Header().Set(overrideHeader, request.Header.Get(overrideHeader))
  response.Header().Set(receivedHeader, l)
  newToken, err := getNewRWToken() // neuen Token besorgen -- hier Token Exchange anbinden

  if err == nil { // alles OK, dann Token in den Upstream Header setzen!
   response.Header().Set("Authorization", "Bearer "+newToken.Token)
   response.Header().Set("x-auth-request-access-token", "Bearer "+newToken.Token)
   response.WriteHeader(http.StatusOK)
  } else { // Ansonsten Zugriff verweigern
   response.WriteHeader(http.StatusForbidden)
  }
 } else { // auch Zugriff verweigern - ein paar informative Header setzen
  response.Header().Set(resultHeader, resultDenied)
  response.Header().Set(overrideHeader, request.Header.Get(overrideHeader))
  response.Header().Set(receivedHeader, l)
  response.WriteHeader(http.StatusForbidden)
  response.Write([]byte(denyBody))
 }
}
```

## Provider Deployen

Der Provider wird im K8s deployed hier im Namespace **extauth**

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: extauth
  namespace: extauth
spec:
  selector:
    matchLabels:
      app: extauth
  template:
    metadata:
      labels:
        app: extauth
        version: v1
    spec:
      containers:
      - name: extauth
        image: omarohn/extauth 
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
        ports:
        - containerPort: 8000 # http
        - containerPort: 9000 # grpc
```

### Service - für beide Protokolle, wenn in der Demo auch nur HTTP verwendet wird

```yaml
apiVersion: v1
kind: Service
metadata:
  name: extauth
  namespace: extauth
spec:
  ports:
  - name: http
    port: 8000
    protocol: TCP
    targetPort: 8000
  - name: grpc
    port: 9000
    protocol: TCP
    targetPort: 9000
  selector:
    app: extauth
  sessionAffinity: None
  type: ClusterIP
```

### Provider in ISTIO bekannt machen

Der Provider muss im ISTIO-System bekannt gemacht werden, dazu muss die `configmap` angepasst werden. Das doppelte *extauth* kommt vom gleichlautenden Namespace. 
>Wichtig ist, dass sowohl die Header die an den Service gesendet werden sollen, als auch die die der Service wieder in Richtung Upstreams schicken darf definiert sind!

`kubectl edit configmap istio -n istio-system`

```yaml
extensionProviders:
    - name: "sample-ext-authz-http"
      envoyExtAuthzHttp:
        service: "extauth.extauth.svc.cluster.local" # namespace.service.
        port: "8000"
        includeRequestHeadersInCheck: ["x-ext-authz","x-request-id","x-b3-sampled","x-b3-spanid", "x-b3-traceid"]
        headersToUpstreamOnAllow: ["authorization", "x-auth-request-access-token"] # Header -> Service
```

Danch durchstarten von ISTIO

`kubectl rollout restart deployment/istiod -n istio-system`

## Policies einrichten

Die Zugriffsbeschränkung erfolgt bei Istio über den Envoy-Proxy, der wiederum wird durch verschiedenen ISTIO-CRDs konfiguriert.

### Zugriffe müssen über ein JWT Token verfügen das von der Audience *demoapi.rebelofbavaria.de* ausgestellt wurde

#### [Request-Authentication](https://istio.io/latest/docs/reference/config/security/request_authentication/)

```yaml
apiVersion: security.istio.io/v1beta1
kind: RequestAuthentication
metadata:
  name: gocoaster
  namespace: gocoaster
spec:
  selector:
    matchLabels:
      app: gocoaster       
  jwtRules: # wo sind die JSON Web Key Sets um die Signatur des JWT checken zu können
    - issuer: "https://dev-vdt9zz3q.us.auth0.com/"
      jwksUri: "https://dev-vdt9zz3q.us.auth0.com/.well-known/jwks.json"      
      audiences: # Zugriff diese Ressource(n)
        - "demoapi.rebelofbavaria.de"
      forwardOriginalToken: true
```

### Zugriffe auf Basis eines Claims - hier den scope - beschränken

#### [Authorisation-Policy](https://istio.io/latest/docs/reference/config/security/authorization-policy/)

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
 name: gocoaster
 namespace: gocoaster
spec:
 selector:
   matchLabels:
     app: gocoaster
 action: ALLOW
 rules:
 - from:
   - source:
      requestPrincipals: ["*"] # Benötigt gültiges JW-Token
   to:
   - operation:
      methods: ["GET","POST"] # für GET + POST reicht ...
   when:
   - key: request.auth.claims[scope]
     values:
       - "read:sample" # ... read
 - from:
   - source:    
      requestPrincipals: ["*"]
   to:
   - operation:
      methods: ["DELETE"] # für DELETE ...
   when:
   - key: request.auth.claims[scope]
     values:
       - "write:sample"  # ... muss es write sein!
```

## Externen AuthZ Provider in der Anwendung verwenden

Der Provider muss mit Hilfe einer Policy am jeweiligen Service aktiviert werden. In diesem Beispiel an einem Container mit dem Label *gocoaster* und zwar nur für die Methode `DELETE`

>Demo Use-Case: Der Anwender hat normalerweise nur ein lesendes Token, soll aber über den externen AuthZ-Provider eine Token mit Löschberechtigung erlangen können.

### Policy

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
 name: gocoaster-custom 
 namespace: gocoaster
spec:
 selector:
   matchLabels:
     app: gocoaster
 action: CUSTOM
 provider:
   name: sample-ext-authz-http
 rules:
 - to:
   - operation:
      methods: ["DELETE"]
```

Nachdem diese Policy deployed ist wird der Versuch einen Datensatz zu löschen fehlschlagen.

```bash
curl -X DELETE http://kom.io:80/mem/coasters/id850 -H "Accept: application/json"
```

erzeugt folgende Fehlermeldung:

```bash
denied by ext_authz for not found header `x-ext-authz: allow` in the request
```

weil der geforderte Header nicht angegeben wurde.

```bash
curl -X DELETE http://kom.io:80/mem/coasters/id850 -H "Accept: application/json" -H "x-ext-authz: allow"
```

führt zum löschen des Datensatzes. Im Log des extauth kann nachvollzogen werden das sich ein Token beschafft wurde. (Log natürlich nur für Demobetrieb)

```log
2022/07/19 12:15:00 [HTTP][allowed]: DELETE kom.io/mem/coasters/id813, headers: map[Content-Length:[0] X-B3-Parentspanid:[1fa4d650284b9d06] X-B3-Sampled:[1] X-B3-Spanid:[0956fcd7352d666c] X-B3-Traceid:[ba2000c66684e6b7bc052f5f3f0a59d8] X-Envoy-Expected-Rq-Timeout-Ms:[600000] X-Envoy-Internal:[true] X-Ext-Authz:[allow] X-Forwarded-For:[10.1.8.184] X-Request-Id:[e4ba0e7b-4966-938e-9daa-e7bcd8d79e0b]], body: []
2022/07/19 12:15:01 Token:= eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjNxdFFwY1ZZN0RlRXhOamJBV095cCJ9.eyJpc3MiOiJodHRwczovL2Rldi12ZHQ5enozcS51cy5hdXRoMC5jb20vIiwic3ViIjoiT0VQcDdmazdSa0s4V1ozVlA1N2RPVWVqNnpSdHZYMHpAY2xpZW50cyIsImF1ZCI6ImRlbW9hcGkucmViZWxvZmJhdmFyaWEuZGUiLCJpYXQiOjE2NTgyMzI5MDEsImV4cCI6MTY1ODMxOTMwMSwiYXpwIjoiT0VQcDdmazdSa0s4V1ozVlA1N2RPVWVqNnpSdHZYMHoiLCJzY29wZSI6InJlYWQ6c2FtcGxlIHdyaXRlOnNhbXBsZSIsImd0eSI6ImNsaWVudC1jcmVkZW50aWFscyIsInBlcm1pc3Npb25zIjpbInJlYWQ6c2FtcGxlIiwid3JpdGU6c2FtcGxlIl19.N6kWW44znpE6tTY0PznNaaYh4gpsSgVOF56rP8o0Wq-nzqVmh5I1KdfdTrYfL9MY_Doq8l_sBZ1PE8z3priFkJ8RTK3x5vXSQWczPbi-s5T70UQuCxjWKhEQOIuouaLuRyrXsXK8T3BuS_M5tWskJaUGWeO0UHQNJdfepPoUU8bkaCRzLB8-SW3XPno6pTshxeCbZXYtzFCw6S399NxKe_EFscaNA8YBQao5jXZYJe7-u9UE_eOBnI2kb3cW_MqJSicV0Fnm_BPTpdy0TmsTCTqowvFEI_vjZNNSiUW_WW3v_rF3csoXluAJxedRnVSfUK07yCfFiUrxVQ5KToRDUQ
```

Das Token im Klartext:

```json
{
  "iss": "https://dev-vdt9zz3q.us.auth0.com/",
  "sub": "OEPp7fk7RkK8WZ3VP57dOUej6zRtvX0z@clients",
  "aud": "demoapi.rebelofbavaria.de",
  "iat": 1658232901,
  "exp": 1658319301,
  "azp": "OEPp7fk7RkK8WZ3VP57dOUej6zRtvX0z",
  "scope": "read:sample write:sample",
  "gty": "client-credentials",
  "permissions": [
    "read:sample",
    "write:sample"
  ]
}
```

## ToDos - Ziel 1. Ablösen des Proxies , Ziel 2. Vereinfachung **silent login**

1. Was ist mit Caching - wenn alle Zugriffe auf den ExtAuth gehen, haben wir da nicht das nächste Performanceproblem? Könnte hier ein Redis helfen
2. Wie skaliert der ExtAuth-Service horizontal?
3. Entscheiden ob JWT Auswertung nur in ISTIO oder aber auch zuätzlich nochmal in der Anwendung
4. Fail fast, fail cheap - frühstmögliche Entscheidung das ein Call abgelehnt werden muss, wird für die CUSTOM Zugriffe ggf. nicht funktionieren, da diese immer ausgewertet werden. (Kann da eine geschickte d)
5. Ist die Reihenfolge der Auswertung der Filter evtl. konfigurierbar?
6. Was machen wir noch alles im aktuellen Proxy UND muss das an der Stelle gemacht werden?
7. Was ist mit dem Thema kein Token in den Browser sondern nur Session? --> Karsten nochmal dazu einladen und Std. der Dinge klären?!
8. ....
