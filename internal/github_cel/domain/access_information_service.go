package domain

import (
	"context"
	"fmt"
)

var NotValidOrganizations = fmt.Errorf("organizations not valid")

type AccessControlRepository interface {
	Scan(ctx context.Context, repository Repository) *RepositoryAccessInformation
}
type GitRepository interface {
	GetOrganizationUsersId(ctx context.Context, organization string) ([]UserID, error)
	GetRepositoryUsers(ctx context.Context, repository string) ([]RepositoryUser, error)
}

type AccessInformationService struct {
	newGitRepository           func(string) GitRepository
	newAccessControlRepository func(string) (AccessControlRepository, error)
}

func NewAccessInformationService(newGitRepository func(string) GitRepository, newAccessControlRepository func(string) (AccessControlRepository, error)) *AccessInformationService {
	return &AccessInformationService{
		newGitRepository:           newGitRepository,
		newAccessControlRepository: newAccessControlRepository,
	}
}

func (c *AccessInformationService) Run(ctx context.Context, gitType GitType, accessControlType AccessControlType, repositoriesUrl []string, organizations string, scanText string, gitToken string) (*AccessInformationResponse, error) {
	//di by git repo key
	gitRepository := c.newGitRepository(gitToken)
	accessControlRepository, err := c.newAccessControlRepository(scanText)

	if err != nil {
		return nil, err
	}

	organizationUsersMap, err := organizationUsersMap(ctx, gitRepository, organizations)
	if err != nil {
		return nil, NotValidOrganizations
	}

	repositories := repositoryUsers(ctx, gitRepository, repositoriesUrl, organizationUsersMap)
	//build response (move to func)
	var accessInformationResponse AccessInformationResponse
	accessInformationResponse.ScanText = scanText
	accessInformationResponse.Organizations = organizations
	repositoriesAccessInformation := make([]RepositoryAccessInformation, 0, len(repositories))

	for _, repository := range repositories {
		repositoryAccessInformation := RepositoryAccessInformation{
			RepositoryUrl: repository.URL,
		}
		if repository.Error != nil {
			repositoryAccessInformation.Error = repository.Error
			repositoriesAccessInformation = append(repositoriesAccessInformation, repositoryAccessInformation)
			continue
		}
		accessControlRepository := accessControlRepository.Scan(ctx, repository)
		repositoriesAccessInformation = append(repositoriesAccessInformation, *accessControlRepository)
	}
	fmt.Printf("accessInformationResponse %+v\n", accessInformationResponse)
	accessInformationResponse.Repositories = repositoriesAccessInformation
	return &accessInformationResponse, nil
}
func organizationUsersMap(ctx context.Context, gitRepository GitRepository, organizations string) (map[UserID]interface{}, error) {
	oUsersID, err := gitRepository.GetOrganizationUsersId(ctx, organizations)
	if err != nil {
		return nil, NotValidOrganizations
	}
	organizationUsersMap := make(map[UserID]interface{}, len(oUsersID))
	for _, oUserID := range oUsersID {
		organizationUsersMap[oUserID] = nil
	}
	return organizationUsersMap, nil
}
func filterOrganizationUsers(users []RepositoryUser, organizationUsersIDMap map[UserID]interface{}) []RepositoryUser {
	usersRepo := make([]RepositoryUser, 0, len(users))
	for _, user := range users {
		_, ok := organizationUsersIDMap[user.ID]
		if !ok {
			continue
		}
		usersRepo = append(usersRepo, user)
	}
	return usersRepo
}
func repositoryUsers(ctx context.Context, gitRepository GitRepository, repositoriesUrl []string, organizationUsersMap map[UserID]interface{}) []Repository {
	repositories := make([]Repository, 0, len(repositoriesUrl))
	for _, repositoryURL := range repositoriesUrl {
		var repository Repository
		repositoryUsers, err := gitRepository.GetRepositoryUsers(ctx, repositoryURL)
		repository.URL = repositoryURL
		if err != nil {
			repository.Error = err
			repositories = append(repositories, repository)
			continue
		}
		usersRepo := filterOrganizationUsers(repositoryUsers, organizationUsersMap)
		repository.Users = usersRepo
		repositories = append(repositories, repository)
	}
	return repositories
}
