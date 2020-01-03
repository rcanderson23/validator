package main

type Config struct {
	ValidatorSpec *ValidatorSpec `yaml:"validatorSpec"`
	TlsCert       string         `yaml:"tlsCert"`
	TlsKey        string         `yaml:"tlsKey"`
}
type ValidatorSpec struct {
	Deployment *Fields `yaml:"Deployment"`
	ReplicaSet *Fields `yaml:"ReplicaSet"`
	Pod        *Fields `yaml:"Pod"`
}

type Fields struct {
	Labels map[string]string `yaml:"labels"`
	Image  string            `yaml:"image"`
}
