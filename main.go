package main

import (
	"log"

	"github.com/secnex/ms-toolbox/webhook-finder/api"
	"github.com/secnex/ms-toolbox/webhook-finder/api/teams"
	"github.com/secnex/ms-toolbox/webhook-finder/config"
)

// Fix the pagination issue for listing Teams - v.0.1.2-beta
func teamNextLink(graph *api.MsGraph, api *teams.TeamsAPIBuilder, url string) *teams.ListTeamsResponse {
	log.Printf("Next Link: %s\n", url)
	resp, _ := graph.GET(url)
	listTeamsResponse, _ := api.GetListTeamsResponse(resp)
	return listTeamsResponse
}

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
	if listAppsResponse == nil {
		log.Println("No Apps found!")
		return
	}
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
	if listTeamsResponse == nil {
		log.Println("No Teams found!")
		return
	}
	teams := listTeamsResponse.Value
	teamsNextLink := false
	for teamsNextLink || listTeamsResponse.ODataNextLink != "" {
		if listTeamsResponse.ODataNextLink != "" {
			listTeamsResponse = teamNextLink(graph, teamsAPI, listTeamsResponse.ODataNextLink)
			teams = append(teams, listTeamsResponse.Value...)
			teamsNextLink = true
		} else {
			teamsNextLink = false
		}
	}
	log.Printf("Teams found: %d\n", len(teams))
	for _, team := range teams {
		log.Printf("Team: %s\n", team.DisplayName)
	}
	affectedTeams := 0
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
				affectedTeams++
			}
		}
	}
	log.Printf("Teams affected: %d/%d\n", affectedTeams, len(teams))
}
