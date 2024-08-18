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
	fmt.Println("*** creating service", serviceName, "***")
	createCommand := exec.Command("incus", "create", incusComposeService.Image, serviceName)
	createOutput, err := createCommand.CombinedOutput()
	if (err != nil) {
		fmt.Println(err)
		return
	}
	fmt.Println(string(createOutput))
}

func createProxy(incusProxy IncusProxy, serviceName string) {
	proxyPortString := strings.Split(incusProxy.Connect, ":")
	proxyPort := proxyPortString[len(proxyPortString)-1]
	fmt.Println("*** creating proxy", incusProxy.Listen, incusProxy.Connect, "***")
	proxyCommand := exec.Command("incus", "config", "device", "add", serviceName, fmt.Sprintf("%s-proxy-%s", serviceName, proxyPort), "proxy", fmt.Sprintf("listen=%s", incusProxy.Listen), fmt.Sprintf("connect=%s", incusProxy.Connect))
	proxyOutput, err := proxyCommand.CombinedOutput()
	if (err != nil) {
		fmt.Println(fmt.Sprint(err))
		return
	}
	fmt.Println(string(proxyOutput))
}

func startService(serviceName string) {
	fmt.Println("*** starting service", serviceName, "***")
	launchCommand := exec.Command("incus", "start", serviceName)
	launchOutput, err := launchCommand.CombinedOutput()
	if (err != nil) {
		fmt.Println(err)
		return
	}
	fmt.Println(string(launchOutput))
	fmt.Println("*** service", serviceName, "started ***")
}

func up() {
	data, err := os.ReadFile("incus-compose.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	var incusCompose IncusCompose
	err = yaml.Unmarshal(data, &incusCompose)
	if err != nil {
		fmt.Println(err)
		return
	}

	for key := range incusCompose.Services {
		createService(incusCompose.Services[key])

		fmt.Println("*** configuring service ***")
		for _, proxy := range incusCompose.Services[key].Proxies {
			createProxy(proxy, key)
		}

		startService(key)
	}
}