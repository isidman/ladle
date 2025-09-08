package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
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
func parseHexColor(hex string) (Color, error) {
	if len(hex) == 7 && hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return Color{}, fmt.Errorf("invalid hex color format")
	}

	r, err1 := strconv.ParseInt(hex[0:2], 16, 0)
	g, err2 := strconv.ParseInt(hex[2:4], 16, 0)
	b, err3 := strconv.ParseInt(hex[4:6], 16, 0)

	if err1 != nil || err2 != nil || err3 != nil {
		return Color{}, fmt.Errorf("invalid hex color format")
	}

	rgb := ColorRGB{R: int(r), G: int(g), B: int(b)}

	return Color{
		Hex: "#" + hex,
		RGB: rgb,
		HSL: rgbToHsl(int(r), int(g), int(b)),
		HSV: rgbToHsv(int(r), int(g), int(b)),
	},nil
}

func parseRGBColor(value interface{}) (Color, error) {
	rgb, ok := value.(map[string]interface{})
	if !ok {
		return Color{}, fmt.Errorf("invalid RGB format")
	}

	r :=int(rgb["r"].(float64))
	g :=int(rgb["g"].(float64))
	b :=int(rgb["b"].(float64))

	return Color{
		RGB: ColorRGB{R: r, G: g, B: b},
		Hex: rgbToHex(r, g, b),
		HSL: rgbToHsl(r, g, b),
		HSV: rgbToHsv(r, g, b),
	},nil
}

func parseHSLColor(value interface{}) (Color, error) {
	hsl, ok := value.(map[string]interface{})
	if !ok {
		return Color{}, fmt.Errorf("invalid HSL format")
	}

	h := int(hsl["h"].(float64))
	s := int(hsl["s"].(float64))
	l := int(hsl["l"].(float64))

	rgb := hslToRgb(h, s, l)

	return Color{
		HSL: ColorHSL{H: h, S: s, L: l},
		RGB: rgb,
		Hex: rgbToHex(rgb.R, rgb.G, rgb.B),
		HSV: rgbToHsv(rgb.R, rgb.G, rgb.B),
	}, nil
}

func parseHSVColor(value interface{}) (Color, error) {
	hsv, ok := value.(map[string]interface{})
	if !ok {
		return Color{}, fmt.Errorf("invalid HSV format")
	}

	h := int(hsv["h"].(float64))
	s := int(hsv["s"].(float64))
	v := int(hsv["v"].(float64))

	rgb := hsvToRgb(h, s, v)

	return Color{
		HSV: ColorHSV{H: h, S: s, V: v},
		RGB: rgb,
		Hex: rgbToHex(rgb.R, rgb.G, rgb.B),
		HSL: rgbToHsl(rgb.R, rgb.G, rgb.B),
	}, nil
}

//Color conversion algorithmics

func hsvToRgb(h, s, v int) ColorRGB {
	hf := float64(h) / 360.0
	sf := float64(s) / 100.0
	vf := float64(v) / 100.0

	c := vf * sf
	x := c * (1- math.Abs(math.Mod(hf*6, 2) - 1))
	m := vf - c

	var r, g, b float64

	switch int(hf * 6) {
	case 0:
		r, g, b = c, x, 0
	case 1:
		r, g, b = x, c, 0
	case 2:
		r, g, b = 0, c, x
	case 3: 
		r, g, b = 0, x, c
	case 4:
		r, g, b = x, 0, c
	case 5:
		r, g, b = c, 0, x
	}

	return ColorRGB{
		R: int(math.Round((r + m) * 255)),
		G: int(math.Round((g + m) * 255)),
		B: int(math.Round((b + m) * 255)),
	}
}

func rgbToHsv(r, g, b int) ColorHSV {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	var h, s, v float64

	//Hue calculations
	if delta == 0 {
		h = 0
	} else if max == rf {
		h = 60 * math.Mod((gf-bf)/delta, 6)
	} else if max == gf {
		h = 60 * ((bf-rf)/delta + 2)
	} else {
		h = 60 * ((rf-gf)/delta + 4)
	}

	//Saturation calculations
	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}

	//Value calcuations
	v = max

	return ColorHSV{
		H: int(math.Round(h)),
		S: int(math.Round(s * 100)),
		V: int(math.Round(v * 100)),
	}
}

