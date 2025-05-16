package domain

type GitType int

const Github GitType = 0

type AccessControlType int

const CEL AccessControlType = 0

type UserType int

const (
	Guest UserType = 0
	User  UserType = 1
	Admin UserType = 2
)

type UserID string

type Repository struct {
	URL   string
	Error error
	Users []RepositoryUser
}

type Permissions struct {
	Key   string
	Value bool
}
type RepositoryUser struct {
	ID          UserID
	UserName    string
	Email       string
	Type        UserType
	Permissions []Permissions
}
type UserInformation struct {
	OK             bool
	Message        string
	Error          error
	RepositoryUser RepositoryUser
}
type RepositoryAccessInformation struct {
	RepositoryUrl   string
	Error           error
	UserInformation []UserInformation
}

type AccessInformationResponse struct {
	Repositories  []RepositoryAccessInformation
	Organizations string
	ScanText      string
}
