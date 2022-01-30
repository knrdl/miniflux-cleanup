package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const EntryProcessingLimit = 5000

func hasEnv(varname string) bool {
	return len(os.Getenv(varname)) > 0
}

func getMinifluxUrl() string {
	return os.Getenv("MINIFLUX_URL")
}

func getUsername(request *http.Request) string {
	if hasEnv("AUTH_PROXY_HEADER") {
		val := request.Header.Get(os.Getenv("AUTH_PROXY_HEADER"))
		if len(val) > 0 {
			return val
		} else {
			return os.Getenv("DEFAULT_USERNAME")
		}
	} else {
		return os.Getenv("DEFAULT_USERNAME")
	}
}

func getCronInterval() (time.Duration, error) {
	dur, err := time.ParseDuration(os.Getenv("CRONJOB_INTERVAL"))
	if err != nil {
		return dur, err
	}
	if dur < 1*time.Second {
		return dur, fmt.Errorf("duration of CRONJOB_INTERVAL too short")
	}
	return dur, nil
}

func checkEnv() {
	if !hasEnv("MINIFLUX_URL") {
		log.Fatal("Environment Variable MINIFLUX_URL must be set!")
	}
	if !hasEnv("CRONJOB_INTERVAL") {
		log.Fatal("Environment Variable CRONJOB_INTERVAL must be set!")
	}
	if _, err := getCronInterval(); err != nil {
		log.Fatal(err)
	}
	if !hasEnv("AUTH_PROXY_HEADER") && !hasEnv("DEFAULT_USERNAME") {
		log.Fatal("Either Environment Variable AUTH_PROXY_HEADER or DEFAULT_USERNAME must be set!")
	}

}
