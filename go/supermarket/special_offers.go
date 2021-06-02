package supermarket

type SpecialOfferType int

const (
	TenPercentDiscount SpecialOfferType = iota
	ThreeForTwo
	TwoForAmount
	FiveForAmount
	BundleDiscount
)

type SpecialOffer struct {
	offerType SpecialOfferType
	product   Product
	argument  float64
}

type Discount struct {
	product        Product
	description    string
	discountAmount float64
}

type Bundle struct {
	offerType SpecialOfferType
	products  []ProductQuantity
	argument  float64
}
