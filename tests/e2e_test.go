package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080/api"

var senderToken, receiverToken string
var senderName, receiverName = "sender_user", "receiver_user"

func authenticate(t *testing.T, username, password string) string {
	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	require.NoError(t, err)

	resp, err := http.Post(baseURL+"/auth", "application/json", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var authResp struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	require.NoError(t, err)
	require.NotEmpty(t, authResp.Token)

	return authResp.Token
}

func sendRequest(t *testing.T, method, url string, body any, token string) *http.Response {
	var reqBody *bytes.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewReader(data)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	return resp
}

func TestE2E_BuyItem(t *testing.T) {
	senderToken = authenticate(t, senderName, "qwerty123")
	itemName := "cup"
	resp := sendRequest(t, "GET", baseURL+"/buy/"+itemName, nil, senderToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestE2E_SendCoin(t *testing.T) {
	senderToken = authenticate(t, senderName, "qwerty123")
	receiverToken = authenticate(t, receiverName, "password")

	reqBody := map[string]interface{}{
		"toUser": receiverName,
		"amount": 10,
	}

	resp := sendRequest(t, "POST", baseURL+"/sendCoin", reqBody, senderToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
