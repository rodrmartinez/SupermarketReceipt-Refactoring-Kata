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

func checkProductsInBundle(cart []ProductQuantity, bundles []ProductQuantity) bool {
	var itemsInBoundle []ProductQuantity
	for _, item := range cart {
		for _, bundle := range bundles {
			if item.product == bundle.product && item.quantity >= bundle.quantity {
				itemsInBoundle = append(itemsInBoundle, item)
			}
		}
	}
	if len(itemsInBoundle) == len(bundles) {
		return true
	}

	return false
}

func TestDiscounts(t *testing.T) {
	var catalog = NewFakeCatalog()

	var teller = NewTeller(catalog)

	var Tests = []struct {
		products []ProductQuantity
		offers   []SpecialOffer
		bundles  []Bundle
		expected float64
	}{
		{
			//No disscount
			products: []ProductQuantity{{toothbrush, 1}},
			offers:   []SpecialOffer{},
			expected: 0.99,
		},
		{
			// 2x1 deal
			products: []ProductQuantity{{toothbrush, 2}},
			offers:   []SpecialOffer{{TwoForAmount, toothbrush, catalog._prices["toothbrush"]}},
			expected: 0.99,
		},
		{
			//10% discount
			products: []ProductQuantity{{apples, 2.5}},
			offers:   []SpecialOffer{{TenPercentDiscount, apples, 20.0}},
			expected: 3.98,
		},
		{
			//20% discount
			products: []ProductQuantity{{rice, 2}},
			offers:   []SpecialOffer{{TenPercentDiscount, rice, 10.0}},
			expected: 4.48,
		},
		{
			// 5xAmount
			products: []ProductQuantity{{toothpaste, 5}},
			offers:   []SpecialOffer{{FiveForAmount, toothpaste, 7.49}},
			expected: 7.49,
		},
		{
			//3x2 deal
			products: []ProductQuantity{{cherrytomatoes, 3}},
			offers:   []SpecialOffer{{ThreeForTwo, cherrytomatoes, 0.69}},
			expected: 1.38,
		},
		{
			//Multiple disocunts
			products: []ProductQuantity{{cherrytomatoes, 3}, {toothbrush, 2}},
			offers:   []SpecialOffer{{ThreeForTwo, cherrytomatoes, 0.69}, {TwoForAmount, toothbrush, catalog._prices["toothbrush"]}},
			expected: 2.37,
		},
		{
			// 5xAmount with 6 products
			products: []ProductQuantity{{toothpaste, 6}},
			offers:   []SpecialOffer{{FiveForAmount, toothpaste, 7.49}},
			expected: 9.28,
		},
		{
			// products quantity equals bundle
			products: []ProductQuantity{{toothpaste, 1}, {toothbrush, 1}},
			bundles:  []Bundle{{BundleDiscount, []ProductQuantity{{toothpaste, 1}, {toothbrush, 1}}, 10}},
			expected: 2.5,
		},
		{
			// products quantity equals bundle 2 and 3
			products: []ProductQuantity{{toothpaste, 2}, {toothbrush, 3}},
			bundles:  []Bundle{{BundleDiscount, []ProductQuantity{{toothpaste, 2}, {toothbrush, 3}}, 10}},
			expected: 5.9,
		},
	}
	for _, test := range Tests {
		var cart = NewShoppingCart()
		for _, product := range test.products {
			cart.addItemQuantity(product.product, product.quantity)
		}
		for _, offer := range test.offers {
			teller.addSpecialOffer(offer.offerType, offer.product, offer.argument)
		}
		for _, item := range test.products {
			for _, bundle := range test.bundles {
				if checkProductsInBundle(test.products, bundle.products) {
					teller.addSpecialOffer(bundle.offerType, item.product, bundle.argument)
				}
			}
		}
		var receipt = teller.checksOutArticlesFrom(cart)
		assert.Equal(t, test.expected, receipt.totalPrice())
		require.Equal(t, len(test.products), len(receipt.items))
	}

}
