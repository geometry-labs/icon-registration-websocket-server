package registration

import (
	"testing"
)

func TestRegistraterBroadcaster(t *testing.T) {

	// Set Register URL
	registration_url_env := "35.160.158.191:8008" // os.Getenv("ICON_REGISTRATION_WEBSOCKET_REGISTRATION_URL"
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
