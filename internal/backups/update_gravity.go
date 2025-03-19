package backups

import (
	"fmt"
	"net/http"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/shared"
)

func UpdateGravity(client *http.Client, baseUrl string, token string) {
	apiUrl := shared.ConcatBaseUrlAndUri(baseUrl, "/api/action/gravity")

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		fmt.Printf("Error trying to update gravity in the child node: %s\n", err)
	}
	req.Header.Add("X-FTL-SID", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error trying to update gravity in the child node: %s\n", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error trying to update gravity in the child node: %s\n", resp.Status)
	} else {
		fmt.Println("Gravity updated successfully")
	}
	DeauthorizeToken(client, token, baseUrl)
}
