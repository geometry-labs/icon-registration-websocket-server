package registration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BroadcasterID string

var registrationURL string = ""

func SetRegistrationURL(url string) {

	registrationURL = url
}

func RegisterBroadcaster(register_json []byte) (BroadcasterID, error) {

	if registrationURL == "" {
		return "", errors.New("registration url not set")
	}

	endpoint := "http://" + registrationURL + "/broadcaster/register"

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(register_json))
	if err != nil {
		return "", err
	}

	// Read response
	defer resp.Body.Close()
	resp_body_raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Cannot read reponse from registraion")
	}

	if resp.StatusCode != 200 {
		error_msg := fmt.Sprintf("Invalid registration response: %d - %s", resp.StatusCode, string(resp_body_raw))

		return "", errors.New(error_msg)
	}

	// Parse response
	resp_map := make(map[string]interface{})
	json.Unmarshal(resp_body_raw, &resp_map)

	broadcaster_id := BroadcasterID(resp_map["broadcaster_id"].(string))

	return broadcaster_id, nil
}
