package infrastructure

import (
	"context"
	"testing"
)

func TestCheck(t *testing.T) {
	c := `userName == "alex" && "admin" in userPermissions && userPermissions["admin"] == true`

	celRepository, err := NewCELRepository(c)
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()

	_ = celRepository
	_ = ctx
}
