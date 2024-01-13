package models4sportus

// Model defines a struct that holds info about a specific model of some brand
type Model struct {
	Brand string   `json:"brand" firestore:"brand"`
	Title string   `json:"title" firestore:"title"`
	Kinds []string `json:"kinds" firestore:"kinds"`
	YearRange
}

// Validate validates record
func (v Model) Validate() error {
	return nil
}

// Product DTO
type Product struct {
	Kinds     []string // e.g. ["kiting", "kiteboard", "twintip"]
	Brand     string   // e.g. "Duotone"
	ModelName string   // e.g. "Neo SLS"
	ModelYear int      // e.g. 2021
}

// Quiver DTO
type Quiver struct {
	UserID string
}

// QuiverItem DTO
type QuiverItem struct {
	Item
	Product
	Location    []string `json:"location,omitempty" firestore:"location,omitempty"` // e.g. ["IE", "ie/dublin"]
	PinHoles    int      `json:"pinHoles,omitempty" firestore:"pinHoles,omitempty"`
	Repairs     int      `json:"repairs,omitempty" firestore:"repairs,omitempty"`
	Currency    string   `json:"currency,omitempty" firestore:"currency,omitempty"`
	Condition   int      `json:"condition,omitempty" firestore:"condition,omitempty"`
	PriceNew    int      `json:"priceNew,omitempty" firestore:"priceNew,omitempty"`
	PriceBought int      `json:"priceBought,omitempty" firestore:"priceBought,omitempty"`
	PriceSale   int      `json:"priceSale,omitempty" firestore:"priceSale,omitempty"`
	PrivateNote string   `json:"privateNote,omitempty" firestore:"privateNote,omitempty"`
	SaleDesc    string   `json:"saleDesc,omitempty" firestore:"saleDesc,omitempty"`
}

// QuiverWantedCollection defines collection name
const QuiverWantedCollection = "QuiverWanted"

// Wanted DTO
type Wanted struct {
	Item
	Locations   []string
	Brands      []string `json:"brands,omitempty" firestore:"brands,omitempty"`
	Models      []string `json:"models,omitempty" firestore:"models,omitempty"`
	Description string   `json:"description,omitempty" firestore:"description,omitempty"`
	YearRange
	PriceRange
	SizeRange
	LengthRange
	WidthRange
	WeightRange
	RepairsRange
	PinholesRange
}
