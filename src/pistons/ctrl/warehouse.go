package ctrl

import (
	"bytes"
	encJson "encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/json"
	"net/http"
)

type warehouseController struct{}

func (wc *warehouseController) CheckInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	partNumber := vars["partNumber"]

	args := checkInventoryArgs{partNumber}

	msg, _ := json.EncodeClientRequest("Inventory.CheckInventory", args)
	resp, _ := http.Post("http://localhost:1234/rpc", "application/json", bytes.NewReader(msg))

	var result int
	json.DecodeClientResponse(resp.Body, &result)
	ir := inventoryResult{result}
	data, _ := encJson.Marshal(ir)

	w.Header().Add("Content-Type", "applicaton/json")
	w.Write(data)
}

type checkInventoryArgs struct {
	PartNumber string
}

type inventoryResult struct {
	Inventory int `json:"inventory"`
}