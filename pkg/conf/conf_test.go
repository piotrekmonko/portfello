package conf

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		wantErr bool
		config  *Config
	}{
		{
			wantErr: true,
			config:  &Config{},
		},
		{
			wantErr: true,
			config:  &Config{DatabaseDSN: "some dsn"},
		},
		{
			wantErr: true,
			config:  &Config{DatabaseDSN: "some dsn", Auth: Auth0{Provider: "invalid provider"}},
		},
		{
			wantErr: true,
			config:  &Config{DatabaseDSN: "some dsn", Auth: Auth0{Provider: AuthProviderAuth0}},
		},
		{
			wantErr: false,
			config:  &Config{DatabaseDSN: "some dsn", Auth: Auth0{Provider: AuthProviderAuth0, ClientID: "123"}},
		},
		{
			wantErr: true,
			config:  &Config{DatabaseDSN: "some dsn", Auth: Auth0{Provider: AuthProviderLocal}},
		},
		{
			wantErr: false,
			config:  &Config{DatabaseDSN: "some dsn", Auth: Auth0{Provider: AuthProviderLocal, ClientSecret: "123"}},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("validation test %d", i), func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTestConfig(t *testing.T) {
	tc := NewTestConfig()
	if tc == nil {
		t.Errorf("NewTestConfig() must not return nil")
	}
}

func TestNew_fromEnv(t *testing.T) {
	prefix := "TEST_PRTFL"
	envs := map[string]string{
		prefix + "_DATABASE_DSN":  "db-dsn",
		prefix + "_AUTH.PROVIDER": "auth-provider",
	}
	for k, v := range envs {
		require.Nil(t, os.Setenv(k, v))
		e, ok := os.LookupEnv(k)
		assert.True(t, ok)
		assert.Equal(t, v, e)
	}
	require.Nil(t, InitConfig("", prefix))
	c := New()
	assert.Equal(t, envs[prefix+"_DATABASE_DSN"], c.DatabaseDSN)
	assert.Equal(t, envs[prefix+"_AUTH_PROVIDER"], c.Auth.Provider)
}
