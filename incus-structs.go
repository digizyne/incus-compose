package main

type IncusProxy struct {
	Listen string `json:"listen"`
	Connect string `json:"connect"`
}

type IncusComposeService struct {
	Image string `json:"image"`
	Proxies []IncusProxy `json:"proxies"`
}

type IncusCompose struct {
	Services map[string]IncusComposeService `json:"services"`
}