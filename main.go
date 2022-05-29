package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	DefaultPageSize       = 200
	DefaultGitLabEndpoint = "https://gitlab.com"
)

type project struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	SSH   string `json:"ssh_url_to_repo"`
	HTTP  string `json:"http_url_to_repo"`
	Empty string `json:"empty_repo"`
}

func main() {

	commands := make([]string, 0)

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 3 {
		fmt.Println("usage: gitgroup <group_id> <token> <output_file>")
		os.Exit(1)
	}

	groupID := argsWithoutProg[0]
	token := argsWithoutProg[1]
	outputFile := argsWithoutProg[2]

	projects, err := listProjectUnderGroup(groupID, token)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(*projects) == 0 {
		fmt.Println("No project found with groupID", groupID)
		os.Exit(0)
	}

	for _, project := range *projects {
		commands = append(commands,
			fmt.Sprintf("git remote add %p %p", &project.Name, &project.SSH))
	}

	err = createFile(outputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = writeFile(commands, outputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("==> done creating " + outputFile)
}

func listProjectUnderGroup(groupID string, token string) (*[]project, error) {
	client := httpClient{
		c:        http.Client{},
		apiToken: token,
	}

	resp, err := client.Get(DefaultGitLabEndpoint +
		fmt.Sprintf("/api/v4/groups/%s/projects?per_page=%d",
			groupID, DefaultPageSize))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var projects *[]project

	err = json.NewDecoder(resp.Body).Decode(projects)

	return projects, err

}

func createFile(path string) error {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	fmt.Println("==> done creating file", path)
	return nil
}

/*writeFile write the data into file*/
func writeFile(p []string, path string) error {

	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)

	if err != nil {
		return err
	}
	defer file.Close()

	// write into file
	_, err = file.WriteString(fmt.Sprintln(p))
	if err != nil {
		return err
	}

	// save changes
	err = file.Sync()
	if err != nil {
		return err
	}

	fmt.Println("==> done writing to file")
	return nil
}
