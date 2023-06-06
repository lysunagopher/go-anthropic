package mock

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	art := assert.New(t)
	client := NewHTTPClient()
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
