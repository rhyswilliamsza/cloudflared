package tunnel

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestHostnameFromURI(t *testing.T) {
	assert.Equal(t, "awesome.warptunnels.horse:22", hostnameFromURI("ssh://awesome.warptunnels.horse:22"))
	assert.Equal(t, "awesome.warptunnels.horse:22", hostnameFromURI("ssh://awesome.warptunnels.horse"))
	assert.Equal(t, "awesome.warptunnels.horse:2222", hostnameFromURI("ssh://awesome.warptunnels.horse:2222"))
	assert.Equal(t, "localhost:3389", hostnameFromURI("rdp://localhost"))
	assert.Equal(t, "localhost:3390", hostnameFromURI("rdp://localhost:3390"))
	assert.Equal(t, "", hostnameFromURI("trash"))
	assert.Equal(t, "", hostnameFromURI("https://awesomesauce.com"))
}

func TestShouldRunQuickTunnel(t *testing.T) {
	tests := []struct {
		name        string
		flags       map[string]string
		expectError bool
	}{
		{
			name:        "Quick tunnel with URL set",
			flags:       map[string]string{"url": "http://127.0.0.1:8080", "quick-service": "https://fakeapi.trycloudflare.com"},
			expectError: false,
		},
		{
			name:        "Quick tunnel with unix-socket set",
			flags:       map[string]string{"unix-socket": "/tmp/socket", "quick-service": "https://fakeapi.trycloudflare.com"},
			expectError: false,
		},
		{
			name:        "Quick tunnel with hello-world flag",
			flags:       map[string]string{"hello-world": "true", "quick-service": "https://fakeapi.trycloudflare.com"},
			expectError: false,
		},
		{
			name:        "Quick tunnel with proxy-dns (invalid combo)",
			flags:       map[string]string{"url": "http://127.0.0.1:9090", "proxy-dns": "true", "quick-service": "https://fakeapi.trycloudflare.com"},
			expectError: true,
		},
		{
			name:        "No quick-service set",
			flags:       map[string]string{"url": "http://127.0.0.1:9090"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock RunQuickTunnel Function
			mockCalled := false
			mockQuickTunnelRunner := func(sc *subcommandContext) error {
				mockCalled = true
				return nil
			}

			// Mock App Context
			app := &cli.App{}
			set := flagSetFromMap(tt.flags)
			context := cli.NewContext(app, set, nil)

			// Call TunnelCommand
			err := tunnelCommandImpl(context, mockQuickTunnelRunner)

			// Validate
			if tt.expectError {
				assert.False(t, mockCalled)
				require.Error(t, err)
			} else {
				assert.True(t, mockCalled)
				require.NoError(t, err)
			}
		})
	}
}

func flagSetFromMap(flags map[string]string) *flag.FlagSet {
	set := flag.NewFlagSet("test", 0)
	for key, value := range flags {
		set.String(key, "", "")
		set.Set(key, value)
	}
	return set
}
