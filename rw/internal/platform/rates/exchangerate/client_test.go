package exchangerate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	type args struct {
		token string
		url   string
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{
			name: "Should initiate new client correctly",
			args: args{
				token: "tokeeeen",
				url:   "https://url.com",
			},
			want: Client{
				token: "tokeeeen",
				url:   "https://url.com",
			},
		},
		{
			name: "Should initiate new client correctly",
			args: args{
				token: "tok352533))____$",
				url:   "https://url2.com",
			},
			want: Client{
				token: "tok352533))____$",
				url:   "https://url2.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClient(tt.args.token, tt.args.url)
			require.Equal(t, tt.want, got)
		})
	}
}
