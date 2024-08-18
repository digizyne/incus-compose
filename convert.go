package main

import (
	"fmt"
	"os"
	"strings"
	"sigs.k8s.io/yaml"
)

func convertPortsToProxies(ports []string) []IncusProxy {
	proxies:= []IncusProxy{}
		for _, port := range ports {
			portPairs := strings.Split(port, ":")
			proxy := IncusProxy{
				Listen: fmt.Sprintf("tcp:127.0.0.1:%s", portPairs[0]),
				Connect: fmt.Sprintf("tcp:0.0.0.0:%s", portPairs[1]),
			}
			proxies = append(proxies, proxy)
		}
	return proxies
		}

func convertDockerComposeToIncusCompose(inputFile string) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
	}

	var dockerCompose DockerCompose
	err = yaml.Unmarshal(data, &dockerCompose)
	if err != nil {
		fmt.Println(err)
		return
	}

	incusCompose := IncusCompose{
		Services: make(map[string]IncusComposeService),
	}

	for key := range dockerCompose.Services {

		proxies := convertPortsToProxies(dockerCompose.Services[key].Ports)

		incusComposeService := IncusComposeService{
			Image: fmt.Sprintf("docker:%s", dockerCompose.Services[key].Image),
			Proxies: proxies,
		}
		incusCompose.Services[key] = incusComposeService
	}

	incusComposeYaml, err := yaml.Marshal(incusCompose)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(incusComposeYaml))

	os.WriteFile("incus-compose.yaml", incusComposeYaml, 0644)
}