package main

type DockerServiceVolumeLongform struct {
	Type string `json:"type"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type DockerServiceVolume struct {
	// Longforms []DockerServiceVolumeLongform `json:"longforms"`
	// Shortforms []string `json:"shorthands"`
}

type DockerVolume struct {
	External bool `json:"external"`
}



type DockerComposeService struct {
	Image string `json:"image"`
	Ports []string `json:"ports"`
	Volumes []string `json:"volumes"`
}

type DockerCompose struct {
	Services map[string]DockerComposeService `json:"services"`
	Volumes map[string]DockerVolume `json:"volumes"`
}