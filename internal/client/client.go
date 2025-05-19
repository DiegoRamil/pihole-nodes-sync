package client

import (
	"strconv"
	"sync"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/backups"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/shared"
)

func SyncBetweenNodes() {
	timeout, err := strconv.Atoi(shared.RetrieveEnvVar("TIMEOUT"))
	if err != nil {
		timeout = 30
	}
	client := CreateHttpClient(timeout)
	backup := backups.CreateBackup(client)

	// Create a goroutine for each child node
	var wg sync.WaitGroup
	nodes := backups.GetChildNodes()
	
	for _, node := range nodes {
		wg.Add(1)
		go func(n backups.ChildNode) {
			defer wg.Done()
			backups.RestoreBackupInChild(client, backup, n)
		}(node)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}
