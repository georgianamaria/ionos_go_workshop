package main

import (
	"context"
	"github.com/google/go-github/v56/github"
	"github.com/pkg/errors"
	"html/template"
	"os"
)

var gitHubClient *github.Client
var repoAdmins = []string{
	"sebX2",
	"mspoeri",
	"ionos-landgraf-vin",
	"lpape-ionos",
}
var repositoryPrefix = "si-workshop-"
var teamName = "testteam"
var bootstrapTeamplate = "templates/argo-team-template.yaml"

func main() {
	// ENV
	// admin:org, repo and SSO for ionos-cloud
	gitHubToken := os.Getenv("GITHUB_TOKEN")

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
	gitHubClient.Repositories.AddCollaborator(ctx, "ionos-cloud", *repository.Name, "ionos-cloud", opts)
}

func addRepoAdmins(repositoryName string, users []string) {
	ctx := context.Background()
	opts := &github.RepositoryAddCollaboratorOptions{
		Permission: "admin",
	}
	for _, user := range users {
		gitHubClient.Repositories.AddCollaborator(ctx, "ionos-cloud", repositoryName, user, opts)
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
