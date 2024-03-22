package cleura

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetToken - Get a new token for the user
func (c *Client) GetToken() (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}

	auth := AuthStructWrapper{
		Auth: c.Auth,
	}
	rb, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/tokens", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, 200)
	if err != nil {
		return nil, err
	}
	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	if ar.Result!="login_ok"{
		return nil, fmt.Errorf("error: get token failed, expected result: `login_ok`, got: %s,",ar.Result)
	}
	return &ar, nil
}

// RevokeToken - Revoke client token
func (c *Client) RevokeToken() error {
	//https://rest.cleura.cloud/auth/v1/tokens
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/auth/v1/tokens", c.HostURL), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req, 204)
	if err != nil {
		return err
	}
	return nil
}

// Validate client token
func (c *Client) ValidateToken() error {
	//https://rest.cleura.cloud/auth/v1/tokens/validate

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/tokens/validate", c.HostURL), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req, 204)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Request2FactorCode() error {
	//https://rest.cleura.cloud/auth/v1/tokens
	// get verification code
	var authVerificationResult AuthVerificationResult
	auth := AuthStructWrapper{
		Auth: c.Auth,
	}
	rb, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/tokens", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	body, err := c.doRequest(req, 200)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &authVerificationResult); err != nil {
		return err
	}
	request2FaDetails := &AuthRequestTwoFactorDetails{
		Login:        c.Auth.Username,
		Verification: authVerificationResult.Verification,
	}
	request2Fa := &AuthRequestTwoFactor{
		Request2FA: *request2FaDetails,
	}
	rb, err = json.Marshal(request2Fa)
	if err != nil {
		return err
	}
	//https://rest.cleura.cloud/auth/v1/tokens/request2facode
	req, err = http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/tokens/request2facode", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	_, err = c.doRequest(req, 204)
	if err != nil {
		return err
	}
	c.Auth.VerificationCode = authVerificationResult.Verification
	//request 2factor code via sms
	return nil
}

func (c *Client) GetTokenWith2FA(twoFACodeFromSms int) error {
	//issue token request with two-factor code received via sms
	var ar AuthResponse
	verify2FaDetails := &AuthVerifyTwoFactorDetails{
		Login:        c.Auth.Username,
		Verification: c.Auth.VerificationCode,
		Code:         twoFACodeFromSms,
	}
	verify2Fa := &AuthVerifyTwoFactor{
		Verify2FA: *verify2FaDetails,
	}
	rb, err := json.Marshal(verify2Fa)
	if err != nil {
		return err
	}
	//https://rest.cleura.cloud/auth/v1/tokens/verify2fa
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/tokens/verify2fa", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	jsonByte, err := c.doRequest(req, 200)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonByte, &ar)
	if err != nil {
		return err
	}
	c.Token = ar.Token
	return nil
}
