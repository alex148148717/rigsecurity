package infrastructure

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"rigsecurity/internal/github_cel/domain"
	"strings"
)

type GitHubRepository struct {
	githubClient *github.Client
}

func NetGitBuilder() func(token string) domain.GitRepository {
	return NewGitHubRepository
}

func NewGitHubRepository(token string) domain.GitRepository {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	githubClient := github.NewClient(tc)
	r := GitHubRepository{githubClient: githubClient}
	return &r
}

func GetOwner(organizations string) string {
	owner := strings.TrimPrefix(organizations, "https://github.com/")
	owner = strings.TrimPrefix(owner, "git@github.com:")
	owner = strings.TrimPrefix(owner, "github.com/")
	return owner
}

func (c *GitHubRepository) GetOrganizationUsersId(ctx context.Context, organization string) ([]domain.UserID, error) {
	owner := GetOwner(organization)
	gitUsersID, _, err := c.githubClient.Organizations.ListMembers(ctx, owner, nil)
	if err != nil {
		return nil, err
	}
	usersID := make([]domain.UserID, 0, len(gitUsersID))
	for _, gitUserID := range gitUsersID {

		if gitUserID.ID == nil {
			continue
		}
		uid := domain.UserID(fmt.Sprintf("%d", *gitUserID.ID))
		usersID = append(usersID, uid)

	}
	return usersID, err
}
func getOwnerRepository(repository string) (owner string, repositoryID string) {
	repository = strings.TrimPrefix(repository, "https://github.com/")
	repository = strings.TrimPrefix(repository, "git@github.com:")
	repository = strings.TrimPrefix(repository, "github.com/")
	parts := strings.Split(repository, "/")
	if len(parts) == 2 {
		owner = parts[0]
		repositoryID = parts[1]
	}
	return
}
func (c *GitHubRepository) GetRepositoryUsers(ctx context.Context, repository string) ([]domain.RepositoryUser, error) {
	owner, repo := getOwnerRepository(repository)
	gitUsers, _, err := c.githubClient.Repositories.ListCollaborators(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}
	users := make([]domain.RepositoryUser, 0, len(gitUsers))
	for _, gitUser := range gitUsers {
		var user domain.RepositoryUser
		if gitUser.Login != nil {
			user.UserName = *gitUser.Login
		}
		if gitUser.Type != nil {
			var userType domain.UserType
			switch *gitUser.Type {
			case "User":
				userType = domain.User
			}
			user.Type = userType
		}
		if gitUser.ID != nil {
			user.ID = domain.UserID(fmt.Sprintf("%d", *gitUser.ID))
		}
		if gitUser.Email != nil {
			user.Email = *gitUser.Email
		}
		if gitUser.Permissions != nil {
			var permissions []domain.Permissions
			for key, value := range *gitUser.Permissions {
				permissions = append(permissions, domain.Permissions{Key: key, Value: value})
			}
			user.Permissions = permissions
		}
		users = append(users, user)
	}
	return users, err

}
