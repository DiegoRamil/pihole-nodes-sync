package backups

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/backups/model"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/deserializers"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/shared"
)

func AuthorizeWithPihole(baseUrl string, pwd string, client *http.Client) string {
	fmt.Printf("Doing authorization in %s...\n", baseUrl)
	path := shared.ConcatBaseUrlAndUri(baseUrl, "/api/auth")

	resp, err := client.Post(path, "application/json", strings.NewReader(`{"password": "`+pwd+`"}`))
	if err != nil {
		fmt.Printf("Error in authorization: %s\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("authorization failed with status: %s", resp.Status)
	}

	authorization := &model.AuthorizationResponse{}
	if err := deserializers.JsonDeserialize(resp.Body, authorization); err != nil {
		fmt.Printf("Error deserializing the body: %s\n", err)
	}
	return authorization.Session.Sid
}

func DeauthorizeToken(client *http.Client, token string, baseUrl string) {
	fmt.Printf("Removing auth token in %s...\n", baseUrl)
	path := shared.ConcatBaseUrlAndUri(baseUrl, "/api/auth")
	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		fmt.Printf("Error when trying to remove the token: %s\n", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-FTL-SID", token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error when trying to remove the token: %s\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		fmt.Printf("removing auth token failed with status: %s", resp.Status)
	}
	fmt.Printf("Auth token removed in %s...\n", baseUrl)
}
