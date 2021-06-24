package main

func OSHookConfig(cfgJSON []byte) []byte {
	return cfgJSON
}

func main() {
	serverLoop(startV2rayConfigRunner())
}
