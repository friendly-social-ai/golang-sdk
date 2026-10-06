package sdk

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestGetFileURL_Success(t *testing.T) {
	descriptor := &FileDescriptor{Id: MockFileId(1), AccessHash: MockFileAccessHash("hash")}
	client := NewClient()

	url, err := client.GetFileURL(descriptor)
	require.NoError(t, err)
	require.Equal(t, "https://api.getfriend.ly/files/download/1/hash", url)
}

func TestGetFileURL_InvalidURL(t *testing.T) {
	descriptor := &FileDescriptor{Id: MockFileId(1), AccessHash: MockFileAccessHash("hash")}
	client := NewClient().
		WithBaseURL("::invalid")

	_, err := client.GetFileURL(descriptor)
	require.Error(t, err)
}

func TestUploadFile_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/files/upload", r.URL.Path)
		require.Equal(t, "11", r.Header.Get("X-File-Size"))
		require.Equal(t, "1", r.Header.Get("X-User-Id"))
		require.Equal(t, "token", r.Header.Get("X-Token"))

		err := r.ParseMultipartForm(1 << 20)
		require.NoError(t, err)
		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close() //nolint:errcheck

		require.Equal(t, "file.txt", header.Filename)
		data, err := io.ReadAll(file)
		require.NoError(t, err)
		require.Equal(t, "hello world", string(data))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"accessHash": "hash"
		}`))
	}))
	defer ts.Close()

	client := NewClient().
		WithBaseURL(ts.URL)

	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	file := strings.NewReader("hello world")
	fd, err := client.UploadFile(context.Background(), auth, "file.txt", file, 11)

	require.NoError(t, err)
	require.Equal(t, &FileDescriptor{Id: MockFileId(123), AccessHash: MockFileAccessHash("hash")}, fd)
}

func TestUploadFile_Canceled(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/files/upload").
		Reply(200).
		Delay(100 * time.Hour)

	client := NewClient()
	file := strings.NewReader("hello world")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.UploadFile(ctx, nil, "file.txt", file, 11)
	require.Error(t, err)
}

func TestUploadFile_AlreadyClosed(t *testing.T) {
	client := NewClient()
	file := strings.NewReader("hello world")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.UploadFile(ctx, nil, "file.txt", file, 11)
	require.Error(t, err)
}

func TestUploadFile_InvalidURL(t *testing.T) {
	client := NewClient().
		WithBaseURL("::invalid")

	file := strings.NewReader("hello world")
	_, err := client.UploadFile(context.Background(), nil, "file.txt", file, 11)

	require.Error(t, err)
}

func TestUploadFile_NewRequestFailed(t *testing.T) {
	client := NewClient()
	file := strings.NewReader("hello world")

	_, err := client.UploadFile(nil, nil, "file.txt", file, 11) //nolint:staticcheck
	require.Error(t, err)
}

func TestDownloadFile_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/files/download/123/hash", r.URL.Path)
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte("hello"))
	}))
	defer ts.Close() //nolint:errcheck

	client := NewClient().
		WithBaseURL(ts.URL)

	fd := &FileDescriptor{Id: MockFileId(123), AccessHash: MockFileAccessHash("hash")}
	reader, err := client.DownloadFile(context.Background(), fd)
	require.NoError(t, err)

	defer reader.Close() //nolint:errcheck
	data, err := io.ReadAll(reader)

	require.NoError(t, err)
	require.Equal(t, "hello", string(data))
}

func TestDownloadFile_InvalidURL(t *testing.T) {
	fd := &FileDescriptor{Id: MockFileId(1), AccessHash: MockFileAccessHash("hash")}
	client := NewClient().
		WithBaseURL("::invalid")

	_, err := client.DownloadFile(context.Background(), fd)
	require.Error(t, err)
}

func TestDownloadFile_Cancel(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/files/download/123/hash").
		Reply(200).
		Delay(100 * time.Hour)

	client := NewClient()
	fd := &FileDescriptor{Id: MockFileId(123), AccessHash: MockFileAccessHash("hash")}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.DownloadFile(ctx, fd)
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestDownloadFile_Error(t *testing.T) {
	cases := []struct {
		name string
		code int
		body []byte
	}{
		{"Unauthorized", 401, []byte("invalid auth")},
		{"Forbidden", 403, []byte("you are not admin")},
		{"Not Found", 404, []byte("not found")},
		{"I'm a teapot", 418, []byte("lol")},
	}

	for _, tc := range cases {
		defer gock.Off()

		gock.New("https://api.getfriend.ly").
			Get("/files/download/123/hash").
			Reply(tc.code).
			Body(io.NopCloser(strings.NewReader(string(tc.body))))

		client := NewClient()
		fd := &FileDescriptor{Id: MockFileId(123), AccessHash: MockFileAccessHash("hash")}
		_, err := client.DownloadFile(context.Background(), fd)

		var apiError APIError
		require.ErrorAs(t, err, &apiError)
		require.Equal(t, tc.code, apiError.Code)
		require.Equal(t, tc.body, apiError.Body)
	}
}

func TestDownloadFile_NewRequestFailed(t *testing.T) {
	client := NewClient()
	fd := &FileDescriptor{Id: MockFileId(123), AccessHash: MockFileAccessHash("hash")}

	_, err := client.DownloadFile(nil, fd) //nolint:staticcheck
	require.Error(t, err)
}

func TestDownloadFile_FailedToReadError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/files/download/123/hash").
		Reply(418).
		Map(func(resp *http.Response) *http.Response {
			resp.Body = io.NopCloser(&errorReader{})
			return resp
		})

	client := NewClient()
	fd := &FileDescriptor{Id: MockFileId(123), AccessHash: MockFileAccessHash("hash")}
	_, err := client.DownloadFile(context.Background(), fd)
	require.Error(t, err)
}

func TestPreuploadFile_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/files/preupload", r.URL.Path)
		require.Equal(t, "11", r.Header.Get("X-File-Size"))
		require.Empty(t, r.Header.Get("X-User-Id"))

		err := r.ParseMultipartForm(1 << 20)
		require.NoError(t, err)
		_, header, err := r.FormFile("file")
		require.NoError(t, err)
		require.Equal(t, "avatar.png", header.Filename)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":123,"accessHash":"hash"}`))
	}))
	defer ts.Close()

	client := NewClient().WithBaseURL(ts.URL)
	fd, err := client.PreuploadFile(context.Background(), nil, "avatar.png", strings.NewReader("hello world"), 11)

	require.NoError(t, err)
	require.Equal(t, &FileDescriptor{Id: MockFileId(123), AccessHash: MockFileAccessHash("hash")}, fd)
}

func TestPreuploadFile_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/files/preupload").
		Reply(400)

	client := NewClient()
	_, err := client.PreuploadFile(context.Background(), nil, "avatar.png", strings.NewReader("hello world"), 11)
	require.Error(t, err)
}
