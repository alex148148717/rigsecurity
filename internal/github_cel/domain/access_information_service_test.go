package domain_test

import (
	"context"
	"fmt"
	"rigsecurity/internal/github_cel/domain"
	"rigsecurity/internal/github_cel/infrastructure"
	"testing"
)

func TestRun(t *testing.T) {

	accessInformationService := domain.NewAccessInformationService(infrastructure.NewGitHubRepository, infrastructure.NewCELRepository)
	ctx := context.Background()
	repos := []string{}
	repos = []string{"https://github.com/testalexgreenman/testme"}

	organizations := "https://github.com/testalexgreenman"
	scanText := `userName == "alex148148717" `
	token := "token"

	ans, err := accessInformationService.Run(ctx, domain.Github, domain.CEL, repos, organizations, scanText, token)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ans)
}
