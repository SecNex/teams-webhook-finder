package main

import (
	"log"

	"github.com/secnex/ms-toolbox/webhook-finder/api"
	"github.com/secnex/ms-toolbox/webhook-finder/api/teams"
	"github.com/secnex/ms-toolbox/webhook-finder/config"
)

func main() {
	c := config.NewCustomConfig("config.json")
	file, _ := c.Load()
	log.Printf("Client ID: %s\n", file.Client.ClientID)
	log.Printf("Client Secret: %s\n", file.Client.ClientSecret)
	log.Printf("Tenant ID: %s\n", file.Client.TenantID)

	graph := api.NewMsGraph(file.Client.ClientID, file.Client.ClientSecret, file.Client.TenantID)
	graph.NewAPIClient()
	graph.Authenticate(api.SCOPE_MS_GRAPH_DEFAULT)
	teamsAPI := teams.NewTeamsAPIBuilder(api.MS_GRAPH_API)
	listApps := teamsAPI.ListAppCatalog(true)
	resp, _ := graph.GET(listApps)
	listAppsResponse, _ := teamsAPI.GetListAppCatalogResponse(resp)
	appCatalog := listAppsResponse.Value
	var incomingWebhookAppId string
	for _, app := range appCatalog {
		if app.DisplayName == "Incoming Webhook" {
			incomingWebhookAppId = app.Id
		}
	}
	listTeamsUrl := teamsAPI.ListTeams()
	resp, _ = graph.GET(listTeamsUrl)
	listTeamsResponse, _ := teamsAPI.GetListTeamsResponse(resp)
	teams := listTeamsResponse.Value
	for _, team := range teams {
		listAppsTeamsUrl := teamsAPI.ListAppsInTeam(team.Id, true)
		resp, _ := graph.GET(listAppsTeamsUrl)
		listAppsInTeamResponse, _ := teamsAPI.GetListAppsInTeamResponse(resp)
		if listAppsInTeamResponse == nil {
			continue
		}
		apps := listAppsInTeamResponse.Value
		for _, app := range apps {
			if app.AppDefinition.TeamsAppId == incomingWebhookAppId {
				log.Printf("Incoming Webhook App found in Team: %s\n", team.DisplayName)
			}
		}
	}
}
