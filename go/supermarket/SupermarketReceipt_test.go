package supermarket

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var toothbrush = Product{name: "toothbrush", unit: Each}
var apples = Product{name: "apples", unit: Kilo}
var rice = Product{name: "rice", unit: Each}
var toothpaste = Product{name: "toothpaste", unit: Each}
var cherrytomatoes = Product{name: "cherrytomatoes", unit: Each}

type FakeCatalog struct {
	_products map[string]Product
	_prices   map[string]float64
}

func (c FakeCatalog) unitPrice(product Product) float64 {
	return c._prices[product.name]
}

func (c FakeCatalog) addProduct(product Product, price float64) {
	c._products[product.name] = product
	c._prices[product.name] = price
}

func NewFakeCatalog() *FakeCatalog {
	var c FakeCatalog
	c._products = make(map[string]Product)
	c._prices = make(map[string]float64)

	c.addProduct(toothbrush, 0.99)
	c.addProduct(apples, 1.99)
	c.addProduct(rice, 2.49)
	c.addProduct(toothpaste, 1.79)
	c.addProduct(cherrytomatoes, 0.69)

	return &c
}

type Test struct {
	product  Product
	offer    SpecialOffer
	quantity float64
}

func TestNoDiscounts(t *testing.T) {

	var catalog = NewFakeCatalog()

	var teller = NewTeller(catalog)

	var cart = NewShoppingCart()

	test1 := Test{toothbrush, SpecialOffer{}, 1}

	cart.addItemQuantity(test1.product, test1.quantity)

	// ACT
	var receipt = teller.checksOutArticlesFrom(cart)

	// ASSERT
	assert.Equal(t, 0.99, receipt.totalPrice())
	assert.Equal(t, 0, len(receipt.discounts))
	require.Equal(t, 1, len(receipt.items))
}

func TestTwoForAmount(t *testing.T) {

	var catalog = NewFakeCatalog()

	var teller = NewTeller(catalog)

	var cart = NewShoppingCart()

	test1 := Test{
		toothbrush,
		SpecialOffer{
			TwoForAmount,
			toothbrush,
			catalog._prices["toothbrush"]},
		2,
	}

	teller.addSpecialOffer(test1.offer.offerType, test1.offer.product, test1.offer.argument)
	cart.addItemQuantity(test1.product, test1.quantity)
	// ACT
	var receipt = teller.checksOutArticlesFrom(cart)

	// ASSERT
	assert.Equal(t, 0.99, receipt.totalPrice())
	assert.Equal(t, 1, len(receipt.discounts))
	require.Equal(t, 1, len(receipt.items))
}

func TestTwentyPercentDiscount(t *testing.T) {

	var catalog = NewFakeCatalog()

	var teller = NewTeller(catalog)
	teller.addSpecialOffer(TenPercentDiscount, apples, 20.0)

	var cart = NewShoppingCart()
	cart.addItemQuantity(apples, 2.5)

	// ACT
	var receipt = teller.checksOutArticlesFrom(cart)

	// ASSERT
	assert.Equal(t, 3.98, receipt.totalPrice())
	assert.Equal(t, 1, len(receipt.discounts))
	require.Equal(t, 1, len(receipt.items))
}

func TestTenPercentDiscount(t *testing.T) {

	var catalog = NewFakeCatalog()

	var teller = NewTeller(catalog)

	var cart = NewShoppingCart()
	cart.addItemQuantity(rice, 2)
	teller.addSpecialOffer(TenPercentDiscount, rice, 10.0)

	// ACT
	var receipt = teller.checksOutArticlesFrom(cart)

	// ASSERT
	assert.Equal(t, 4.48, receipt.totalPrice())
	assert.Equal(t, 1, len(receipt.discounts))
	require.Equal(t, 1, len(receipt.items))
}

func TestFiveForAmount(t *testing.T) {

	var catalog = NewFakeCatalog()

	var teller = NewTeller(catalog)

	var cart = NewShoppingCart()
	cart.addItemQuantity(toothpaste, 5)
	teller.addSpecialOffer(FiveForAmount, toothpaste, 7.49)

	// ACT
	var receipt = teller.checksOutArticlesFrom(cart)

	// ASSERT
	assert.Equal(t, 7.49, receipt.totalPrice())
	assert.Equal(t, 1, len(receipt.discounts))
	require.Equal(t, 1, len(receipt.items))
}

func TestThreeForTwo(t *testing.T) {

	var catalog = NewFakeCatalog()

	var teller = NewTeller(catalog)

	var cart = NewShoppingCart()
	cart.addItemQuantity(cherrytomatoes, 3)
	teller.addSpecialOffer(ThreeForTwo, cherrytomatoes, 0.69)

	// ACT
	var receipt = teller.checksOutArticlesFrom(cart)

	// ASSERT
	assert.Equal(t, 1.38, receipt.totalPrice())
	assert.Equal(t, 1, len(receipt.discounts))
	require.Equal(t, 1, len(receipt.items))
}
