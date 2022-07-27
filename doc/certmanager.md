# Verwalten von Zertifikaten mit dem Cert-Manager

Jeder Traffic in den Cluster muss tls-Verschlüsselt sein, dazu benötigt der Cluster genauer der Ingress-Controller SSL-Zertifikate. Da diese i.a.R. eine Laufzeit haben, und die Anzahl der Cluster und Zertifikate schnell unübersichtlich wird, ist es sinnvoll diese automatisiert zu verwalten. Unter K8s eignet sich dazu der [*cert-manager*](https://cert-manager.io/docs/)

## Installation

Ich verwende dazu das **Krew** Plugin cert-manager. [Details zur installation](https://cert-manager.io/docs/installation/)

`kubectl cert-manager x install`

## Beispiel - self signed Zertifikat an ISTIO Ingress Controller anbinden

Cert-Manager wird über **CRDs** configuriert und gesteuert. Ein **Issuer** - nur im angelegten NS verwendbar - oder **ClusterIssuer** - Clusterweit verwendbar - sorgt dafür, das Zertifikate ausgestellt werden können. Der einfachste Fall ist ein Self Signed Issuer.

```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: ssgnd-issuer # Name des Issues auf den Bezug genommen wird
  namespace: istio-system # an den ISTIO NS gebunden also kein ClusterIssuer, da der Ingress Controller die Secrets im NS von ISTIO benötigt.
spec:
  selfSigned: {}
```

### Das eigentliche Zeritifikat anfordern

Der Issuer erstellt die Zertifikate für uns, welches Zertifikat wir benötigen wird in einem anderen **CRD** beschreiben. Das gewünschte Zertifikat wird im gewünschten **Secret** gespeichert.

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: cert-web-ssl
  namespace: istio-system # secret muss im ns von ingress-controller liegen 
spec:
  secretName: web-ssl # das ist der Name des Secrets, welches durch den cert-manager angelegt werden soll

  dnsNames:
  - kom.io # für welchen / welche DNS soll das Zertifikat verwendet werden
  issuerRef: # wer soll das Zertifikat für uns besorgen?
    name: ssgnd-issuer
    kind: Issuer

```

### Im Gateway bekannt machen

Der Ingress-Controller muss wissen welches Zertifikat bei welchem Hostnamen verwendet werden soll.

```yaml
kind: Gateway
apiVersion: networking.istio.io/v1alpha3
metadata:
  name: gocoaster-gateway
  namespace: gocoaster # im NS der Anwendung
spec:
  selector:
    istio: ingressgateway

  servers:
    - port:
        name: https
        number: 443
        protocol: https
      tls:
        mode: SIMPLE
        credentialName: web-ssl # das ist das Secret, was uns der cert-manager angelegt hat
      hosts:
        - kom.io  
```

# cert-manager und let´s encrypt

Funktioniert ziemlich gleich einzig der Issuer ist anders zu konfigurieren und es können mehr Parameter bei der Erstellung des Zertifikats angegeben werden. (Um ehrlich zu sein habe ich das beim ss nicht ausprobiert :-) )

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
  namespace: istio-system
spec:
  acme:
    # You must replace this email address with your own.
    # Let's Encrypt will use this to contact you about expiring
    # certificates, and issues related to your account.
    email: omarohn@gmail.com
    #server: https://acme-staging-v02.api.letsencrypt.org/directory # stage
    server: https://acme-v02.api.letsencrypt.org/directory # prod
    privateKeySecretRef:
      # Secret resource that will be used to store the account's private key.
      name: le-issuer-account-key
    # Add a single challenge solver, HTTP01 using nginx
    solvers:
    - http01:
        ingress:
          class: istio
```

und das Zertifikat

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: omarohn-domain-cert-prod
  namespace: istio-system
spec:
  secretName: web-ssl
  duration: 2160h # 90d
  renewBefore: 360h # 15d
  isCA: false
  privateKey:
    algorithm: RSA
    encoding: PKCS1
    size: 2048
  usages:
    - server auth
    - client auth
  dnsNames:
    - "omarohn.de"
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
    group: cert-manager.io
```
