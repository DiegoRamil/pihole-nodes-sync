package client

import (
	"strconv"

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

	backups.RestoreBackupInChilds(client, backup)
}
