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
				os.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				os.Setenv(rateWatchAddrEnvKey, "rw:8080")
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(portEnvKey, "3030")
				os.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				os.Setenv(mailerFromAddrEnvKey, "from@from.com")
			},
			want: &Config{
				MailerAddr:      "mailer:80",
				RateWatcherAddr: "rw:8080",
				LogLevel:        "debug",
				Port:            "3030",
				Dsn:             "mysql://test:tests@(db:testse)/shgsoh",
				MailerFromAddr:  "from@from.com",
			},
			wantErr: false,
		},
		{
			name: "Should not parse config when mailer addr is missing",
			setup: func() {
				os.Setenv(mailerServiceAddrEnvKey, "")
				os.Setenv(rateWatchAddrEnvKey, "rw:8080")
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(portEnvKey, "3030")
				os.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				os.Setenv(mailerFromAddrEnvKey, "from@from.com")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when rw addr is missing",
			setup: func() {
				os.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				os.Setenv(rateWatchAddrEnvKey, "")
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(portEnvKey, "3030")
				os.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				os.Setenv(mailerFromAddrEnvKey, "from@from.com")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when log level is missing",
			setup: func() {
				os.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				os.Setenv(rateWatchAddrEnvKey, "rw:2209")
				os.Setenv(logLevelEnvKey, "")
				os.Setenv(portEnvKey, "3030")
				os.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				os.Setenv(mailerFromAddrEnvKey, "from@from.com")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when port is missing",
			setup: func() {
				os.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				os.Setenv(rateWatchAddrEnvKey, "rw:2209")
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(portEnvKey, "")
				os.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				os.Setenv(mailerFromAddrEnvKey, "from@from.com")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when dsn is missing",
			setup: func() {
				os.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				os.Setenv(rateWatchAddrEnvKey, "rw:2209")
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(portEnvKey, "801")
				os.Setenv(dsnEnvKey, "")
				os.Setenv(mailerFromAddrEnvKey, "from@from.com")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when mailer from is missing",
			setup: func() {
				os.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				os.Setenv(rateWatchAddrEnvKey, "rw:2209")
				os.Setenv(logLevelEnvKey, "debug")
				os.Setenv(portEnvKey, "2424")
				os.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				os.Setenv(mailerFromAddrEnvKey, "")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				os.Unsetenv(mailerServiceAddrEnvKey)
				os.Unsetenv(rateWatchAddrEnvKey)
				os.Unsetenv(logLevelEnvKey)
				os.Unsetenv(portEnvKey)
				os.Unsetenv(dsnEnvKey)
				os.Unsetenv(mailerFromAddrEnvKey)
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
