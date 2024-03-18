package config

import (
	"os"
	"strconv"
)

type Config struct {
	StoragePath string

	CounterFlushIntervalSeconds int
	CounterWindowSeconds        int

	ServerAddress string
}

const DEFAULTPATH = "data.txt"
const DEFAULTFLUSHINTERVAL = 66
const DEFAULTWINDOW = 60
const DEFAULTADDRESS = ":8080"

func NewFromEnv() Config {
	c := Config{
		StoragePath:                 "data.txt",
		CounterFlushIntervalSeconds: DEFAULTFLUSHINTERVAL,
		CounterWindowSeconds:        DEFAULTWINDOW,
		ServerAddress:               DEFAULTADDRESS,
	}

	storagePath := os.Getenv("SC_STORAGEPATH")
	if storagePath != "" {
		c.StoragePath = storagePath
	}

	counterFlushIntervalSeconds := os.Getenv("SC_FLUSHINTERVALSECONDS")
	if counterFlushIntervalSeconds != "" {
		interval, err := strconv.Atoi(counterFlushIntervalSeconds)
		if err != nil {
			c.CounterFlushIntervalSeconds = interval
		}
	}

	counterWindowSeconds := os.Getenv("SC_WINDOWSIZESECONDS")
	if counterWindowSeconds != "" {
		window, err := strconv.Atoi(counterFlushIntervalSeconds)
		if err != nil {
			c.CounterWindowSeconds = window
		}
	}

	serverAddress := os.Getenv("SC_ADDRESS")
	if serverAddress != "" {
		c.ServerAddress = serverAddress
	}

	return c
}
