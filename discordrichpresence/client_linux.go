//go:build !windows

package discordrichpresence

import (
	"fmt"
	"net"
	"os"
	"time"
)

// connectToPipe connects to the appropriate IPC mechanism based on the platform.
func (c *Client) connectToPipe(pipeName string) (net.Conn, error) {
	// Unix-like systems use Unix domain sockets
	address := c.getSocketAddress(pipeName)
	return net.DialTimeout("unix", address, 2*time.Second)
}

// getSocketAddress returns the appropriate socket address for Unix-like systems.
func (c *Client) getSocketAddress(pipeName string) string {
	tmpDir := os.Getenv("XDG_RUNTIME_DIR")
	if tmpDir == "" {
		tmpDir = os.Getenv("TMPDIR")
	}
	if tmpDir == "" {
		tmpDir = "/tmp"
	}

	if _, err := os.Stat(tmpDir + "/app/com.discordapp.Discord/" + pipeName); err == nil {
		return tmpDir + "/app/com.discordapp.Discord/" + pipeName
	}

	if _, err := os.Stat(tmpDir + "/.flatpak/com.discordapp.Discord/xdg-run/" + pipeName); err == nil {
		return tmpDir + "/.flatpak/com.discordapp.Discord/xdg-run/" + pipeName
	}

	return fmt.Sprintf("%s/%s", tmpDir, pipeName)
}
