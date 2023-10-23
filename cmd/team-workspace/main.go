package main

import (
	"context"
	"github.com/google/go-github/v56/github"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/pkg/errors"
	"html/template"
	"os"
)

var gitHubClient *github.Client
var apiClient *ionoscloud.APIClient

var repoAdmins = []string{
	"sebX2",
	"mspoeri",
	"ionos-landgraf-vin",
	"lpape-ionos",
}
var repositoryPrefix = "si-workshop-"
var teamName = "testteam"
var bootstrapTeamplate = "templates/argo-team-template.yaml"

type cloudUser struct {
	FirstName string
	LastName  string
}

var cloudUsers = []cloudUser{
	{"Markus", "Spoeri"},
	{"Lucas", "Pape"},
}

func main() {
	// ENV
	// admin:org, repo and SSO for ionos-cloud
	gitHubToken := os.Getenv("GITHUB_TOKEN")
	ionosCloudUser := os.Getenv("IONOS_CLOUD_USER")
	ionosCloudPassword := os.Getenv("IONOS_CLOUD_PASSWORD")
	workshopPassword := os.Getenv("WORKSHOP_PASSWORD")

	// Bootstrap
	tmpl, err := template.ParseFiles(bootstrapTeamplate)
	if err != nil {
		panic(errors.Wrap(err, "parsing template file"))
	}
	values := make(map[string]any)
	createBootstrap(teamName, tmpl, values)

	// GitHub
	gitHubClient = github.NewClient(nil).WithAuthToken(gitHubToken)
	repositoryName := repositoryPrefix + teamName
	createRepository(repositoryName)
	addRepoAdmins(repositoryName, repoAdmins)

	// Ionos Cloud
	cfg := ionoscloud.NewConfiguration(ionosCloudUser, ionosCloudPassword, "", "")
	cfg.Debug = false
	apiClient = ionoscloud.NewAPIClient(cfg)
	createUsers(workshopPassword)
}

func createRepository(name string) {
	repo := &github.Repository{
		Name:    github.String(name),
		Private: github.Bool(true),
	}
	ctx := context.Background()
	repository, _, err := gitHubClient.Repositories.Create(ctx, "ionos-cloud", repo)
	if err != nil {
		panic(errors.Wrap(err, "creating repository"))
	}
	opts := &github.RepositoryAddCollaboratorOptions{
		Permission: "admin",
	}
	_, _, err = gitHubClient.Repositories.AddCollaborator(ctx, "ionos-cloud", *repository.Name, "ionos-cloud", opts)
	if err != nil {
		panic(errors.Wrap(err, "adding ionos-cloud as collaborator"))
	}
}

func addRepoAdmins(repositoryName string, users []string) {
	ctx := context.Background()
	opts := &github.RepositoryAddCollaboratorOptions{
		Permission: "admin",
	}
	for _, user := range users {
		_, _, err := gitHubClient.Repositories.AddCollaborator(ctx, "ionos-cloud", repositoryName, user, opts)
		if err != nil {
			panic(errors.Wrap(err, "adding collaborator"))
		}
	}

}

func createBootstrap(teamName string, tmpl *template.Template, values map[string]any) {

	values["TeamName"] = teamName
	values["RepoPrefix"] = repositoryPrefix

	output, err := os.Create("bootstrap/" + teamName + ".yaml")
	if err != nil {
		panic(errors.Wrap(err, "creating output file"))
	}
	defer output.Close()
	err = tmpl.Execute(output, values)
	if err != nil {
		panic(errors.Wrap(err, "executing template file"))
	}
}

func createUsers(workshopPassword string) {
	for _, user := range cloudUsers {
		createUser(user.FirstName, user.LastName, workshopPassword)
	}
}

func createUser(firstname, lastname, workshopPassword string) {

	admin := true
	secAuth := false
	active := true
	mailAddress := (firstname + "." + lastname + "+workshop" + "@workspace.ionos.com")

	properties := ionoscloud.UserPropertiesPost{
		Firstname:     &firstname,
		Lastname:      &lastname,
		Email:         &mailAddress,
		Administrator: &admin,
		Active:        &active,
		ForceSecAuth:  &secAuth,
		Password:      &workshopPassword,
	}
	user := ionoscloud.UserPost{
		Properties: &properties,
	}
	_, _, err := apiClient.UserManagementApi.UmUsersPost(context.Background()).User(user).Depth(0).Execute()
	if err != nil {
		panic(errors.Wrap(err, "creating user"))
	}
}
