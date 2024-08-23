package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"sigs.k8s.io/yaml"
)

func createService(incusComposeService IncusComposeService) {
	serviceName := strings.Split(incusComposeService.Image, ":")[1]
	fmt.Println(Blue, "*** creating service", serviceName, "***")
	createCommand := exec.Command("incus", "create", incusComposeService.Image, serviceName)
	createOutput, err := createCommand.CombinedOutput()
	if (err != nil) {
		fmt.Println(Red, err)
		return
	}
	fmt.Println(Green, string(createOutput))
}

func createProxy(incusProxy IncusProxy, serviceName string) {
	proxyPortString := strings.Split(incusProxy.Connect, ":")
	proxyPort := proxyPortString[len(proxyPortString)-1]
	fmt.Println(Blue, "*** creating proxy", incusProxy.Listen, incusProxy.Connect, "***")
	proxyCommand := exec.Command("incus", "config", "device", "add", serviceName, fmt.Sprintf("%s-proxy-%s", serviceName, proxyPort), "proxy", fmt.Sprintf("listen=%s", incusProxy.Listen), fmt.Sprintf("connect=%s", incusProxy.Connect))
	proxyOutput, err := proxyCommand.CombinedOutput()
	if (err != nil) {
		fmt.Println(Red, fmt.Sprint(err))
		return
	}
	fmt.Println(Green, string(proxyOutput))
}

func createVolume(volumeName string) {
	fmt.Println(Blue, "*** creating volume", volumeName, "***")
	volumeCommand := exec.Command("incus", "storage", "volume", "create", "default", volumeName)
	volumeOutput, err := volumeCommand.CombinedOutput()
	if (err != nil) {
		fmt.Println(Red, err)
		return
	}
	fmt.Println(Green, string(volumeOutput))
}

func mountVolume(volume IncusServiceVolume, serviceName string) {
	fmt.Println(Blue, "*** mounting volume", volume.Source, "to", volume.Target, "***")
	mountCommand := exec.Command("incus", "storage", "volume", "attach", "default", volume.Source, serviceName, volume.Target)
	mountOutput, err := mountCommand.CombinedOutput()
	if (err != nil) {
		fmt.Println(Red, err)
		return
	}
	fmt.Println(Green, string(mountOutput))
}

func bindMount(volume IncusServiceVolume, serviceName string) {
	fmt.Println(Blue, "*** bind mounting directory", volume.Source, "to", volume.Target, "***")
	bindMountCommand := exec.Command("incus", "file", "mount", fmt.Sprintf("%s/%s", serviceName, volume.Target), volume.Source)
	// var stdout, stderr bytes.Buffer
	bindMountCommand.Stdout = nil
	bindMountCommand.Stderr = nil
	bindMountCommand.Stdin = nil
	// go func() {
		err2 := bindMountCommand.Run()
		if (err2 != nil) {
			fmt.Println(Red, err2)
			return
		}
	// }()
	// fmt.Println(Green, string(bindMountOutput))
}

func startService(serviceName string) {
	fmt.Println(Blue, "*** starting service", serviceName, "***")
	launchCommand := exec.Command("incus", "start", serviceName)
	launchOutput, err := launchCommand.CombinedOutput()
	if (err != nil) {
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
		if (!incusCompose.Volumes[key].External) {
			createVolume(key)
		}
	}

	//* Create services
	for key := range incusCompose.Services {
		createService(incusCompose.Services[key])

		fmt.Println(Blue, "*** configuring service", key, "***")
		for _, proxy := range incusCompose.Services[key].Devices.Proxies {
			createProxy(proxy, key)
		}

		for _, volume := range incusCompose.Services[key].Volumes {
			if (volume.Type == "volume") {
				mountVolume(volume, key)
			}
			if (volume.Type == "bind") {
				bindMount(volume, key)
			}
		}

		startService(key)
	}
}