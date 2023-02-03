package example

type Order struct {
	CustomerName string `json:"customer_name"`
}

type OrderSlice []Order

type Orders struct {
	Orders OrderSlice `json:"orders"`
}
