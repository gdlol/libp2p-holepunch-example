package main

type Config struct {
	Role            Role   `envconfig:"ROLE"`
	ServerAddress   string `envconfig:"SERVER_ADDRESS"`
	ListenerAddress string `envconfig:"LISTENER_ADDRESS"`
	ExternalAddress string `envconfig:"EXTERNAL_ADDRESS"`
	NATRange        string `envconfig:"NAT_RANGE"`
	GatewayAddress  string `envconfig:"GATEWAY_ADDRESS"`
	ClientAddress   string `envconfig:"CLIENT_ADDRESS"`
}

const defaultPort = 80
