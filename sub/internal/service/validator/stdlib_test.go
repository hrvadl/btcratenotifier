package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStdlibValidate(t *testing.T) {
	t.Parallel()
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should parse correct email",
			args: args{
				email: "email@gmail.com",
			},
			want: true,
		},
		{
			name: "Should parse correct email",
			args: args{
				email: "hrvadleo@gmail.com",
			},
			want: true,
		},
		{
			name: "Should parse correct email",
			args: args{
				email: "test@gmail.com",
			},
			want: true,
		},
		{
			name: "Should not parse incorrect email",
			args: args{
				email: "test@gmail.",
			},
			want: false,
		},
		{
			name: "Should not parse incorrect email",
			args: args{
				email: "test@.com",
			},
			want: false,
		},
		{
			name: "Should not parse incorrect email",
			args: args{
				email: "@gmail.com",
			},
			want: false,
		},
		{
			name: "Should not parse incorrect email",
			args: args{
				email: "",
			},
			want: false,
		},
		{
			name: "Should not parse incorrect email",
			args: args{
				email: "youngwwad",
			},
			want: false,
		},
		{
			name: "Should not parse incorrect email",
			args: args{
				email: "",
			},
			want: false,
		},
		{
			name: "Should not parse incorrect email",
			args: args{
				email: "youngwwad@m",
			},
			want: false,
		},
		{
			name: "Should parse correct email with subdomain",
			args: args{
				email: "youngwwad@m.c.c",
			},
			want: true,
		},
		{
			name: "Should parse correct email",
			args: args{
				email: "youngwwad@m.c",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := Stdlib{}
			got := r.Validate(tt.args.email)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNewStdlib(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want *Stdlib
	}{
		{
			name: "Should initialize validator correctly",
			want: &Stdlib{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewStdlib()
			require.Equal(t, tt.want, got)
		})
	}
}
