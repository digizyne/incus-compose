package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
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

func convertServiceVolumeShorthand(volume string, topLevelVolumes map[string]DockerVolume) (IncusServiceVolume, error) {
	volumePairs := strings.Split(volume, ":")
	if _, ok := topLevelVolumes[volumePairs[0]]; ok {
		incuseSrviceVolume := IncusServiceVolume{
			Type: "volume",
			Source: volumePairs[0],
			Target: volumePairs[1],
		}
		if (len(volumePairs) == 3) {
			incuseSrviceVolume.Read_Only = strings.Contains(volumePairs[2], "ro")
		}

		return incuseSrviceVolume, nil
	}

	return IncusServiceVolume{}, errors.New("bind mounts are not yet supported in incus-compose")
}

func convertServiceVolumeLongform(volume DockerServiceVolume, topLevelVolumes map[string]DockerVolume) (IncusServiceVolume, error) {
	if _, ok := topLevelVolumes[volume.Source]; ok {
		if (volume.Type == "volume") {
			incusCompose := IncusServiceVolume(volume)
			return incusCompose, nil
		}
	}

	return IncusServiceVolume{}, errors.New("bind mounts are not yet supported in incus-compose")
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
			if (reflect.TypeOf(volume).Kind() == reflect.String) {
				newVolume, err := convertServiceVolumeShorthand(volume.(string), dockerCompose.Volumes)
				if (err != nil) {
					continue
				}
				serviceVolumes = append(serviceVolumes, newVolume)
			} else {
				var coercedVolume DockerServiceVolume
				err := mapstructure.Decode(volume, &coercedVolume)
				if (err != nil) {
					continue
				}
				newVolume, err := convertServiceVolumeLongform(coercedVolume, dockerCompose.Volumes)
				if (err != nil) {
					continue
				}
				serviceVolumes = append(serviceVolumes, newVolume)
			}
		}

		incusComposeService := IncusComposeService{
			Image: fmt.Sprintf("docker:%s", dockerCompose.Services[key].Image),
			Devices: IncusDevices{
				Proxies: proxies,
			},
			Volumes: serviceVolumes,
		}

		if (dockerCompose.Services[key].Container_Name != "") {
			incusComposeService.Container_Name = dockerCompose.Services[key].Container_Name
		}

		if (len(dockerCompose.Services[key].Environment) > 0) {
			incusComposeService.Environment = dockerCompose.Services[key].Environment
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