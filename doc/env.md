# Gutes Tooling für K8s

Um mit Kubernetes (K8s) zu arbeiten empfiehlt es sich eine effiziente Toolumgebung zu schaffen. Das fängt mit der Frage des OS an. Am einfachsten wird das ganze mit einem Linux Derivat, auch unter Windows sollte man sich das **WSL2** (Windows Subsystem für Linux) installieren. Denn auch wenn es für eingefleischte Windows-Anwender hart wird, der Weg lohnt sich.

>Natürlich geht auch fast alles mit Windows und/oder PowerShell. Aber auch der Weg ist hart ;-)

## Visual Studio Code

Für mich ist Visual Studio Code seit langen einer meiner Lieblingseditoren, und auch die Integration von K8s hat MS aus meiner Sicht gut gelöst. Angenehm ist auch das ich in der Wahl meines Terminal-Dialektes (DOS/PS/WSL) frei bin.

Sinnvolle Plugins:

* Kubernetes
* Kubernetes Support
* Azure Kubernetes Service
* Remote WSL
* Remote SSH
* Checkov
* Docker
* YAML

## Windows Terminal

Die Flexibilität mit verschiedenen Terminal-Derivaten zu arbeiten wird von Windows mit dem Tool **Windows Terminal** unterstützt.

[Download Windows Terminal](https://apps.microsoft.com/store/detail/windows-terminal/9N0DX20HK701?hl=de-de&gl=DE)

### Terminal einrichten

Es gibt einige Tools die helfen effizienter mit den K8s Clustern zu arbeiten.

tbd: **kommt noch**

## kubectl - one tool to rule them all

Zur Kommunikation mit der K8s API wird das Tool `kubectl` benötigt. Die Version von kubectl sollte nicht zu weit von der verwendeten K8s Version abweiche, da sonst ggf. nicht alle Features unterstützt werden.

[Download kubectl](https://kubernetes.io/de/docs/tasks/tools/install-kubectl/)

## KREW - Pluginmanager für kubectl

Auch kubectl lässt sich einfach über ein Pluginsystem erweitern. Zur Verwaltung der Plugins empfiehlt es sich Krew einzusetzen.

[Download Krew](https://krew.sigs.k8s.io/)

Sinnvolle Plugins:

* [oidc-login](https://github.com/int128/kubelogin) Ermöglich Login über einen OIDC Provider bsw. Keycloak
* [access-matrix](https://github.com/corneliusweig/rakkess) Anzeigen der RBAC Zugriffsmatrix
* [cert-manager](https://github.com/cert-manager/cert-manager) Erleichtert das arbeiten mit dem cert-manager
* [ctx](https://github.com/ahmetb/kubectx) Schnelles wechseln der K8s Cluster
* [ns](https://github.com/ahmetb/kubectx) Schnelles wechseln der Namespaces innerhalb einens Clusters
* [neat](https://github.com/itaysk/kubectl-neat) Entfernt "unnötige Informationen / Annotationen" aus Ressourcedefinitionen
