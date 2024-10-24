package password_test

import (
	"testing"

	"github.com/ncfex/dcart/auth-service/internal/core/services/password"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewPasswordService(t *testing.T) {
	tests := []struct {
		name         string
		cost         int
		expectedCost int
	}{
		{
			name:         "zero cost uses default",
			cost:         0,
			expectedCost: bcrypt.DefaultCost,
		},
		{
			name:         "custom cost",
			cost:         12,
			expectedCost: 12,
		},
		{
			name:         "minimum cost",
			cost:         bcrypt.MinCost,
			expectedCost: bcrypt.MinCost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := password.NewPasswordService(tt.cost)
			assert.NotNil(t, service)

			hash, err := service.HashPassword("testpassword")
			assert.NoError(t, err)

			cost, err := bcrypt.Cost([]byte(hash))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCost, cost)
		})
	}
}
