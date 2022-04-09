package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"testing"

	"github.com/fuskovic/networker/internal/test"
	"github.com/stretchr/testify/require"
)

func TestRequestCommand(t *testing.T) {
	t.Run("ShouldFail", func(t *testing.T) {
		test.WithNetworker(t, "url is not provided", func(t *testing.T) {
			cmd := exec.Command("networker", "request")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), `accepts 1 arg(s), received 0`)
		})
		test.WithNetworker(t, "protocol is not included in url", func(t *testing.T) {
			cmd := exec.Command("networker", "request", "google.com")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "no protocol specified in url endpoint")
		})
		test.WithNetworker(t, "invalid request method", func(t *testing.T) {
			cmd := exec.Command("networker", "request", "-m", "invalid", "https://google.com")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "invalid request method")
		})
		test.WithNetworker(t, "unsupported file extension designated for request body", func(t *testing.T) {
			cmd := exec.Command("networker", "request", "-m", "post", "-b", "unsupported.txt", "https://google.com")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "unsupported file extension")
		})
		test.WithNetworker(t, "json file designated to use for request body doesnt exist", func(t *testing.T) {
			cmd := exec.Command("networker", "request", "-m", "post", "-b", "doesntexist.json", "https://google.com")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "no such file or directory")
		})
		test.WithNetworker(t, "multi-part form arg missing equals sign", func(t *testing.T) {
			cmd := exec.Command("networker", "request", "-m", "post", "--files", "formname", "https://google.com")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "invalid multipart form upload format")
		})
		test.WithNetworker(t, "multi-part form arg doesnt designate any files", func(t *testing.T) {
			cmd := exec.Command("networker", "request", "-m", "post", "--files", "formname=", "https://google.com")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "no upload files designated")
		})
		test.WithNetworker(t, "multi-part form arg doesnt designate a formname", func(t *testing.T) {
			cmd := exec.Command("networker", "request", "-m", "post", "--files", "=file.png", "https://google.com")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "no multipart form name specified")
		})
		test.WithNetworker(t, "multi-part form arg designates a file that does not exist", func(t *testing.T) {
			cmd := exec.Command("networker", "request", "-m", "post", "--files", "formname=doesntexist.png", "https://google.com")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "no such file or directory")
		})
		test.WithNetworker(t, "add headers with uneven number of key/value pairs", func(t *testing.T) {
			cmd := exec.Command("networker", "request", "-H", "key:value,key", "https://google.com")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "uneven number of key/value pairs")
		})
	})
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithNetworker(t, "CRUD", func(t *testing.T) {
			test.WithMockServer(t, func(t *testing.T, testserverURL string) {
				// create a new object with a json string literal
				cmd := exec.Command("networker", "request",
					"-H", "Authorization: Bearer doesntmatter",
					"-m", "post",
					"-b", `{"field": "doesntmatter"}`,
					"-s",
					testserverURL,
				)
				output, err := cmd.CombinedOutput()
				require.NoError(t, err)
				object := new(test.MockObject)
				require.NoError(t, json.Unmarshal(output, &object))
				require.Equal(t, "doesntmatter", object.Field)
				require.NotZero(t, object.ID)

				// delete the object
				cmd = exec.Command("networker", "request",
					"-H", "Authorization: Bearer doesntmatter",
					"-m", "delete",
					fmt.Sprintf("%s?id=%d", testserverURL, object.ID),
				)
				output, err = cmd.CombinedOutput()
				require.NoError(t, err)
				require.Contains(t, string(output), "status_code: 200")
				projectRoot := test.ProjectRoot(t)

				// create another but this time using a json file
				cmd = exec.Command("networker", "request",
					"-H", "Authorization: Bearer doesntmatter",
					"-m", "post",
					"-b", path.Join(projectRoot, "internal/test/body.json"),
					"-s",
					testserverURL,
				)
				output, err = cmd.CombinedOutput()
				require.NoError(t, err)
				object = new(test.MockObject)
				require.NoError(t, json.Unmarshal(output, object))
				require.Equal(t, "name", object.Field)
				require.NotZero(t, object.ID)

				// get the object
				cmd = exec.Command("networker", "request",
					"-H", "Authorization: Bearer doesntmatter",
					"-s",
					fmt.Sprintf("%s?id=%d", testserverURL, object.ID),
				)
				output, err = cmd.CombinedOutput()
				require.NoError(t, err)
				retrievedObject := new(test.MockObject)
				require.NoError(t, json.Unmarshal(output, retrievedObject))
				require.Equal(t, object, retrievedObject)

				// delete the object
				cmd = exec.Command("networker", "request",
					"-H", "Authorization: Bearer doesntmatter",
					"-m", "delete",
					fmt.Sprintf("%s?id=%d", testserverURL, retrievedObject.ID),
				)
				output, err = cmd.CombinedOutput()
				require.NoError(t, err)
				require.Contains(t, string(output), "status_code: 200")

				root := test.ProjectRoot(t)
				file1 := path.Join(root, "internal/test/cat_1.jpeg")
				file2 := path.Join(root, "internal/test/cat_2.jpeg")
				multiPartFormUploadArg := fmt.Sprintf("files=%s,%s", file1, file2)

				// upload a file
				cmd = exec.Command("networker", "request",
					"-H", "Authorization: Bearer doesntmatter",
					"-m", "put",
					"--files", multiPartFormUploadArg,
					testserverURL,
				)
				output, err = cmd.CombinedOutput()
				require.NoError(t, err)
				require.Contains(t, string(output), "status_code: 201")
			})
		})
	})
}
