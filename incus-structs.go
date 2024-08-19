package main

type IncusProxy struct {
	Listen string `json:"listen"`
	Connect string `json:"connect"`
}

type IncusDevices struct {
	Proxies []IncusProxy `json:"proxies"`
}

type IncusServiceVolume struct {
	Type string `json:"type"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type IncusVolume struct {
	External bool `json:"external"`
}



type IncusComposeService struct {
	Image string `json:"image"`
	Devices IncusDevices `json:"devices"`
	Volumes []IncusServiceVolume `json:"volumes"`
}

type IncusCompose struct {
	Services map[string]IncusComposeService `json:"services"`
	Volumes map[string]IncusVolume `json:"volumes"`
}