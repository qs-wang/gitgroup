package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type project struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	SSH   string `json:"ssh_url_to_repo"`
	HTTP  string `json:"http_url_to_repo"`
	Empty bool   `json:"empty_repo"`
}

type subgroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func listProjectUnderGroup(groupID int, token string) (*[]project, error) {
	client := httpClient{
		c:        http.Client{},
		apiToken: token,
	}

	resp, err := client.Get(DefaultGitLabEndpoint +
		fmt.Sprintf("/api/v4/groups/%d/projects?per_page=%d",
			groupID, DefaultPageSize))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("request failed with code %d", resp.StatusCode)
	}

	var projects []project

	// body, err := io.ReadAll(resp.Body)

	// fmt.Println(body)

	err = json.NewDecoder(resp.Body).Decode(&projects)

	return &projects, err

}

func listSubGroupsUnderGroup(groupID string, token string) (*[]subgroup, error) {
	client := httpClient{
		c:        http.Client{},
		apiToken: token,
	}

	resp, err := client.Get(DefaultGitLabEndpoint +
		fmt.Sprintf("/api/v4/groups/%s/subgroups?per_page=%d&all_available=true",
			groupID, DefaultPageSize))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("request failed with code %d", resp.StatusCode)
	}

	var subgroups []subgroup

	// body, err := io.ReadAll(resp.Body)

	// fmt.Println(body)

	err = json.NewDecoder(resp.Body).Decode(&subgroups)

	return &subgroups, err

}
