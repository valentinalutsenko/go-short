package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handlePost(t *testing.T) {
	storage = make(map[string]string)
	pattern := `^http://localhost:8080/[a-zA-Z0-9]{5}$`

	tests := []struct {
		name        string
		contentType string
		body        string
		statusCode  int
	}{
		{
			name:        "Success case",
			contentType: "text/plain",
			body:        "https://yandex.ru",
			statusCode:  http.StatusCreated,
		},
		{
			name:        "Wrong Content-Type",
			contentType: "application/json",
			body:        "https://yandex.ru",
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "Empty body",
			contentType: "text/plain",
			body:        "  ",
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "Invalid URL (no schema)",
			contentType: "text/plain",
			body:        "google.com",
			statusCode:  http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.contentType)

			w := httptest.NewRecorder()
			handlePost(w, request)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, test.statusCode, res.StatusCode)

			if test.statusCode == http.StatusCreated {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Regexp(t, pattern, string(resBody))
				assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
			}

		})
	}
}
