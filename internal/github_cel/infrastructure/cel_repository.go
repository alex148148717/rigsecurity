package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/google/cel-go/cel"
	"rigsecurity/internal/github_cel/domain"
)

type CELRepository struct {
	cellRepoEnv *cel.Env
	prg         cel.Program
}

const (
	userID          = "userID"
	userName        = "userName"
	userEmail       = "userEmail"
	userPermissions = "userPermissions"
)

func NewCelBuilder() func(scanS string) (domain.AccessControlRepository, error) {
	return NewCELRepository
}
func NewCELRepository(scanS string) (domain.AccessControlRepository, error) {
	cellRepoEnv, err := cel.NewEnv(
		cel.Variable(userID, cel.StringType),
		cel.Variable(userName, cel.StringType),
		cel.Variable(userEmail, cel.StringType),
		cel.Variable(userPermissions, cel.MapType(cel.StringType, cel.BoolType)),
	)
	if err != nil {
		return nil, err
	}
	expr, issues := cellRepoEnv.Compile(scanS)
	if issues != nil && issues.Err() != nil {
		return nil, issues.Err()
	}
	prg, err := cellRepoEnv.Program(expr)
	if err != nil {
		return nil, err
	}

	c := CELRepository{cellRepoEnv: cellRepoEnv, prg: prg}
	return &c, nil
}
func (c *CELRepository) Init(scanS string) error {

	return nil
}

func (c *CELRepository) Scan(ctx context.Context, repository domain.Repository) *domain.RepositoryAccessInformation {

	var repositoryAccessInformation domain.RepositoryAccessInformation
	usersInformation := make([]domain.UserInformation, 0, len(repository.Users))
	repositoryAccessInformation.RepositoryUrl = repository.URL
	prg := c.prg
	for _, user := range repository.Users {
		var userInformation domain.UserInformation
		userInformation.RepositoryUser = user

		permissionsMap := make(map[string]bool, len(user.Permissions))
		for _, permission := range user.Permissions {
			permissionsMap[permission.Key] = permission.Value
		}
		out, _, err := prg.Eval(map[string]interface{}{
			userID:          user.ID,
			userName:        user.UserName,
			userEmail:       user.Email,
			userPermissions: permissionsMap,
		})
		if err != nil {
			userInformation.Error = err
			usersInformation = append(usersInformation, userInformation)
			continue
		}
		if status, ok := out.Value().(bool); ok {
			userInformation.OK = status
		} else {
			b, _ := json.Marshal(out)
			userInformation.Message = string(b)
		}
		usersInformation = append(usersInformation, userInformation)
	}
	repositoryAccessInformation.UserInformation = usersInformation
	return &repositoryAccessInformation
}
