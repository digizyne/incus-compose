package main

type DockerServiceVolume struct {
	Type string `json:"type"`
	Source string `json:"source"`
	Target string `json:"target"`
	Read_Only bool `json:"read_only"`
}

type DockerVolume struct {
	External bool `json:"external"`
}

type Volume interface {}

type DockerComposeService struct {
	Image string `json:"image"`
	Ports []string `json:"ports"`
	Volumes []Volume `json:"volumes"`
}

type DockerCompose struct {
	Services map[string]DockerComposeService `json:"services"`
	Volumes map[string]DockerVolume `json:"volumes"`
}