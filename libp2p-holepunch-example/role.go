package main

type Role string

const (
	Server   Role = "Server"
	NAT      Role = "NAT"
	Listener Role = "Listener"
	Dialer   Role = "Dialer"
)
