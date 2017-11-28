package main

import (
	"context"
	"fmt"
	"log"
	"os"

	aw "github.com/deanishe/awgo"
	"github.com/google/go-github/github"

	"golang.org/x/oauth2"
)

func main() {
	aw.Run(run)
}

func run() {
	var (
		token string = os.Getenv("GITHUB_TOKEN")
		owner string = os.Getenv("GITHUB_OWNER")
		repo  string = os.Getenv("GITHUB_REPO")
	)
	wf := aw.New()

	iconAvailable := &aw.Icon{Value: os.Getenv("GITHUB_ICON_PATH")}

	args := wf.Args()
	var query string
	if len(args) > 0 {
		query = args[0]
	}

	if token == "" {
		wf.NewWarningItem("github token is missing!", "subtitle")
		wf.SendFeedback()
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	projects, _, err := client.Repositories.ListProjects(context.Background(), owner, repo, &github.ProjectListOptions{})
	if err != nil {
		wf.FatalError(err)
	}

	for _, project := range projects {
		url := fmt.Sprintf("https://github.com/%s/%s/projects/%d", owner, repo, project.GetNumber())
		wf.NewItem(project.GetName()).Subtitle(url).Arg(url).Valid(true).Icon(iconAvailable)
	}

	if query != "" {
		res := wf.Filter(query)
		log.Printf("%d results match \"%s\"", len(res), query)
		for i, r := range res {
			log.Printf("%02d. score=%0.1f sortkey=%s", i+1, r.Score, wf.Feedback.SortKey(i))
		}
	}

	wf.WarnEmpty("No matching", "Try a different query")
	wf.SendFeedback()
}
