package handlers

import (
  "fmt"
  "strings"
  "net/http"
)

func (t *mainSuite) TestQueryTx() {
  urlPath := fmt.Sprintf("/tx/%s", TxHash)
  resp, err := http.Get(server.URL + urlPath)
  t.Must(t.Nil(err))

  t.Must(t.Equal(resp.StatusCode, http.StatusOK))  
}

func (t *mainSuite) TestQueryRawTx() {
  urlPath := fmt.Sprintf("/tx/%s/raw", TxHash)
  resp, err := http.Get(server.URL + urlPath)
  t.Must(t.Nil(err))

  t.Must(t.Equal(resp.StatusCode, http.StatusOK))  
}

func (t *mainSuite) TestSearchTxByBlock() {
  urlPath := fmt.Sprintf("/block/%d/tx", Height)
  resp, err := http.Get(server.URL + urlPath)
  t.Must(t.Nil(err))

  t.Must(t.Equal(resp.StatusCode, http.StatusOK))  
}

func (t *mainSuite) TestSearchCoinTxByAccount() {
  urlPath := fmt.Sprintf("/account/%s/tx/coin", Address)
  resp, err := http.Get(server.URL + urlPath)
  t.Must(t.Nil(err))

  t.Must(t.Equal(resp.StatusCode, http.StatusOK))  
}

func (t *mainSuite) TestDecodeRaw() {
  urlPath := fmt.Sprintf("/tx/decode")
  payload := strings.NewReader(RawTx)
  resp, err := http.Post(server.URL + urlPath, "", payload)
  t.Must(t.Nil(err))

  t.Must(t.Equal(resp.StatusCode, http.StatusOK))  
}
