package example

type Product struct {
	Title string `json:"title"`
}

type ProductSlice []Product

type Products struct {
	Products ProductSlice `json:"products"`
}
