package registration

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestRegistraterBroadcaster(t *testing.T) {

	// mock registration api
	mock_registration_api_port := ":8888"
	go func() {
		http.HandleFunc("/broadcaster/register", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"broadcaster_id": "test-broadcaster-id"}`)
		})
		http.HandleFunc("/broadcaster/unregister", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"err": ""}`)
		})
		http.ListenAndServe(mock_registration_api_port, nil)
		t.Logf("Failed to mock registration api")
		t.Fail()
	}()

	// Wait for mock registration api
	time.Sleep(1 * time.Second)

	// Set Register URL
	registration_url_env := "localhost" + mock_registration_api_port
	SetRegistrationURL(registration_url_env)

	// Test json config
	register_json := `
	{
					"connection_type": "ws",
					"endpoint": "wss://test",
					"transaction_events": [
							{
									"to_address": "cx0000000000000000000000000000000000000000"
							}
					]
	}
	`

	// Register
	broadcaster_id, err := RegisterBroadcaster([]byte(register_json))
	if err != nil {
		t.Logf("Failed to register broadcaster")
		t.Fail()
	}

	if string(broadcaster_id) == "" {
		t.Logf("Fail: broadcaster_id returned empty")
		t.Fail()
	}

	// Unregister
	err = UnregisterBroadcaster(broadcaster_id)
	if err != nil {
		t.Logf("Failed to unregister broadcaster")
		t.Fail()
	}

	// Pass
	return
}
