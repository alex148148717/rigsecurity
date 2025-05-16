package infrastructure

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"testing"
)

func TestGetOrganizationUsers(t *testing.T) {
	token := "token"
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	gitHubRepository := GitHubRepository{githubClient: client}
	ctx := context.Background()
	users, err := gitHubRepository.GetOrganizationUsersId(ctx, "https://github.com/testalexgreenman")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(users)
}

func TestGetRepositoryUsers(t *testing.T) {
	token := "token"
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	gitHubRepository := GitHubRepository{githubClient: client}
	ctx := context.Background()
	users, err := gitHubRepository.GetRepositoryUsers(ctx, "https://github.com/testalexgreenman/testme")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(users)
}
