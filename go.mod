module github.com/joelanford/multicache-operator

go 1.15

require (
	github.com/go-logr/logr v0.4.0
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	sigs.k8s.io/controller-runtime v0.8.3
)

replace sigs.k8s.io/controller-runtime => github.com/joelanford/controller-runtime v0.9.0-alpha.1.0.20210426142113-4f6cd8b22fd6
