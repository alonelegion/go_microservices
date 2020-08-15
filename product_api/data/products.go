package data

import (
	"fmt"
)

// ErrProductNotFound in an error raised when a product
// can not be found in the database
// ErrProductNotFound вызывается при возникновении ошибки,
// когда товар не может быть найден
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for this user
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this product
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this product
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`
}

// Products is a collection of Product
// Products это коллекция из Product
type Products []*Product

// GetProducts returns all products from the database
// GetProducts возвращает все товары из хранилища данных
func GetProducts() Products {
	return productList
}

// GetProductByID returns a single product which matches the id from the
// database
// If a product is not found this function returns a ProductNotFound error
// GetProductByID возвращает один товар, который соответствует идентификатору
// из база данных
// Если товар не был найден, эта функция возвращает ошибку ProductNotFound
func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

// UpdateProduct replaces a product in the database with the given id
// If a product with the given id does not exist in the database
// this function returns a ErrProductNotFound error
// UpdateProduct заменяет продукт в базе данных по заданному id
// Если продукта с указанным идентификатором нет в базе
// эта функция возвращает ошибку ErrProductNotFound
func UpdateProduct(p Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	productList[i] = &p

	return nil
}

// AddProduct adds a new product to the database
// AddProducts добавляет новый товар в базу данных
func AddProduct(p Product) {
	// get the next id in sequence
	// получить следующий идентификатор в последовательности
	maxID := productList[len(productList)-1].ID
	p.ID = maxID + 1
	productList = append(productList, &p)
}

// DeleteProduct deletes a product from the database
// DeleteProduct удаляет товар из базы данных
func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// return -1 when no product can be found
// findIndex находит индекс товара в базе данных
// возвращает -1 когда нет нужного товара
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

// productList is a hard coded list of products for this
// example data source
// productList это захардкоженный список товаров для
// примера источника данных
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
	&Product{
		ID:          3,
		Name:        "Cappuccino",
		Description: "Is an espresso-based coffee drink that originated in Italy",
		Price:       3.15,
		SKU:         "ljbgj",
	},
}
