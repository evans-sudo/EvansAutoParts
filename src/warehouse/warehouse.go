package main

import (
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"net/http"
)

var inventory map[string]int

func init() {
	inventory = map[string]int{
		"8472983": 42,
		"8479342": 27,
	}
}

func main() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(InventoryService), "Inventory")
	http.Handle("/rpc", s)

	http.ListenAndServe(":1234", nil)
}

type CheckInventoryArgs struct {
	PartNumber string
}

type InventoryService struct{}

func (is *InventoryService) CheckInventory(r *http.Request, args *CheckInventoryArgs,
	reply *int) error {

	*reply = inventory[args.PartNumber]
	println(reply)
	return nil

}
