package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"testing"

	"github.com/fuskovic/networker/internal/test"
	"github.com/stretchr/testify/require"
)

func TestRequestCrafting(t *testing.T) {
	t.Run("ShouldFail", func(t *testing.T) {
		t.Parallel()
		t.Run("if request URL is not specified", func(t *testing.T) {
			t.Parallel()
			_, err := NewNetworkerCraftedHTTPRequest(
				&Config{URL: "", Method: http.MethodGet},
			)
			require.Error(t, err)
			require.Equal(t, err, errUrlUnset)
		})
		t.Run("if request URL is missing protocol", func(t *testing.T) {
			t.Parallel()
			_, err := NewNetworkerCraftedHTTPRequest(
				&Config{URL: "urlwithmissingprotocol.com", Method: http.MethodGet},
			)
			require.Error(t, err)
			require.Equal(t, err, errProtocolUnset)
		})
		t.Run("if request method is invalid", func(t *testing.T) {
			t.Parallel()
			_, err := NewNetworkerCraftedHTTPRequest(
				&Config{URL: "http://validurl.com", Method: "invalid"},
			)
			require.Error(t, err)
			require.Equal(t, err, errInvalidRequestMethod)
		})
		t.Run("if file extension for provisioning request body is not .json", func(t *testing.T) {
			t.Parallel()
			_, err := NewNetworkerCraftedHTTPRequest(
				&Config{
					URL:    "http://validurl.com",
					Method: http.MethodPost,
					Body:   "unsupported.txt",
				},
			)
			require.Error(t, err)
			require.Equal(t, err,
				fmt.Errorf(
					"failed to add request body: %w",
					errUnsupportedFileExtension,
				),
			)
		})
		t.Run("if multipart form data arg does not contain an equals sign character", func(t *testing.T) {
			t.Parallel()
			_, err := NewNetworkerCraftedHTTPRequest(
				&Config{
					URL:           "http://validurl.com",
					Method:        http.MethodPost,
					MultiPartForm: `"missingequalssign"`,
				},
			)
			require.Error(t, err)
			require.Equal(t, err,
				fmt.Errorf(
					"failed to add multi-part form data to request: %w",
					errInvalidMultiPartFormDataFormat,
				),
			)
		})
		t.Run("if multipart form data arg does not designate any files", func(t *testing.T) {
			t.Parallel()
			_, err := NewNetworkerCraftedHTTPRequest(
				&Config{
					URL:           "http://validurl.com",
					Method:        http.MethodPost,
					MultiPartForm: "formname=",
				},
			)
			require.Error(t, err)
			require.Equal(t, err,
				fmt.Errorf(
					"failed to add multi-part form data to request: %w",
					errNoUploadFilesDesignated,
				),
			)
		})
		t.Run("if multipart form data arg does not designate a formname", func(t *testing.T) {
			t.Parallel()
			_, err := NewNetworkerCraftedHTTPRequest(
				&Config{
					URL:           "http://validurl.com",
					Method:        http.MethodPost,
					MultiPartForm: "=file.png",
				},
			)
			require.Error(t, err)
			require.Equal(t, err,
				fmt.Errorf(
					"failed to add multi-part form data to request: %w",
					errMultiPartFormNameUnset,
				),
			)
		})
		t.Run("if multipart form data arg designates files that dont exist", func(t *testing.T) {
			t.Parallel()
			_, err := NewNetworkerCraftedHTTPRequest(
				&Config{
					URL:           "http://validurl.com",
					Method:        http.MethodPost,
					MultiPartForm: `"formname=doesntexist.png"`,
				},
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), "no such file or directory")
		})
		t.Run("if uneven number of key/value pairs for request headers is specified", func(t *testing.T) {
			t.Parallel()
			_, err := NewNetworkerCraftedHTTPRequest(
				&Config{
					URL:     "http://validurl.com",
					Headers: []string{"key:value", "key"},
					Method:  http.MethodPost,
				},
			)
			require.Error(t, err)
			require.Equal(t, err,
				fmt.Errorf("failed to add headers: %w",
					errUnevenNumberOfHeaderKeyValuePairs,
				))
		})
	})
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithMockServer(t, func(t *testing.T, testserverURL string) {
			// create a new object using a JSON string literal
			cfg := &Config{
				Headers: []string{"auth:doesntmatter"},
				URL:     testserverURL,
				Method:  http.MethodPost,
				Body:    `{"field": "doesntmatter"}`,
			}
			req, err := NewNetworkerCraftedHTTPRequest(cfg)
			require.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, resp.StatusCode)
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			var createdObject test.MockObject
			require.NoError(t, json.Unmarshal(data, &createdObject))

			// get the object
			cfg.URL = fmt.Sprintf("%s?id=%d", cfg.URL, createdObject.ID)
			cfg.Method = http.MethodGet
			cfg.Body = ""
			req, err = NewNetworkerCraftedHTTPRequest(cfg)
			require.NoError(t, err)
			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			defer resp.Body.Close()
			var newObject test.MockObject
			data, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(data, &newObject))
			require.Equal(t, createdObject, newObject)

			// delete the object
			cfg.Method = http.MethodDelete
			req, err = NewNetworkerCraftedHTTPRequest(cfg)
			require.NoError(t, err)
			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			// make sure its gone
			cfg.Method = http.MethodGet
			req, err = NewNetworkerCraftedHTTPRequest(cfg)
			require.NoError(t, err)
			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusNotFound, resp.StatusCode)

			root := test.ProjectRoot(t)

			// create a new object with a json file to use for the request body
			cfg = &Config{
				URL:     testserverURL,
				Headers: []string{"auth:doesntmatter"},
				Method:  http.MethodPost,
				Body:    path.Join(root, "internal/test/body.json"),
			}
			req, err = NewNetworkerCraftedHTTPRequest(cfg)
			require.NoError(t, err)
			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, resp.StatusCode)
			defer resp.Body.Close()
			data, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			newlyCreatedObject := new(test.MockObject)
			require.NoError(t, json.Unmarshal(data, newlyCreatedObject))
			require.True(t, newlyCreatedObject.ID != 0)
			require.Equal(t, "name", newlyCreatedObject.Field)

			file1 := path.Join(root, "internal/test/cat_1.jpeg")
			file2 := path.Join(root, "internal/test/cat_2.jpeg")
			multiPartFormUploadArg := fmt.Sprintf("files=%s,%s", file1, file2)

			// upload
			cfg = &Config{
				URL:           testserverURL,
				Headers:       []string{"auth:doesntmatter"},
				Method:        http.MethodPut,
				MultiPartForm: multiPartFormUploadArg,
			}
			req, err = NewNetworkerCraftedHTTPRequest(cfg)
			require.NoError(t, err)
			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, resp.StatusCode)
		})
	})
}
