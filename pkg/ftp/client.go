package ftp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/RobertMNewton/bambu-golang-api/pkg/types/config"
	"github.com/jlaffaye/ftp"
)

const (
	ftpPort        = 990
	connectTimeout = 10 * time.Second
)

// Client represents an FTP client for file transfers
type Client struct {
	config config.PrinterConfig
	conn   *ftp.ServerConn
}

func NewClient(config config.PrinterConfig) *Client {
	return &Client{
		config: config,
	}
}

func (client *Client) Connect(ctx context.Context) error {
	// Create a channel for the connection result
	connChan := make(chan error, 1)

	go func() {
		addr := fmt.Sprintf("%s:%d", client.config.GetDeviceIPAddress(), ftpPort)

		conn, err := ftp.Dial(addr, ftp.DialWithTimeout(connectTimeout))
		if err != nil {
			connChan <- fmt.Errorf("ftp dial failed: %w", err)
			return
		}

		if err := conn.Login("bblp", client.config.GetDeviceAccessCode()); err != nil {
			conn.Quit()
			connChan <- fmt.Errorf("ftp login failed: %w", err)
			return
		}

		client.conn = conn
		connChan <- nil
	}()

	select {
	case err := <-connChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("connection timeout: %w", ctx.Err())
	}
}

func (client *Client) UploadFile(ctx context.Context, localPath, remotePath string) error {
	if client.conn == nil {
		return fmt.Errorf("not connected")
	}

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer file.Close()

	remoteDir := filepath.Dir(remotePath)
	if err := client.createDirectory(remoteDir); err != nil {
		return fmt.Errorf("failed to create remote directory: %w", err)
	}

	err = client.conn.Stor(remotePath, file)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (client *Client) ListFiles(ctx context.Context, path string) ([]string, error) {
	if client.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	entries, err := client.conn.List(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFile {
			files = append(files, entry.Name)
		}
	}

	return files, nil
}

func (client *Client) Disconnect() error {
	if client.conn != nil {
		return client.conn.Quit()
	}
	return nil
}

func (client *Client) createDirectory(path string) error {
	if err := client.conn.ChangeDir(path); err == nil {
		return client.conn.ChangeDir("/")
	}

	// Directory doesn't exist, create it
	if err := client.conn.MakeDir(path); err != nil {
		return err
	}

	return nil
}
