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

			got := Must(tt.args.cfg, tt.args.err)
			require.Equal(t, tt.want, got)
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
				require.NoError(t, os.Setenv(addrEnvKey, "0.0.0.0:80"))
				require.NoError(t, os.Setenv(rateWatchAddrEnvKey, "rw:3333"))
				require.NoError(t, os.Setenv(subServiceAddrEnvKey, "ss:6666"))
			},
			want: &Config{
				LogLevel:        "debug",
				Addr:            "0.0.0.0:80",
				RateWatcherAddr: "rw:3333",
				SubAddr:         "ss:6666",
			},
			wantErr: false,
		},
		{
			name: "Should not parse config when log level is missing",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.Setenv(logLevelEnvKey, ""))
				require.NoError(t, os.Setenv(addrEnvKey, "0.0.0.0:80"))
				require.NoError(t, os.Setenv(rateWatchAddrEnvKey, "rw:3333"))
				require.NoError(t, os.Setenv(subServiceAddrEnvKey, "ss:6666"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when addr is missing",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.Setenv(logLevelEnvKey, "debug"))
				require.NoError(t, os.Setenv(addrEnvKey, ""))
				require.NoError(t, os.Setenv(rateWatchAddrEnvKey, "rw:3333"))
				require.NoError(t, os.Setenv(subServiceAddrEnvKey, "ss:6666"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when rw addr is missing",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.Setenv(logLevelEnvKey, "debug"))
				require.NoError(t, os.Setenv(addrEnvKey, "0.0.0.0:80"))
				require.NoError(t, os.Setenv(rateWatchAddrEnvKey, ""))
				require.NoError(t, os.Setenv(subServiceAddrEnvKey, "ss:6666"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when ss addr is missing",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.Setenv(logLevelEnvKey, "debug"))
				require.NoError(t, os.Setenv(addrEnvKey, "0.0.0.0:80"))
				require.NoError(t, os.Setenv(rateWatchAddrEnvKey, "rw:3333"))
				require.NoError(t, os.Setenv(subServiceAddrEnvKey, ""))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				require.NoError(t, os.Unsetenv(logLevelEnvKey))
				require.NoError(t, os.Unsetenv(addrEnvKey))
				require.NoError(t, os.Unsetenv(rateWatchAddrEnvKey))
				require.NoError(t, os.Unsetenv(subServiceAddrEnvKey))
			})

			tt.setup(t)
			got, err := NewFromEnv()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}
