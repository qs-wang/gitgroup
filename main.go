package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultPageSize        = 200
	DefaultGitLabEndpoint  = "https://gitlab.com"
	DefaultMonoBuildScript = " ~/git/monorepo-tools/monorepo_build.sh"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 4 {
		fmt.Println("usage: gitgroup <group_id> <token> <output_file> main-rep")
		os.Exit(1)
	}

	groupID := argsWithoutProg[0]
	token := argsWithoutProg[1]
	outputFile := argsWithoutProg[2]
	mainRepo := argsWithoutProg[3]

	group, _ := strconv.Atoi(groupID)

	projects, err := listProjectUnderGroup(group, token)

	if err != nil {
		fmt.Println("Error while geeting  projects from group", groupID)
		os.Exit(0)
	}

	subGroups, err := listSubGroupsUnderGroup(groupID, token)
	if err != nil {
		fmt.Println("Error while geeting  subgroups from group", groupID)
		os.Exit(0)
	}

	for _, group := range *subGroups {
		subProjects, _ := listProjectUnderGroup(group.ID, token)
		*projects = append(*projects, *subProjects...)
	}

	fmt.Println("Found project total:", len(*projects))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(*projects) == 0 {
		fmt.Println("No project found with groupID", groupID)
		os.Exit(0)
	}

	commands := make([]string, 0)
	removeCommands := make([]string, 0)

	repoNames := make([]string, 0)
	repoNames = append(repoNames, mainRepo)
	for _, project := range *projects {
		if project.Empty {
			fmt.Println("Empty repo", project.Name)
			continue
		}
		name := strings.ReplaceAll(project.Name, " ", "-")
		commands = append(commands,
			fmt.Sprintf("git remote add %s %s", name, project.SSH))
		removeCommands = append(removeCommands, fmt.Sprintf("git remote remove %s", name))
		repoNames = append(repoNames, name)
	}

	commands = append(commands, fmt.Sprintln("git fetch --all --no-tags"))
	monoCommd := DefaultMonoBuildScript

	for _, cmd := range repoNames {
		monoCommd += fmt.Sprintf(" %s ", cmd)
	}

	monoCommd += fmt.Sprintln("")

	commands = append(commands, monoCommd)

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

	err = createFile(outputFile + "-remove")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = writeFile(removeCommands, outputFile+"-remove")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("==> done creating " + outputFile)
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
func writeFile(lines []string, path string) error {

	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)

	if err != nil {
		return err
	}
	defer file.Close()

	// write into file
	for _, l := range lines {
		_, err = file.WriteString(fmt.Sprintln(l))
		if err != nil {
			return err
		}
	}

	// save changes
	err = file.Sync()
	if err != nil {
		return err
	}

	fmt.Println("==> done writing to file")
	return nil
}
