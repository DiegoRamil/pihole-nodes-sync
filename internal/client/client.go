package client

import (
	"os"
	"strconv"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/backups"
)

func SyncBetweenNodes() {
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		timeout = 30
	}
	client := CreateHttpClient(timeout)
	backup := backups.CreateBackup(client)

	backups.RestoreBackupInChilds(client, backup)
}
