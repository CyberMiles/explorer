package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/tendermint/tmlibs/common"
)

// queryStatus is to query block chain status
func queryStatus(w http.ResponseWriter, r *http.Request) {
	c := commands.GetNode()
	status, err := c.Status()
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if err := printResult(w, status); err != nil {
		common.WriteError(w, err)
	}
}

// mux.Router registrars

func RegisterQueryStatus(r *mux.Router) error {
	r.HandleFunc("/status", queryStatus).Methods("GET")
	return nil
}

func RegisterStatus(r *mux.Router) error {
	funcs := []func(*mux.Router) error{
		RegisterQueryStatus,
	}

	for _, fn := range funcs {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

// End of mux.Router registrars
