package teams

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const TEAMS_LIST = "/teams"
const TEAMS_LIST_APPS_IN_TEAM = "/teams/%s/installedApps"
const TEAMS_LIST_APP_CATALOG = "/appCatalogs/teamsApps"

type TeamsAPIBuilder struct {
	BaseURL string
	ApiPath string
}

type ListTeamsResponse struct {
	Value []Team `json:"value"`
}

type ListAppsInTeamResponse struct {
	Value []InstalledApp `json:"value"`
}

type ListAppCatalogResponse struct {
	Value []AppCatalog `json:"value"`
}

type Team struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type InstalledApp struct {
	Id            string         `json:"id"`
	AppDefinition AppDefinitions `json:"teamsAppDefinition"`
}

type AppCatalog struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	// AppDefinitions AppDefinitions `json:"appDefinitions"`
}

type AppDefinitions struct {
	Id          string `json:"id"`
	TeamsAppId  string `json:"teamsAppId"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

func NewTeamsAPIBuilder(baseURL string) *TeamsAPIBuilder {
	return &TeamsAPIBuilder{
		BaseURL: baseURL,
	}
}

func (t *TeamsAPIBuilder) ListTeams() string {
	return t.BaseURL + "/v1.0" + TEAMS_LIST + "?$select=id,displayName,description&$top=999"
}

func (t *TeamsAPIBuilder) ListAppsInTeam(teamID string, appDefinition bool) string {
	if appDefinition {
		return t.BaseURL + "/v1.0" + fmt.Sprintf(TEAMS_LIST_APPS_IN_TEAM, teamID) + "?$expand=teamsAppDefinition"
	}
	return t.BaseURL + "/v1.0" + fmt.Sprintf(TEAMS_LIST_APPS_IN_TEAM, teamID)
}

func (t *TeamsAPIBuilder) ListAppCatalog(appDefinition bool) string {
	if appDefinition {
		return t.BaseURL + "/v1.0" + TEAMS_LIST_APP_CATALOG + "?$expand=appDefinitions"
	}
	return t.BaseURL + "/v1.0" + TEAMS_LIST_APP_CATALOG
}

func (t *TeamsAPIBuilder) GetListTeamsResponse(response *http.Response) (*ListTeamsResponse, error) {
	var listTeamsResponse ListTeamsResponse
	err := json.NewDecoder(response.Body).Decode(&listTeamsResponse)
	if err != nil {
		return nil, err
	}
	return &listTeamsResponse, nil
}

func (t *TeamsAPIBuilder) GetListAppsInTeamResponse(response *http.Response) (*ListAppsInTeamResponse, error) {
	var listAppsInTeamResponse ListAppsInTeamResponse
	err := json.NewDecoder(response.Body).Decode(&listAppsInTeamResponse)
	if err != nil {
		return nil, err
	}
	return &listAppsInTeamResponse, nil
}

func (t *TeamsAPIBuilder) GetListAppCatalogResponse(response *http.Response) (*ListAppCatalogResponse, error) {
	var listAppCatalogResponse ListAppCatalogResponse
	err := json.NewDecoder(response.Body).Decode(&listAppCatalogResponse)
	if err != nil {
		return nil, err
	}
	return &listAppCatalogResponse, nil
}
