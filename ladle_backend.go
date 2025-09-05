package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

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

func main() {
	//Enabling CORS for all routes
	http.HandleFunc("/api/convert",
		corsHandler(convertColorHandler))
	http.HandleFunc("/api/palette",
		corsHandler(generatePaletteHandler))
	http.HandleFunc("/api/random",
		corsHandler(randomColorHandler))

	fmt.Println("Ladle API server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func corsHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Autorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Conversion of colors between different formats
func convertColorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ColorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var color Color
	var err error

	switch req.Type {
	case "hex":
		color, err = parseHexColor(req.Value.(string))
	case "rgb":
		color, err = parseRGBColor(req.Value)
	case "hsl":
		color, err = parseHSLColor(req.Value)
	case "hsv":
		color, err = parseHSVColor(req.Value)
	default:
		http.Error(w, "Unsupported color type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(color)
}

// Generating the color palette
func generatePaletteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PaletteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	palette := generatePalette(req.BaseColor, req.Type, req.Count)

	w.Header.Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]Color{
		"palette": palette,
	})
}

// Generating a random color
func randomColorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Generating random HSV values
	h := int(math.Round(math.Mod(float64(rand.Int()), 360)))

	//50-100%
	s := 50 + int(math.Round(math.Mod(float64(rand.Int()), 50)))

	//50-100%
	v := 50 + int(math.Round(math.Mod(float64(rand.Int()), 50)))

	color := Color{
		HSV: ColorHSV{H: h, S: s, V: v},
	}

	//Conversion to other formats
	color.RGB = hsvToRgb(h, s, v)
	color.HSL = rgbToHsl(color.RGB.R, color.RGB.G, color.RGB.B)
	color.Hex = rgbToHex(color.RGB.R, color.RGB.G, color.RGB.B)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(color)
}

//Color conversion functions & parsing
//   |
//  |
// V
