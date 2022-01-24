package clone

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

type cloneData struct {
	repos    []string
	cloneDir string
}

func Clone(ctx context.Context, args []string) error {
	c := gitHubAuth(ctx)
	err := cloneUserRepos(ctx, c)
	return err
}
func gitHubAuth(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken()},
	)
	tc := oauth2.NewClient(ctx, ts)
	c := github.NewClient(tc)
	return c
}

func ghToken() string {
	if pat := os.Getenv("GHPAT"); strings.EqualFold(pat, "") {
		log.Fatal("github pat not set")
		os.Exit(1)
	}
	return os.Getenv("GHPAT")
}

func cloneUserRepos(ctx context.Context, c *github.Client) error {
	repo := queryRepo(ctx, c)
	repoList := make([]string, len(repo)-1)

	for i := range repo {
		r := repo[i]
		log.Print(*r.GitURL)
		go func(n string) {
			go cloneRepo(ctx, n)
		}(*r.GitURL)
		fmt.Println()
		repoList = append(repoList, *r.GitURL)
	}
	// for r := range repoList {
	// 	go func(n string) {
	// 		go cloneRepo(ctx, n)
	// 	}(repoList[r])
	// }
	return nil
}

func queryRepo(ctx context.Context, c *github.Client) []*github.Repository {
	resp, _, err := c.Repositories.List(ctx, "", nil)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func cloneRepo(ctx context.Context, url string) {
	log.Print("cloning url:", url)
	// exec.Command("git", "clone", *url)
}
