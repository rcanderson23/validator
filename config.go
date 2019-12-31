package main

type Config struct {
	Deployment *Fields `yaml:"Deployment,omitempty"`
	ReplicaSet *Fields `yaml:"ReplicaSet,omitempty"`
	Pod        *Fields `yaml:"Pod,omitempty"`
}

type Fields struct {
	Labels map[string]string `yaml:"labels,omitempty"`
	Image  string            `yaml:"image,omitempty"`
}
