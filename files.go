package sdk

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// FileDescriptor is a helper structure for composing file's ID and AccessHash.
type FileDescriptor struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

type uploadFileResponse struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

// GetFileURL returns file access URL for corresponding descriptor.
func (c *Client) GetFileURL(fd *FileDescriptor) (string, error) {
	url, err := url.JoinPath(c.url, fmt.Sprintf("/files/download/%d/%s", fd.Id.value, fd.AccessHash.value))
	if err != nil {
		return "", fmt.Errorf("invalid path: %s + %s", c.url, fmt.Sprintf("/files/download/%d/%s", fd.Id.value, fd.AccessHash.value))
	}

	return url, nil
}

// DownloadFile opens connection for downloading file and returns corresponding io.ReadCloser.
func (c *Client) DownloadFile(ctx context.Context, fd *FileDescriptor) (io.ReadCloser, error) {
	url, err := c.GetFileURL(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to get file URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp.Body, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP error body")
	}

	return nil, fmt.Errorf("unexpected response: %w", APIError{Code: resp.StatusCode, Body: body})
}

// UploadFile uploads file of size bytes from io.Reader to the server on behalf of Authorization and returns
// corresponding descriptor. It accepts filename by which file will be saved on server, and reader from which file
// will be read.
func (c *Client) UploadFile(ctx context.Context, auth *Authorization, filename string, reader io.Reader, size int64) (*FileDescriptor, error) {
	fd, err := c.sendFile(ctx, auth, "/files/upload", filename, reader, size)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return fd, nil
}

// PreuploadFile uploads file the same way as UploadFile, but for use before the account exists, like the avatar
// passed to Register. Authorization may be nil.
func (c *Client) PreuploadFile(ctx context.Context, auth *Authorization, filename string, reader io.Reader, size int64) (*FileDescriptor, error) {
	fd, err := c.sendFile(ctx, auth, "/files/preupload", filename, reader, size)
	if err != nil {
		return nil, fmt.Errorf("failed to preupload file: %w", err)
	}

	return fd, nil
}

// sendFile posts file of size bytes from reader as multipart form to endpoint and returns its descriptor.
func (c *Client) sendFile(ctx context.Context, auth *Authorization, endpoint, filename string, reader io.Reader, size int64) (*FileDescriptor, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	filename = path.Base(filename)

	go func() {
		var err error
		defer func() {
			_ = writer.Close()
			_ = pw.CloseWithError(err)
		}()

		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			return
		}

		_, err = io.Copy(part, reader)
	}()

	completePath, err := url.JoinPath(c.url, endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s + %s", c.url, endpoint)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", completePath, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-File-Size", strconv.FormatInt(size, 10))
	authorize(req, auth)

	var resp uploadFileResponse
	if err := c.execute(req, &resp); err != nil {
		return nil, err
	}

	return &FileDescriptor{
		Id:         resp.Id,
		AccessHash: resp.AccessHash,
	}, nil
}
