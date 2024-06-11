package cfg

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMust(t *testing.T) {
	t.Parallel()
	type args struct {
		cfg *Config
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Should not panic when err is nil",
			args: args{
				cfg: &Config{},
				err: nil,
			},
			want: &Config{},
		},
		{
			name: "Should panic when err is not nil",
			args: args{
				cfg: nil,
				err: errors.New("failed to parse config"),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.wantErr {
				require.Panics(t, func() {
					Must(tt.args.cfg, tt.args.err)
				})
				return
			}

			require.NotPanics(t, func() {
				got := Must(tt.args.cfg, tt.args.err)
				require.Equal(t, tt.want, got)
			})
		})
	}
}

func TestNewFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		want    *Config
		wantErr bool
	}{
		{
			name: "Should parse config correctly when all env vars are present",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.Setenv(logLevelEnvKey, "debug"))
				require.NoError(t, os.Setenv(portEnvKey, "80"))
				require.NoError(t, os.Setenv(mailerTokenEnvKey, "secret"))
			},
			want: &Config{
				LogLevel:    "debug",
				Port:        "80",
				MailerToken: "secret",
			},
			wantErr: false,
		},
		{
			name: "Should not parse config correctly when log level is missing",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.Setenv(logLevelEnvKey, ""))
				require.NoError(t, os.Setenv(portEnvKey, "80"))
				require.NoError(t, os.Setenv(mailerTokenEnvKey, "secret"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config correctly when port is missing",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.Setenv(logLevelEnvKey, "debug"))
				require.NoError(t, os.Setenv(portEnvKey, ""))
				require.NoError(t, os.Setenv(mailerTokenEnvKey, "secret"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config correctly when token is missing",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.Setenv(logLevelEnvKey, "debug"))
				require.NoError(t, os.Setenv(portEnvKey, "80"))
				require.NoError(t, os.Setenv(mailerTokenEnvKey, ""))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				require.NoError(t, os.Unsetenv(logLevelEnvKey))
				require.NoError(t, os.Unsetenv(portEnvKey))
				require.NoError(t, os.Unsetenv(mailerTokenEnvKey))
			})

			tt.setup(t)
			got, err := NewFromEnv()
			require.Equal(t, tt.want, got)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
