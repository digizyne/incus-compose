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
	Read_Only bool `json:"read_only"`
}

type IncusVolume struct {
	External bool `json:"external"`
}



type IncusComposeService struct {
	Image string `json:"image"`
	Container_Name string `json:"container_name,omitempty"`
	Devices IncusDevices `json:"devices"`
	Volumes []IncusServiceVolume `json:"volumes"`
	Environment []string `json:"environment,omitempty"`
}

type IncusCompose struct {
	Services map[string]IncusComposeService `json:"services"`
	Volumes map[string]IncusVolume `json:"volumes"`
}