func hslToRgb(h, s, l int) ColorRGB {
	hf := float64(h) / 360.0
	sf := float64(s) / 100.0
	lf := float64(l) / 100.0

	c := (1 - math.Abs(2*lf - 1)) * sf
	x := c * (1 - math.Abs(math.Mod(hf*6, 2) - 1))
	m := lf - c/2

	var r, g, b float64

	switch int(hf * 6) {
	case 0:
		r, g, b = c, x, 0
	case 1:
		r, g, b = x, c, 0
	case 2:
		r, g, b = 0, c, x
	case 3:
		r, g, b = 0, x, c
	case 4:
		r, g, b = x, 0, c
	case 5:
		r, g, b = c, 0, x
	}

	return ColorRGB{
		R: int(math.Round((r + m) * 255)),
		G: int(math.Round((g + m) * 255)),
		B: int(math.Round((b + m) * 255)),
	}
}

func rgbToHsl(r, g, b int) ColorHSL {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	var h, s, l float64

	//Light calculation
	l = (max + min) / 2
	
	//Saturation calculation
	if delta == 0 {
		s = 0
	} else if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2 - max - min)
	}

	//Calculate hue
	if delta == 0 {
		h = 0
	} else if max == rf {
		h = 60 * math.Mod((gf-bf)/delta, 6)
	} else if max == gf {
		h = 60 * ((bf-rf)/delta + 2)
	} else {
		h = 60 * ((rf-gf)/delta + 4)
	}

	return ColorHSL{
		H: int(math.Round(h)),
		S: int(math.Round(s * 100)),
		L: int(math.Round(l * 100)),
	}
}

func rgbToHex(r, g, b int) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

//Palette genaration Functions
func generatePalette(baseColor Color, paletteType string, count int)[]Color {
	colors := make([]Color, count)

	switch paletteType {
	case "complementary":
		return generateComplementaryPalette(baseColor, count)
	case "analogous":
		return generateAnalogousPalette(baseColor, count)
	case "triadic":
		return generateTriadicPalette(baseColor, count)
	default:
		return generateAnalogousPalette(baseColor, count)	
	}

	return colors
}

//Complementary palette function
func generateComplementaryPalette(baseColor, count int) []Color {
	colors:= make([Color, count])

	for i := 0; i < count; i++ {
		hue := (baseColor.HSV.H + (180 * i / (count - 1))) % 360
		hsv := ColorHSV{
			H: hue,
			S: baseColor.HSV.S,
			V: baseColor.HSV.V,
		}

		rgb := hsvToRgb(hsv.H, hsv.S, hsv.V)
		colors[i] = Color{
			HSV: hsv,
			RGB: rgb,
			HSL: rgbToHsl(rgb.R, rgb.G, rgb.B),
			Hex: rgbToHex(rgb.R, rgb.G, rgb.B),
		}
	}

	return colors
}

//Analogous palette function
func generateAnalogousPalette(baseColor Color, count int) []Color {
	colors := make([]Color, count)
	
	for i := 0; i < count; i++ {
		hue := (baseColor.HSV.H + (60 * (i - count/2) / count)) % 360
		if hue < 0 {
			hue += 360
		}

		hsv := ColorHSV{
			H: hue,
			S: baseColor.HSV.S,
			V: baseColor.HSV.V,
		}

		rgb := hsvToRgb(hsv.H, hsv.S, hsv.V)
		colors[i] = Color{
			HSV: hsv,
			RGB: rgb,
			HSL: rgbToHsl(rgb.R, rgb.G, rgb.B),
			Hex: rgbToHex(rgb.R, rgb.G, rgb.B),
		}
	}

	return colors
}

//Triadic palette function
func generateTriadicPalette(baseColor Color, count int) []Color {
	colors := make([]Color, count)

	for i := 0; i < count; i++ {
		hue := (baseColor.HSV.H + (120 * i)) % 360

		hsv := ColorHSV{
			H: hue,
			S: baseColor.HSV.S,
			V: baseColor.HSV.V,
		}

		rgb := hsvToRgb(hsv.H, hsv.S, hsv.V)
		colors[i] = Color{
			HSV: hsv,
			RGB: rgb,
			HSL: rgbToHsl(rgb.R, rgb.G, rgb.B),
			Hex: rgbToHex(rgb.R, rgb.G, rgb.B),
		}
	}

	return colors
}

//Monochromatic palette function
func generateMonochromaticPalette(baseColor Color, count int) []Color {
	colors := make ([]Color, count)

	for i := 0; i < count; i++ {
		v := 20 + (80 * i / (count-1)) //Brightness varietion from 20% to 100%
		
		hsv := ColorHSV{
			H: baseColor.HSV.H,
			S: baseColor.HSV.S,
			V: v,
		}

		rgb := hsvToRgb(hsv.H, hsv.S, hsv.V)
		colors[i] = Color{
			HSV: hsv,
			RGB: rgb,
			HSL: rgbToHsl(rgb.R, rgb.G, rgb.B),
			Hex: rgbToHex(rgb.R, rgb.G, rgb.B),
		}
	}

	return colors
}