package registration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRegistraterBroadcaster(t *testing.T) {

	// mock registration api
	go func() {
		http.HandleFunc("/broadcaster/register", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"broadcaster_id": "test-broadcaster-id"}`)
		})
		http.HandleFunc("/broadcaster/unregister", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"err": ""}`)
		})

		http.ListenAndServe(":8888", nil)
	}()

	// Set Register URL
	registration_url_env := "ocalhost:8888"
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
