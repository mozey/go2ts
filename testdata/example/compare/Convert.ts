export interface Order {
  customer_name: string;
}

export interface Orders extends Array<Order>{}

export interface Product {
  title: string;
}

export interface Products extends Array<Product>{}