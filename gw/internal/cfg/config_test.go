package cfg

import (
	"errors"
	"os"
	"reflect"
	"testing"
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
				defer func() {
					if recover() == nil {
						t.Fatal("Expected to panic")
					}
				}()
			}
			if got := Must(tt.args.cfg, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Must() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		want    *Config
		wantErr bool
	}{
		{
			name: "Should parse config correctly when all env vars are present",
			setup: func() {
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(addrEnvKey, "0.0.0.0:80")
				os.Setenv(rateWatchAddrEnvKey, "rw:3333")
				os.Setenv(subServiceAddrEnvKey, "ss:6666")
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
			setup: func() {
				os.Setenv(logLevelEnvKey, "")
				os.Setenv(addrEnvKey, "0.0.0.0:80")
				os.Setenv(rateWatchAddrEnvKey, "rw:3333")
				os.Setenv(subServiceAddrEnvKey, "ss:6666")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when addr is missing",
			setup: func() {
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(addrEnvKey, "")
				os.Setenv(rateWatchAddrEnvKey, "rw:3333")
				os.Setenv(subServiceAddrEnvKey, "ss:6666")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when rw addr is missing",
			setup: func() {
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(addrEnvKey, "0.0.0.0:80")
				os.Setenv(rateWatchAddrEnvKey, "")
				os.Setenv(subServiceAddrEnvKey, "ss:6666")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when ss addr is missing",
			setup: func() {
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(addrEnvKey, "0.0.0.0:80")
				os.Setenv(rateWatchAddrEnvKey, "rw:3333")
				os.Setenv(subServiceAddrEnvKey, "")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				os.Unsetenv(logLevelEnvKey)
				os.Unsetenv(addrEnvKey)
				os.Unsetenv(rateWatchAddrEnvKey)
				os.Unsetenv(subServiceAddrEnvKey)
			})

			tt.setup()
			got, err := NewFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
