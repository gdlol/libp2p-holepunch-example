package main

import (
	"context"
	"log"
	"os"
	"os/exec"
)

func run(cmd string, args ...string) {
	command := exec.Command(cmd, args...)
	log.Println(command.String())
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		panic(err)
	}
}

func configureGateway(gatewayAddress string) {
	run("ip", "route", "replace", "default", "via", gatewayAddress)
}

func runNAT(ctx context.Context, natRange, externalAddress, clientAddress string) {
	run("iptables", "--append", "FORWARD", "--jump", "ACCEPT")
	run("iptables", "--append", "INPUT", "--jump", "DROP")
	run("iptables", "--append", "OUTPUT", "--jump", "DROP")
	run("iptables",
		"--table", "nat",
		"--append", "POSTROUTING",
		"--source", natRange,
		"--jump", "SNAT",
		"--to-source", externalAddress)
	// run("iptables",
	// 	"--table", "nat",
	// 	"--append", "PREROUTING",
	// 	"--destination", externalAddress,
	// 	"--jump", "DNAT",
	// 	"--to-destination", clientAddress)
	<-ctx.Done()
}
