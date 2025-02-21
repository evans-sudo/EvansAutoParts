package model

type CartItem struct {
	Part     *Part
	Quantity int
}

type Cart struct {
	Items []CartItem
}
