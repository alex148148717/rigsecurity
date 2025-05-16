package interfaces

import (
	"context"
	"rigsecurity/internal/github_cel/domain"
)

type GitAccessInformationServerImpl struct {
	accessInformationService *domain.AccessInformationService
}

func NewGitAccessInformationServer(accessInformationService *domain.AccessInformationService) GitAccessInformationV1Server {
	return &GitAccessInformationServerImpl{accessInformationService: accessInformationService}
}

func (s *GitAccessInformationServerImpl) GetAccessInformation(ctx context.Context, request *AccessInformationRequest) (*AccessInformationResponse, error) {
	gitType := domain.GitType(request.GitType)
	accessControlType := domain.AccessControlType(request.AccessControlType)
	repositoriesUrl := request.RepositoriesUrl
	organizations := request.Organizations
	scanText := request.ScanText
	gitToken := request.GitToken
	accessInformationResponse, err := s.accessInformationService.Run(ctx, gitType, accessControlType, repositoriesUrl, organizations, scanText, gitToken)

	if err != nil {
		return nil, err
	}
	repositoriesAccessInformation := make([]*RepositoryAccessInformation, 0, len(accessInformationResponse.Repositories))
	for _, repository := range accessInformationResponse.Repositories {
		usersInformation := make([]*UserInformation, 0, len(repository.UserInformation))
		repositories := RepositoryAccessInformation{}
		for _, ui := range repository.UserInformation {
			var userInformation UserInformation
			ru := ui.RepositoryUser
			permissions := make([]*Permissions, 0, len(ru.Permissions))
			for _, p := range ru.Permissions {
				permissions = append(permissions, &Permissions{
					Key: p.Key, Value: p.Value,
				})
			}

			repositoryUser := RepositoryUser{
				Id:          &UserID{Value: string(ru.ID)},
				UserName:    ru.UserName,
				Email:       ru.Email,
				Type:        UserType(ru.Type),
				Permissions: permissions,
			}
			userInformation.RepositoryUser = &repositoryUser
			userInformation.Ok = ui.OK
			userInformation.Message = ui.Message
			userInformation.Error = ui.Message
			usersInformation = append(usersInformation, &userInformation)
		}
		repositories.RepositoryUrl = repository.RepositoryUrl
		if repository.Error != nil {
			repositories.Error = repository.Error.Error()
		}
		repositories.UserInformation = usersInformation
		repositoriesAccessInformation = append(repositoriesAccessInformation, &repositories)
	}

	return &AccessInformationResponse{
		Organizations: accessInformationResponse.Organizations,
		ScanText:      accessInformationResponse.ScanText,
		Repositories:  repositoriesAccessInformation,
	}, nil
}

func (s *GitAccessInformationServerImpl) mustEmbedUnimplementedGitAccessInformationV1Server() {
	return
}
