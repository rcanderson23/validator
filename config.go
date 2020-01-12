package main

type Config struct {
	ValidatorSpec *ValidatorSpec `yaml:"validatorSpec"`
	TlsCert       string         `yaml:"tlsCert"`
	TlsKey        string         `yaml:"tlsKey"`
}
type ValidatorSpec struct {
	Pod     *PodFields     `yaml:"Pod"`
	Service *ServiceFields `yaml:"Service"`
}

type PodFields struct {
	Labels map[string]string `yaml:"labels"`
	Image  string            `yaml:"image"`
}

type ServiceFields struct {
	Labels              map[string]string `yaml:"labels"`
	DisableLoadBalancer bool              `yaml:"disableLoadBalancer"`
}
