package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"sigs.k8s.io/yaml"
)

func createService(incusComposeService IncusComposeService, serviceName string) {
	configurationOptions := []string{
		"create",
		incusComposeService.Image,
		serviceName,
	}
	if len(incusComposeService.Environment) > 0 {
		for _, env := range incusComposeService.Environment {
			configurationOptions = append(configurationOptions, fmt.Sprintf("-c environment.%s", env))
		}
	}

	fmt.Println(configurationOptions)

	fmt.Println(Blue, "*** creating service", serviceName, "***")
	createCommand := exec.Command("incus", "create", incusComposeService.Image, serviceName)
	createOutput, err := createCommand.CombinedOutput()
	if err != nil {
		fmt.Println(Red, err)
		return
	}
	fmt.Println(Green, string(createOutput))
}

func setEnvVar(envVar string, serviceName string) {
	cmd := exec.Command("incus", "config", "set", serviceName, fmt.Sprintf("environment.%s", envVar))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(Red, err)
		return
	}
	fmt.Println(Green, string(output))
}

func createProxy(incusProxy IncusProxy, serviceName string) {
	proxyPortString := strings.Split(incusProxy.Connect, ":")
	proxyPort := proxyPortString[len(proxyPortString)-1]
	fmt.Println(Blue, "*** creating proxy", incusProxy.Listen, incusProxy.Connect, "***")
	proxyCommand := exec.Command("incus", "config", "device", "add", serviceName, fmt.Sprintf("%s-proxy-%s", serviceName, proxyPort), "proxy", fmt.Sprintf("listen=%s", incusProxy.Listen), fmt.Sprintf("connect=%s", incusProxy.Connect))
	proxyOutput, err := proxyCommand.CombinedOutput()
	if err != nil {
		fmt.Println(Red, fmt.Sprint(err))
		return
	}
	fmt.Println(Green, string(proxyOutput))
}

func createVolume(volumeName string) {
	fmt.Println(Blue, "*** creating volume", volumeName, "***")
	volumeCommand := exec.Command("incus", "storage", "volume", "create", "default", volumeName)
	volumeOutput, err := volumeCommand.CombinedOutput()
	if err != nil {
		fmt.Println(Red, err)
		return
	}
	fmt.Println(Green, string(volumeOutput))
}

func mountVolume(volume IncusServiceVolume, serviceName string) {
	fmt.Println(Blue, "*** mounting volume", volume.Source, "to", volume.Target, "***")
	mountCommand := exec.Command("incus", "config", "device", "add", serviceName, fmt.Sprintf("%s-%s-disk", serviceName, volume.Source), "disk", "pool=default", fmt.Sprintf("source=%s", volume.Source), fmt.Sprintf("path=%s", volume.Target), fmt.Sprintf("readonly=%t", volume.Read_Only))
	mountOutput, err := mountCommand.CombinedOutput()
	if err != nil {
		fmt.Println(Red, err)
		return
	}
	fmt.Println(Green, string(mountOutput))
}

func startService(serviceName string) {
	fmt.Println(Blue, "*** starting service", serviceName, "***")
	launchCommand := exec.Command("incus", "start", serviceName)
	launchOutput, err := launchCommand.CombinedOutput()
	if err != nil {
		fmt.Println(Red, err)
		return
	}
	fmt.Println(string(launchOutput))
	fmt.Println(Green, "*** service", serviceName, "started ***")
}

func up() {
	data, err := os.ReadFile("incus-compose.yaml")
	if err != nil {
		fmt.Println(Red, err)
		return
	}

	var incusCompose IncusCompose
	err = yaml.Unmarshal(data, &incusCompose)
	if err != nil {
		fmt.Println(Red, err)
		return
	}

	//* Create volumes
	for key := range incusCompose.Volumes {
		if !incusCompose.Volumes[key].External {
			createVolume(key)
		}
	}

	//* Create services
	for key := range incusCompose.Services {
		var serviceName string
		if incusCompose.Services[key].Container_Name != "" {
			serviceName = incusCompose.Services[key].Container_Name
		} else {
			serviceName = strings.Split(incusCompose.Services[key].Image, ":")[1]
		}

		createService(incusCompose.Services[key], serviceName)

		fmt.Println(Blue, "*** configuring service", key, "***")

		fmt.Println(Blue, "*** setting env vars for", key, "***")
		for _, envVar := range incusCompose.Services[key].Environment {
			setEnvVar(envVar, serviceName)
		}

		for _, proxy := range incusCompose.Services[key].Devices.Proxies {
			createProxy(proxy, serviceName)
		}

		for _, volume := range incusCompose.Services[key].Volumes {
			if volume.Type == "volume" {
				mountVolume(volume, serviceName)
			}
		}

		startService(serviceName)
	}
}
