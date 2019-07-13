package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"strconv"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	pt "github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/oak/shape"
)

var (
	colorMod color.RGBA
	imagePath string
)

func main() {

	oak.AddCommand("color", func(args []string) {
		if len(args) > 3 {
			r, g, b, a, err := parseRGBA(args)
			if err != nil {
				return
			}
			colorMod = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			src.Generator.(pt.Colorable).SetStartColor(startColor, startColorRand)
		} else {
			fmt.Println("Not enough args for color")
		}
	})

	oak.AddCommand("load", func(args []string) {
		if len(args) > 0 {
			
		}
	})


	oak.Add("demo", func(string, interface{}) {
		
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "demo", nil
	})

	render.SetDrawStack(
		render.NewCompositeR(),
	)

	oak.Init("demo")
}

func parseRGBA(args []string) (r, g, b, a int, err error) {
	if len(args) < 4 {
		return
	}
	r, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}
	g, err = strconv.Atoi(args[1])
	if err != nil {
		return
	}
	b, err = strconv.Atoi(args[2])
	if err != nil {
		return
	}
	a, err = strconv.Atoi(args[3])
	return
}

func parseFloats(args []string) (f1, f2 float64, two bool, err error) {
	if len(args) < 1 {
		err = errors.New("No args")
		return
	}
	f1, err = strconv.ParseFloat(args[0], 64)
	if err != nil {
		return
	}
	if len(args) < 2 {
		return
	}
	f2, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return
	}
	two = true
	return
}

func parseInts(args []string) (i1, i2 int, two bool, err error) {
	if len(args) < 1 {
		err = errors.New("No args")
		return
	}
	i1, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}
	if len(args) < 2 {
		return
	}
	i2, err = strconv.Atoi(args[1])
	if err != nil {
		return
	}
	two = true
	return
}
