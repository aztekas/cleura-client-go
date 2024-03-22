package cleura

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default Cleura API Endpoint prefix
const HostURL string = "https://rest.cleura.cloud" //"http://localhost:8088/"

// Need to know status code for retry
type RequestAPIError struct {
	StatusCode int
	Err        error
}

func (r *RequestAPIError) Error() string {
	return r.Err.Error()
}

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// AuthStruct Wrapper
type AuthStructWrapper struct {
	Auth AuthStruct `json:"auth"`
}

// AuthStruct -
type AuthStruct struct {
	Username         string `json:"login"`
	Password         string `json:"password"`
	TwoFactorMethod  string `json:"twofa_method,omitempty"`
	VerificationCode string `json:"verification,omitempty"`
}

// AuthResponse -
type AuthResponse struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

// AuthVerificationResult
type AuthVerificationResult struct {
	Result       string `json:"result"`
	Type         string `json:"type"`
	Verification string `json:"verification"`
}

type AuthRequestTwoFactor struct {
	Request2FA AuthRequestTwoFactorDetails `json:"request2fa"`
}
type AuthRequestTwoFactorDetails struct {
	Login        string `json:"login"`
	Verification string `json:"verification"`
}

type AuthVerifyTwoFactor struct {
	Verify2FA AuthVerifyTwoFactorDetails `json:"verify2fa"`
}
type AuthVerifyTwoFactorDetails struct {
	Login        string `json:"login"`
	Verification string `json:"verification"`
	Code         int    `json:"code"`
}

// NewClient -
func NewClient(host, username, password *string, twoFactorAuthEnabled bool) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 600 * time.Second},
		// Default API URL
		HostURL: HostURL,
	}
	if host != nil {
		c.HostURL = *host
	}
	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}
	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}
	// Return client without token if two factor is enabled
	if twoFactorAuthEnabled {
		c.Auth.TwoFactorMethod = "sms"
		err := c.Request2FactorCode()
		if err != nil {
			return nil, err
		}
		return &c, nil
	}
	ar, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	c.Token = ar.Token
	return &c, nil
}

func NewClientNoPassword(host, username, token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 600 * time.Second},
		// Default API URL
		HostURL: HostURL,
	}
	if host != nil {
		c.HostURL = *host
	}
	// If username or password not provided, return empty client
	if username == nil || token == nil {
		return &c, nil
	}
	c.Auth = AuthStruct{
		Username: *username,
	}
	c.Token = *token
	return &c, nil
}

func (c *Client) doRequest(req *http.Request, successResponse int) ([]byte, error) {
	token := c.Token
	req.Header.Set("X-AUTH-LOGIN", c.Auth.Username)
	req.Header.Set("X-AUTH-TOKEN", token)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != successResponse {
		rErr := RequestAPIError{
			Err:        fmt.Errorf("actual_status: %d, expected_status: %d, body: %s", res.StatusCode, successResponse, body),
			StatusCode: res.StatusCode,
		}
		return nil, &rErr
	}
	return body, nil
}
