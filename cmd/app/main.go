package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/client"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/shared"
)

func main() {
	timer, err := strconv.Atoi(shared.RetrieveEnvVar("SYNC_HOURS"))
	if err != nil {
		panic(err)
	}

	for {
		client.SyncBetweenNodes()
		currentTime := time.Now().Local().Local().Add(time.Duration(timer) * time.Hour)
		fmt.Println("Sync completed. Next sync at:", currentTime.Format(time.RFC1123))
		time.Sleep(time.Hour * time.Duration(timer))
	}
}
