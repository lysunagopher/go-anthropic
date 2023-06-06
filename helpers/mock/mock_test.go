package mock_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/staropshq/go-anthropic/helpers/mock"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	art := assert.New(t)
	client := mock.NewHTTPClient()
	client.RespondWith([]byte("Hello, World!"), http.StatusOK, nil)
	response, err := client.Do(nil)
	defer response.Body.Close()
	if art.NoError(err) {
		art.Equal(response.StatusCode, http.StatusOK)
		body, err := ioutil.ReadAll(response.Body)
		if art.NoError(err) {
			art.Equal([]byte("Hello, World!"), body)
		}
	}
}
