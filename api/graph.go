package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const MS_GRAPH_API = "https://graph.microsoft.com"
const MS_OAUTH_API = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
const SCOPE_MS_GRAPH_DEFAULT = "https://graph.microsoft.com/.default"

type MsGraphClientCredentials struct {
	ClientID     string
	ClientSecret string
	TenantID     string
}

type MsGraphAPIClient struct {
	Client      *http.Client
	AccessToken string
}

type MsGraph struct {
	ClientCredentials MsGraphClientCredentials
	APIClient         *MsGraphAPIClient
}

type MsGraphTokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

type Scopes []Scope
type Scope string

func NewMsGraph(clientID, clientSecret, tenantID string) *MsGraph {
	return &MsGraph{
		ClientCredentials: MsGraphClientCredentials{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TenantID:     tenantID,
		},
	}
}

func NewScopes(scopes ...Scope) Scopes {
	return scopes
}

func (m *MsGraph) NewAPIClient() *MsGraphAPIClient {
	m.APIClient = &MsGraphAPIClient{
		Client: &http.Client{},
	}
	return m.APIClient
}

func (m *MsGraph) Authenticate(scope Scope) (*MsGraphTokenResponse, error) {
	uri := fmt.Sprintf(MS_OAUTH_API, m.ClientCredentials.TenantID)
	log.Printf("Authenticating with scope: %s\n", scope)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Add("client_id", m.ClientCredentials.ClientID)
	data.Add("client_secret", m.ClientCredentials.ClientSecret)
	data.Add("scope", string(scope))

	u, _ := url.ParseRequestURI(uri)
	urlStr := fmt.Sprintf("%v", u)

	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.APIClient.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed: %s", resp.Status)
	}

	var tokenResponse MsGraphTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, err
	}
	if tokenResponse.AccessToken == "" {
		return nil, fmt.Errorf("authentication failed: empty access token")
	}
	m.APIClient.AccessToken = tokenResponse.AccessToken

	return &tokenResponse, nil
}

func (m *MsGraph) GetAccessToken() string {
	return m.APIClient.AccessToken
}

func (m *MsGraph) GET(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+m.APIClient.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.APIClient.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET request failed: %s", resp.Status)
	}

	return resp, nil
}
