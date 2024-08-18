package main

type DockerComposeService struct {
	Image string `json:"image"`
	Ports []string `json:"ports"`
}

type DockerCompose struct {
	Services map[string]DockerComposeService
}