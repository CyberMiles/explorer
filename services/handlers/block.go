package handlers

import (
  "net/http"

  "github.com/gorilla/mux"
  "github.com/spf13/cast"

  "github.com/tendermint/tmlibs/common"
)

// queryBlock is to query a block by height
func queryBlock(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  height := args["height"]

  c, err := getSecureNode()
  if err != nil {
    common.WriteError(w, err)    
    return
  }

  h := cast.ToInt64(height)
  block, err := c.Block(&h)
  if err != nil {
    common.WriteError(w, err)    
    return
  }
  if err := printResult(w, block); err != nil {
    common.WriteError(w, err)    
  }
}

// mux.Router registrars

func RegisterQueryBlock(r *mux.Router) error {
  r.HandleFunc("/block/{height}", queryBlock).Methods("GET")
  return nil
}

func RegisterBlock(r *mux.Router) error {
  funcs := []func(*mux.Router) error{
    RegisterQueryBlock,
  }

  for _, fn := range funcs {
    if err := fn(r); err != nil {
      return err
    }
  }
  return nil
}

// End of mux.Router registrars
