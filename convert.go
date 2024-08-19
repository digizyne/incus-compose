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

func convertTopLevelVolumes(volumes map[string]DockerVolume) map[string]IncusVolume {
	incusVolumes := make(map[string]IncusVolume)
	for key := range volumes {
		incusVolumes[key] = IncusVolume{
			External: volumes[key].External,
		}
	}
	return incusVolumes
}

func convertServiceVolume(volume string, topLevelVolumes map[string]DockerVolume) IncusServiceVolume {
	volumePairs := strings.Split(volume, ":")

	mountType := "bind"
	if _, ok := topLevelVolumes[volumePairs[0]]; ok {
		mountType = "volume"
	}

	return IncusServiceVolume{
		Type: mountType,
		Source: volumePairs[0],
		Target: volumePairs[1],
	}
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

		serviceVolumes := []IncusServiceVolume{}
		for _, volume := range dockerCompose.Services[key].Volumes {
			serviceVolumes = append(serviceVolumes, convertServiceVolume(volume, dockerCompose.Volumes))
		}

		incusComposeService := IncusComposeService{
			Image: fmt.Sprintf("docker:%s", dockerCompose.Services[key].Image),
			Devices: IncusDevices{
				Proxies: proxies,
			},
			Volumes: serviceVolumes,
		}

		incusCompose.Services[key] = incusComposeService
	}

	incusCompose.Volumes = convertTopLevelVolumes(dockerCompose.Volumes)

	incusComposeYaml, err := yaml.Marshal(incusCompose)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(incusComposeYaml))

	os.WriteFile("incus-compose.yaml", incusComposeYaml, 0644)
}