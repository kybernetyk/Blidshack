package main

import (
	"exp/draw"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"fmt"
	"math"
)

var (
	imagefile = "aktkartepolengrau.jpg"
)


type Location struct {
	Lat struct {
		deg int
		hrs int
		min int
	}
	Lon struct {
		deg int
		hrs int
		min int
	}
}

type Bounds struct {
	lat_min, lat_max int
	lon_min, lon_max int

	px_w, px_h int

	filename string
}

var poland = Bounds{
	lat_min: 493500,
	lat_max: 550500,

	lon_min: 133934,
	lon_max: 242754,

	px_w: 500,
	px_h: 460,

	filename: "poland.jpg",
}

var germany = Bounds{
	lat_min: 464500,
	lat_max: 544500,

	lon_min: 46230,
	lon_max: 154030,

	px_w: 480,
	px_h: 580,

	filename: "germany.jpg",
}
var land = &poland

func (bnds Bounds) GeoLocation(idx int) Location {

	/*
		lat_min := 136019
		lat_max := 235876
		ilat := ((lat_max-lat_min)/500)*x + lat_min

		lon_min := 488029
		lon_max := 548022
		ilon := (((lon_min-lon_max)/460)*-y - lon_max) * -1
	*/
	x := idx % bnds.px_w
	y := idx / bnds.px_h

	//	ilat := ((bnds.lat_max-bnds.lat_min)/bnds.px_w)*x + bnds.lat_min
	//	ilon := (((bnds.lon_min-bnds.lon_max)/bnds.px_h)*-y - bnds.lon_max) * -1
	ilon := ((bnds.lon_max-bnds.lon_min)/bnds.px_w)*x + bnds.lon_min
	ilat := (((bnds.lat_min-bnds.lat_max)/bnds.px_h)*-y - bnds.lat_max) * -1

	ret := Location{}
	ret.Lat.deg = ilat / 10000
	ret.Lat.hrs = (ilat / 100) % 100
	ret.Lat.min = ilat % 100

	/*	if ret.Lat.min > 60 {
			rest := ret.Lat.min % 60
			ret.Lat.hrs ++
			ret.Lat.min = rest
		}

		if ret.Lat.hrs > 60 {
			rest := ret.Lat.hrs % 60
			ret.Lat.deg ++
			ret.Lat.hrs = rest
		}
	*/
	ret.Lon.deg = ilon / 10000
	ret.Lon.hrs = (ilon / 100) % 100
	ret.Lon.min = ilon % 100

	/*if ret.Lon.min > 60 {
		rest := ret.Lon.min % 60
		ret.Lon.hrs ++
		ret.Lon.min = rest
	}

	if ret.Lon.hrs > 60 {
		rest := ret.Lon.hrs % 60
		ret.Lon.deg ++
		ret.Lon.hrs = rest
	}
	*/

	return ret

}


// rgba returns an RGBA version of the image, making a copy only if
// necessary.
func rgba(m image.Image) *image.RGBA {
	if r, ok := m.(*image.RGBA); ok {
		return r
	}
	b := m.Bounds()
	r := image.NewRGBA(b.Dx(), b.Dy())
	draw.Draw(r, b, m, image.ZP)
	return r
}

func stripColors(img image.Image) *image.RGBA {
	rgb := rgba(img)

	fmt.Printf("conv: %#v\n", rgb.Pix[0])

	pix := rgb.Pix[0]
	var avg float64
	min_dist := 30.0
	for i := 0; i < len(rgb.Pix); i++ {
		pix = rgb.Pix[i]

		//first we need to kill everything that is some shade of grey
		avg = (float64(pix.R) + float64(pix.G) + float64(pix.B)) / 3.0

		dist_r := math.Fabs(float64(pix.R) - avg)
		dist_g := math.Fabs(float64(pix.G) - avg)
		dist_b := math.Fabs(float64(pix.B) - avg)

		if dist_r <= min_dist && dist_g <= min_dist && dist_b <= min_dist {
			rgb.Pix[i].R = 0
			rgb.Pix[i].G = 0
			rgb.Pix[i].B = 0
		}

		//now keep only the red data (which is the current one)
		if (pix.R > pix.G && pix.R > pix.B) &&
			(pix.R-pix.G) > 140 && (pix.R-pix.B) > 140 {
			//keep pixel
		} else {
			rgb.Pix[i].R = 0
			rgb.Pix[i].G = 0
			rgb.Pix[i].B = 0
		}
	}

	return rgb
}

func extractData(img image.Image) {
	rgb := rgba(img)

	data := make([]Location, 0, 0)
	for i := 0; i < len(rgb.Pix); i++ {
		if rgb.Pix[i].R != 0 && rgb.Pix[i].G != 0 && rgb.Pix[i].B != 0 {
			//			i := x + y * w

			data = append(data, land.GeoLocation(i))
		}
	}

	//	fmt.Printf("%#v", locForXY(350,210))

	//http://maps.google.com/?ie=UTF8&ll=37.0625,-95.677068&spn=33.710275,65.654297&z=4
	//	loc := locForXY(350,210);
	//loc := data[0]

	land := &poland
	idx := 396 + 371*land.px_w
	loc := land.GeoLocation(idx)
	fmt.Printf("rzeszow: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
		loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
		loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)

	idx = 245 + 50*land.px_w
	loc = land.GeoLocation(idx)
	fmt.Printf("gdansk: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
		loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
		loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)

	idx = 350 + 211*land.px_w
	loc = land.GeoLocation(idx)
	fmt.Printf("warszawa: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
		loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
		loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)

	idx = 100 + 233*land.px_w
	loc = land.GeoLocation(idx)
	fmt.Printf("zelona gora: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
		loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
		loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)

	land = &germany

	idx = 235 + 59*land.px_w
	loc = land.GeoLocation(idx)
	fmt.Printf("kiel: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
		loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
		loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)

	idx = 227 + 121 * land.px_w
	loc = land.GeoLocation(idx)
	fmt.Printf("hamburg: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
		loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
		loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)

	idx = 384 + 198*land.px_w
	loc = land.GeoLocation(idx)
	fmt.Printf("berlin: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
		loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
		loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)

	idx = 81 + 294*land.px_w
	loc = land.GeoLocation(idx)
	fmt.Printf("ddorf: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
		loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
		loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)

	idx = 300 + 510*land.px_w
	loc = land.GeoLocation(idx)
	fmt.Printf("muenchen: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
		loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
		loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)



	//	fmt.Printf("%#v\n", data)
	//	fmt.Printf("%d\n", len(data))
}

func main() {
	file, err := os.Open(land.filename)
	if err != nil {
		panic(err.String())
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err.String())
	}

	fmt.Printf("fmt: %#v\n", img.ColorModel())
	fmt.Printf("loaded: %#v\n", img.At(0, 0))

	rgb := stripColors(img)
	extractData(rgb)

	f, err := os.Create("lol.png")
	if err != nil {
		panic(err.String())
	}
	defer f.Close()

	err = png.Encode(f, rgb)
	if err != nil {
		panic(err.String())
	}
	fmt.Printf("success!\n")
}
