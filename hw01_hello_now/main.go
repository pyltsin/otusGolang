package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {

	now := time.Now()
	ntpNow, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalf("Error getting time: %s", err)
	}

	fmt.Printf("current time: %s\n", now.Round(0))
	fmt.Printf("exact time: %s\n", ntpNow.Round(0))
}
