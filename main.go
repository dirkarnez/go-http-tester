package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
)

type LoginForm struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func main() {
	// Create a new HTTP client with a cookie jar to store the session cookie
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass the certificate verification
			},
		},
	}

	loginURL := "https://localhost:4443/api/login"
	loginRequest := LoginForm{
		UserName: "UserName",
		Password: "Password",
	}

	jsonData, _ := json.Marshal(loginRequest)

	resp, err := client.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error during login1:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error during login2:", err)
		return
	}

	// Print the token received from the login response
	fmt.Printf("Login token: %+v\n", string(body))

	// Now you can use the same HTTP client to make authenticated requests
	protectedURL := "https://localhost:4443/api/users"
	req, _ := http.NewRequest("GET", protectedURL, nil)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error during request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, _ = io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}
