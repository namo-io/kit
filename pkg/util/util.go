package util

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

var (
	serviceId = uuid.NewString()
)

func GetServiceId() string {
	return serviceId
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown"
	}

	return hostname
}

func WaitSignal() os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer signal.Stop(sigs)

	return <-sigs
}
