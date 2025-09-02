package main

// The following struct represents a color in different formats
type Color struct {
	Hex string   `json:"hex"`
	RGB ColorRGB `json:"rgb"`
	HSL ColorHSL `json:"hsl"`
	HSV ColorHSV `json:"hsv"`
}

// Creating structs for all formats of Color

type ColorRGB struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type ColorHSL struct {
	H int `json:"h"`
	S int `json:"s"`
	L int `json:"l"`
}

type ColorHSV struct {
	H int `json:"h"`
	S int `json:"s"`
	V int `json:"v"`
}

// The ColorRequest struct is a part that represents the values for the incoming colors.
type ColorRequest struct {
	Type  string      `json:"type"` //"hex", "rgb", "hsl", "hsv"
	Value interface{} `json:"value"`
}

// The PaletteRequest is kind of a request for a color palette to generate
type PaletteRequest struct {
	BaseColor Color  `json:"baseColor"`
	Type      string `json:"type"` // "Complementary", "Analogous", "Triadic", "Monochrome"
}
