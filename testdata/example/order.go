package example

type Order struct {
	CustomerName string `json:"customer_name"`
}

type Orders []Order
