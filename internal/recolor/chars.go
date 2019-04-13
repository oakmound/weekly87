package recolor

import "image/color"

var (
	WarriorBase = map[color.RGBA]color.RGBA{
		// Armor
		color.RGBA{78, 96, 117, 255}:   color.RGBA{78, 96, 117, 255},
		color.RGBA{103, 129, 139, 255}: color.RGBA{103, 129, 139, 255},
		color.RGBA{121, 138, 157, 255}: color.RGBA{121, 138, 157, 255},
		color.RGBA{145, 162, 178, 255}: color.RGBA{145, 162, 178, 255},
		color.RGBA{150, 182, 196, 255}: color.RGBA{150, 182, 196, 255},
		color.RGBA{174, 211, 221, 255}: color.RGBA{174, 211, 221, 255},
		// Skin
		color.RGBA{124, 99, 65, 255}:   color.RGBA{124, 99, 65, 255},
		color.RGBA{164, 124, 77, 255}:  color.RGBA{164, 124, 77, 255},
		color.RGBA{192, 143, 87, 255}:  color.RGBA{192, 143, 87, 255},
		color.RGBA{205, 167, 117, 255}: color.RGBA{205, 167, 117, 255},
		// Eyes
		color.RGBA{39, 87, 7, 255}: color.RGBA{39, 87, 7, 255},
		// Pants
		color.RGBA{118, 117, 101, 255}: color.RGBA{118, 117, 101, 255},
		color.RGBA{156, 146, 119, 255}: color.RGBA{156, 146, 119, 255},
		color.RGBA{183, 168, 135, 255}: color.RGBA{183, 168, 135, 255},
		color.RGBA{201, 202, 181, 255}: color.RGBA{201, 202, 181, 255},
		color.RGBA{201, 202, 181, 255}: color.RGBA{201, 202, 181, 255},
		color.RGBA{226, 227, 214, 255}: color.RGBA{226, 227, 214, 255},
		// Belt
		color.RGBA{171, 171, 171, 255}: color.RGBA{171, 171, 171, 255},
		color.RGBA{220, 228, 229, 255}: color.RGBA{220, 228, 229, 255},
		color.RGBA{191, 204, 206, 255}: color.RGBA{191, 204, 206, 255},
	}
	WarriorNames = map[string]color.RGBA{
		// Armor
		"Armor1": color.RGBA{78, 96, 117, 255},
		"Armor2": color.RGBA{103, 129, 139, 255},
		"Armor3": color.RGBA{121, 138, 157, 255},
		"Armor4": color.RGBA{145, 162, 178, 255},
		"Armor5": color.RGBA{150, 182, 196, 255},
		"Armor6": color.RGBA{174, 211, 221, 255},
		// Skin
		"Skin1": color.RGBA{124, 99, 65, 255},
		"Skin2": color.RGBA{164, 124, 77, 255},
		"Skin3": color.RGBA{192, 143, 87, 255},
		"Skin4": color.RGBA{205, 167, 117, 255},
		// Eyes
		"Eyes1": color.RGBA{39, 87, 7, 255},
		// Pants
		"Pants1": color.RGBA{118, 117, 101, 255},
		"Pants2": color.RGBA{156, 146, 119, 255},
		"Pants3": color.RGBA{183, 168, 135, 255},
		"Pants4": color.RGBA{201, 202, 181, 255},
		"Pants5": color.RGBA{201, 202, 181, 255},
		"Pants6": color.RGBA{226, 227, 214, 255},
		// Belt
		"Belt1": color.RGBA{171, 171, 171, 255},
		"Belt2": color.RGBA{191, 204, 206, 255},
		"Belt3": color.RGBA{220, 228, 229, 255},
	}
)
