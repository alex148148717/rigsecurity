# Git Access Information Service

This service defines a gRPC API using Protocol Buffers (`proto3`) to fetch and validate user access information for Git repositories (currently supports GitHub) using access control engines like CEL.

## Package Details

- **Proto syntax:** proto3
- **Proto package:** `git_access_information`
- **Go package path:** `internal/github_cel/interfaces`

---

## Enumerations

### `GitType`
Defines the source control platform:
- `Github`: Currently the only supported provider.

### `AccessControlType`
Defines the access control policy engine:
- `CEL`: Google's Common Expression Language.

### `UserType`
Defines the role of the user within the repository:
- `Guest`: Limited access.
- `User`: Standard access.
- `Admin`: Full administrative rights.

---

## Messages

### `Permissions`
Represents a single permission flag:
- `key`: Name of the permission (e.g. "write", "read").
- `value`: Boolean indicating whether the user has it.

### `UserID`
Wraps a string value representing a user ID.

### `RepositoryUser`
Represents a user with repository access:
- `id`: Unique identifier.
- `userName`: GitHub username.
- `email`: User email.
- `type`: One of the enum `UserType`.
- `permissions`: A list of key-value permissions.

### `UserInformation`
Provides the result of a policy evaluation for a single user:
- `ok`: Whether the user passed the access rules.
- `message`: Human-readable result message.
- `error`: Optional error if evaluation failed.
- `repositoryUser`: The user being evaluated.

### `RepositoryAccessInformation`
Access data for a single repository:
- `repositoryUrl`: The Git repo URL.
- `error`: Any error encountered during analysis.
- `userInformation`: List of evaluated users and their statuses.

### `AccessInformationRequest`
Input message sent to the API:
- `gitType`: The type of Git service (e.g. GitHub).
- `accessControlType`: The policy engine to use (e.g. CEL).
- `repositoriesUrl`: List of repository URLs to analyze.
- `organizations`: GitHub organization name.
- `scanText`: Policy expression in CEL.
- `gitToken`: Authentication token.

### `AccessInformationResponse`
The result returned from the API:
- `repositories`: A list of `RepositoryAccessInformation` entries.
- `organizations`: The organization name (echoed).
- `scanText`: The policy text (echoed).

---

## gRPC Service

### `GitAccessInformationV1`
Main service exposing the following method:

#### `GetAccessInformation`
- **Input:** `AccessInformationRequest`
- **Output:** `AccessInformationResponse`

This method processes repository access based on the given parameters and evaluates each user according to the specified policy engine.

---

## Code Generation (Go)

To generate the Go code from the `.proto` file, simply run:

```bash
docker compose up --build
```

This uses the `protoc-builder` service defined in your `docker-compose.yml` to compile the protobuf definitions into Go code automatically.

---

## Example Usage (Go Client)

The following is an example of how to call the gRPC service using Go:

```go
package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"rigsecurity/internal/github_cel/interfaces"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := interfaces.NewGitAccessInformationV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &interfaces.AccessInformationRequest{
		GitType:           interfaces.GitType_Github,
		AccessControlType: interfaces.AccessControlType_CEL,
		RepositoriesUrl:   []string{"https://github.com/testalexgreenman/testme"},
		Organizations:     "https://github.com/testalexgreenman",
		ScanText:          `userName == "alex148148717" `,
		GitToken:          "token",
	}

	resp, err := client.GetAccessInformation(ctx, req)
	if err != nil {
		log.Fatalf("could not get access information: %v", err)
	}

	log.Printf("Access: %v", resp)
}
```

---

## Explanation: AccessInformationRequest Fields (Go Client)

The following explains each field used when creating the `AccessInformationRequest` message:

```go
GitType: interfaces.GitType_Github,
```
Specifies the Git provider. Here, it is set to `Github`, indicating that the repository is hosted on GitHub.

```go
AccessControlType: interfaces.AccessControlType_CEL,
```
Specifies the access control engine. `CEL` (Common Expression Language) allows you to write policy rules in a safe expression language.

```go
RepositoriesUrl: []string{"https://github.com/testalexgreenman/testme"},
```
A list of repository URLs to scan. Each repository will be checked according to the given access rules.

```go
Organizations: "https://github.com/testalexgreenman",
```
The GitHub organization or owner URL associated with the repositories.

```go
ScanText: `userName == "alex148148717" `,
```
The access control rule written in CEL. In this example, it allows access only to users with the username "alex148148717".

```go
GitToken: "token",
```
A personal access token (PAT) used to authenticate with the GitHub API.
