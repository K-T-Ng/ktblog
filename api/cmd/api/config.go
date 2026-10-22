package main

const (
	keyAddr     = "ADDR"
	defaultAddr = "127.0.0.1:8080"
)

type config struct {
	Addr string
}

func loadConfig(getEnv func(string) string) (config, error) {
	addr := getEnv(keyAddr)
	if addr == "" {
		addr = defaultAddr
	}
	return config{Addr: addr}, nil
}
