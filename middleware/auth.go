package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

var (
	authToken   string
	tokenExpiry time.Time
	tokenMutex  sync.Mutex
	// apiKey      = os.Getenv("SANDBOX_API_KEY")
	// apiSecret   = os.Getenv("SANDBOX_API_SECRET")
	sandboxURL = "https://api.sandbox.co.in/authenticate"
	refreshURL = "https://api.sandbox.co.in/authorize"
	apiVersion = "2.0"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenMutex.Lock()
		defer tokenMutex.Unlock()

		if authToken == "" || time.Now().After(tokenExpiry) {
			if err := FetchAuthToken(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate"})
				c.Abort()
				return
			}
		}

		c.Request.Header.Set("Authorization", "Bearer "+authToken)
		c.Next()
	}
}

func FetchAuthToken() error {
	if authToken != "" {
		if err := RefreshAuthToken(); err == nil {
			return nil
		}
	}

	req, err := http.NewRequest("POST", sandboxURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-api-key", os.Getenv("SANDBOX_API_KEY"))
	req.Header.Add("x-api-secret", os.Getenv("SANDBOX_API_SECRET"))
	req.Header.Add("x-api-version", apiVersion)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate: %s", string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return err
	}

	authToken = tokenResp.AccessToken
	tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return nil
}

func RefreshAuthToken() error {
	req, err := http.NewRequest("POST", refreshURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("request_token", authToken)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("x-api-key", os.Getenv("SANDBOX_API_KEY"))
	req.Header.Add("x-api-version", apiVersion)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to refresh token: %s", string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return err
	}

	authToken = tokenResp.AccessToken
	tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return nil
}
