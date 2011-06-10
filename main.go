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

/* 
	did you know that google maps coordinates are strange?
*/

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
	lat_min, lat_max int //decimal encoded floats, 49.3524 -> 493524
	lon_min, lon_max int

	px_w, px_h int

	filename string
	img      image.Image
	data     []Location
}

var poland = Bounds{
	lat_min: 493500,
	lat_max: 550500,

	lon_min: 133934,
	lon_max: 242754,

	px_w: 500,
	px_h: 460,

	filename: "aktkartepolengrau.jpg",
}

var germany = Bounds{
	lat_min: 455000,
	lat_max: 552000,

	lon_min: 50756,
	lon_max: 155346,

	px_w: 480,
	px_h: 580,

	filename: "aktkartegergrau.jpg",
}

func (bnds *Bounds) GeoLocation(idx int) Location {
	x := idx % bnds.px_w
	y := idx / bnds.px_h

	ilon := ((bnds.lon_max-bnds.lon_min)/bnds.px_w)*x + bnds.lon_min
	ilat := (((bnds.lat_min-bnds.lat_max)/bnds.px_h)*-y - bnds.lat_max) * -1

	ret := Location{}
	ret.Lat.deg = ilat / 10000
	ret.Lat.hrs = (ilat / 100) % 100
	ret.Lat.min = ilat % 100

	ret.Lon.deg = ilon / 10000
	ret.Lon.hrs = (ilon / 100) % 100
	ret.Lon.min = ilon % 100

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

//strip everything that is not red
//this is probably more complicated than it needs to be - but it works
func stripColors(img image.Image) *image.RGBA {
	rgb := rgba(img)

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

//get "text" data from an image
func (land *Bounds) ExtractData() {
	rgb := rgba(land.img)

	land.data = make([]Location, 0, 0)
	for i := 0; i < len(rgb.Pix); i++ {
		if rgb.Pix[i].R != 0 && rgb.Pix[i].G != 0 && rgb.Pix[i].B != 0 {

			//loc := land.GeoLocation(i)
			land.data = append(land.data, land.GeoLocation(i))
			/*			fmt.Printf("data [%d]: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n", len(land.data),
						loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
						loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)
			*/
		}
	}

	//some tests

	/*	land = &poland
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

		/*land = &germany

		idx = 235 + 59*land.px_w
		loc = land.GeoLocation(idx)
		fmt.Printf("kiel: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d8&z=6\n",
			loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
			loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)

		idx = 227 + 121*land.px_w
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
		*/
	/*	idx := 180 + 380*land.px_w
		loc := land.GeoLocation(idx)
		fmt.Printf("ffm: http://maps.google.com/?ie=UTF8&ll=%.2d.%.2d%.2d,%.2d.%.2d%.2d&z=6\n",
			loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
			loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)
	
			land.data = append(land.data, loc)*/
}

func (land *Bounds) SaveData() {
	f, err := os.Create(land.filename + "_out.txt")
	if err != nil {
		panic(err.String())
	}
	defer f.Close()

	for _, loc := range land.data {
		s := fmt.Sprintf("%.2d.%.2d%.2d,%.2d.%.2d%.2d;",
			loc.Lat.deg, loc.Lat.hrs, loc.Lat.min,
			loc.Lon.deg, loc.Lon.hrs, loc.Lon.min)
		f.Write([]byte(s))
	}

}

func (land Bounds) DoJob(backchan chan bool) {
	file, err := os.Open(land.filename)
	if err != nil {
		panic(err.String())
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err.String())
	}

	rgb := stripColors(img)
	land.img = rgb
	land.ExtractData()
	land.SaveData()

	f, err := os.Create(land.filename + "_stripped.png")
	if err != nil {
		panic(err.String())
	}
	defer f.Close()

	err = png.Encode(f, land.img)
	if err != nil {
		panic(err.String())
	}

	backchan <- true
}

func main() {

	jobchan := make(chan bool)
	jobcount := 2

	go germany.DoJob(jobchan)
	go poland.DoJob(jobchan)

	for {
		select {
		case <-jobchan:
			jobcount--
		}

		if jobcount == 0 {
			break
		}
	}

}
