package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type nonceResponse struct {
	Nonce string `json:"nonce"`
}

type registerData struct {
	Nonce       string `json:"nonce"`
	Username    string `json:"username"`
	DisplayName string `json:"displayname,omitempty"`
	Password    string `json:"password"`
	Admin       bool   `json:"admin"`
	MAC         string `json:"mac"`
}

type registerResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	Homeserver  string `json:"home_server"`
	DeviceID    string `json:"device_id"`
}

func matrixRegister(username, password string) (string, string, error) {
	// POST /_synapse/admin/v1/register

	nonce, err := getNonce()
	if err != nil {
		return "", "", err
	}

	data := registerData{
		Nonce:    nonce,
		Username: username,
		Password: password,
		MAC:      calcHmac(nonce, username, password),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Post(cfg.MatrixAPI+"/_synapse/admin/v1/register", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			// lmao we're fucked here. double error
			panic(err)
		}
		log.Printf("status not 200 ok (%v): %s\n", resp.StatusCode, string(respBody))
		return "", "", errors.New("not 200 ok")
	}

	var accountInfo registerResponse
	err = json.NewDecoder(resp.Body).Decode(&accountInfo)
	if err != nil {
		return "", "", err
	}

	return accountInfo.UserID, accountInfo.AccessToken, nil
}

func getNonce() (string, error) {
	// GET /_synapse/admin/v1/register
	resp, err := http.Get(cfg.MatrixAPI + "/_synapse/admin/v1/register")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var nonce nonceResponse
	err = json.NewDecoder(resp.Body).Decode(&nonce)
	if err != nil {
		return "", err
	}

	return nonce.Nonce, nil
}

func calcHmac(nonce, user, password string) string {
	// The MAC is the hex digest output of the HMAC-SHA1 algorithm, with the key being the shared secret and the content
	//being the nonce, user, password, either the string "admin" or "notadmin", and optionally the user_type each separated by NULs.

	// TODO: maybe allow creation of admin users?
	data := fmt.Sprintf("%s\x00%s\x00%s\x00notadmin", nonce, user, password)
	mac := hmac.New(sha1.New, []byte(cfg.SharedSecret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
