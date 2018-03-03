package handlers

import (
	"fmt"
	"net/http"
)

func (t *mainSuite) TestQueryAccount() {
	urlPath := fmt.Sprintf("/account/%s", Address)
	resp, err := http.Get(server.URL + urlPath)
	t.Must(t.Nil(err))

	t.Must(t.Equal(resp.StatusCode, http.StatusOK))
}
