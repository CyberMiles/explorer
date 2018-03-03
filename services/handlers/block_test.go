package handlers

import (
	"fmt"
	"net/http"
)

func (t *mainSuite) TestQueryBlock() {
	urlPath := fmt.Sprintf("/block/%d", Height)
	resp, err := http.Get(server.URL + urlPath)
	t.Must(t.Nil(err))

	t.Must(t.Equal(resp.StatusCode, http.StatusOK))
}

func (t *mainSuite) TestQueryValidators() {
	urlPath := fmt.Sprintf("/validators/%d", Height)
	resp, err := http.Get(server.URL + urlPath)
	t.Must(t.Nil(err))

	t.Must(t.Equal(resp.StatusCode, http.StatusOK))
}

func (t *mainSuite) TestQueryRecentBlocks() {
	urlPath := fmt.Sprintf("/blocks/recent")
	resp, err := http.Get(server.URL + urlPath)
	t.Must(t.Nil(err))

	t.Must(t.Equal(resp.StatusCode, http.StatusOK))
}
