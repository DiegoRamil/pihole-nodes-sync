package backups_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/backups"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/backups/model"
	"github.com/stretchr/testify/assert"
)

func setupMockServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/auth":
			// Mock auth response
			json.NewEncoder(w).Encode(model.AuthorizationResponse{
				Session: model.Session{
					Valid:    true,
					Totp:     false,
					Sid:      "test_session",
					Csrf:     "test_csrf",
					Validity: 3600,
				},
				Took: 0.1,
			})
		case "/api/teleporter":
			// Mock teleporter response
			json.NewEncoder(w).Encode(map[string]interface{}{
				"processed": []string{"test"},
				"took":      0.5,
			})
		case "/api/gravity":
			// Mock gravity update response
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestGetChildNodes(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		childUrls      string
		childPasswords string
		shouldPanic    bool
		expected       []backups.ChildNode
	}{
		{
			name:           "Valid nodes",
			childUrls:      "http://pihole1:80,http://pihole2:80",
			childPasswords: "pass1,pass2",
			shouldPanic:    false,
			expected: []backups.ChildNode{
				{URL: "http://pihole1:80", Password: "pass1"},
				{URL: "http://pihole2:80", Password: "pass2"},
			},
		},
		{
			name:           "Empty values",
			childUrls:      "",
			childPasswords: "",
			shouldPanic:    true,
			expected:       nil,
		},
		{
			name:           "Mismatched lengths",
			childUrls:      "http://pihole1:80,http://pihole2:80",
			childPasswords: "pass1",
			shouldPanic:    false,
			expected: []backups.ChildNode{
				{URL: "http://pihole1:80", Password: "pass1"},
			},
		},
		{
			name:           "With spaces",
			childUrls:      " http://pihole1:80 , http://pihole2:80 ",
			childPasswords: " pass1 , pass2 ",
			shouldPanic:    false,
			expected: []backups.ChildNode{
				{URL: "http://pihole1:80", Password: "pass1"},
				{URL: "http://pihole2:80", Password: "pass2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("CHILD_URLS", tt.childUrls)
			os.Setenv("CHILD_PASSWORDS", tt.childPasswords)
			defer os.Unsetenv("CHILD_URLS")
			defer os.Unsetenv("CHILD_PASSWORDS")

			if tt.shouldPanic {
				assert.Panics(t, func() {
					backups.GetChildNodes()
				})
				return
			}

			// Get child nodes
			nodes := backups.GetChildNodes()

			// Assert
			assert.Equal(t, tt.expected, nodes)
		})
	}
}

func TestRestoreBackupInChild(t *testing.T) {
	// Create mock server
	server := setupMockServer(t)
	defer server.Close()

	// Set up environment variables
	os.Setenv("UPDATE_GRAVITY", "true")
	defer os.Unsetenv("UPDATE_GRAVITY")

	// Test cases
	tests := []struct {
		name     string
		node     backups.ChildNode
		backup   *model.BackupResponse
		expected *backups.BackupChildResponse
	}{
		{
			name: "Empty node",
			node: backups.ChildNode{},
			backup: &model.BackupResponse{
				Content: []byte("test"),
			},
			expected: nil,
		},
		{
			name: "Valid node",
			node: backups.ChildNode{
				URL:      server.URL,
				Password: "test_password",
			},
			backup: &model.BackupResponse{
				Content: []byte("test"),
			},
			expected: &backups.BackupChildResponse{
				Processed: []string{server.URL},
				Took:      0.5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test client
			client := server.Client()

			// Call the function
			result := backups.RestoreBackupInChild(client, tt.backup, tt.node)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRestoreBackupInChilds(t *testing.T) {
	// Create mock server
	server := setupMockServer(t)
	defer server.Close()

	// Set up environment variables
	os.Setenv("UPDATE_GRAVITY", "true")
	defer os.Unsetenv("UPDATE_GRAVITY")

	// Test cases
	tests := []struct {
		name           string
		childUrls      string
		childPasswords string
		shouldPanic    bool
		backup         *model.BackupResponse
		expected       *backups.BackupChildResponse
	}{
		{
			name:           "Multiple nodes",
			childUrls:      server.URL + "," + server.URL,
			childPasswords: "pass1,pass2",
			shouldPanic:    false,
			backup: &model.BackupResponse{
				Content: []byte("test"),
			},
			expected: &backups.BackupChildResponse{
				Processed: []string{server.URL, server.URL},
				Took:      1.0, // 0.5 + 0.5 from two nodes
			},
		},
		{
			name:           "Empty values",
			childUrls:      "",
			childPasswords: "",
			shouldPanic:    true,
			backup: &model.BackupResponse{
				Content: []byte("test"),
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("CHILD_URLS", tt.childUrls)
			os.Setenv("CHILD_PASSWORDS", tt.childPasswords)
			defer os.Unsetenv("CHILD_URLS")
			defer os.Unsetenv("CHILD_PASSWORDS")

			// Create a test client
			client := server.Client()

			if tt.shouldPanic {
				assert.Panics(t, func() {
					backups.RestoreBackupInChilds(client, tt.backup)
				})
				return
			}

			// Call the function
			result := backups.RestoreBackupInChilds(client, tt.backup)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}
