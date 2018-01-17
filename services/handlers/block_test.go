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
