package todi

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadFile uploads a file for use in comments.
func (c *Client) UploadFile(ctx context.Context, path, name, projectID string) (Upload, []byte, error) {
	var upload Upload
	file, err := os.Open(path)
	if err != nil {
		return Upload{}, nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if projectID != "" {
		if err := writer.WriteField("project_id", projectID); err != nil {
			return Upload{}, nil, err
		}
	}
	if name != "" {
		if err := writer.WriteField("file_name", name); err != nil {
			return Upload{}, nil, err
		}
	}

	fileName := name
	if fileName == "" {
		fileName = filepath.Base(path)
	}
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return Upload{}, nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return Upload{}, nil, err
	}
	if err := writer.Close(); err != nil {
		return Upload{}, nil, err
	}

	fullURL, err := c.url("/api/v1/uploads", nil)
	if err != nil {
		return Upload{}, nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, body)
	if err != nil {
		return Upload{}, nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	raw, err := c.do(req, &upload)
	return upload, raw, err
}

// DeleteUpload deletes a file upload by file URL.
func (c *Client) DeleteUpload(ctx context.Context, fileURL string) ([]byte, error) {
	params := map[string]string{"file_url": fileURL}
	return c.deleteWithParams(ctx, "/api/v1/uploads", params)
}
