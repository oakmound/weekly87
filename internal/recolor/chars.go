package recolor

import "image/color"

var (
	MageBase = map[color.RGBA]color.RGBA{
		// Armor
		color.RGBA{104, 20, 11, 255}:   color.RGBA{104, 20, 11, 255},
		color.RGBA{152, 38, 28, 255}:   color.RGBA{152, 38, 28, 255},
		color.RGBA{181, 90, 64, 255}:   color.RGBA{181, 90, 64, 255},
		color.RGBA{230, 199, 160, 255}: color.RGBA{230, 199, 160, 255},
		// Skin
		color.RGBA{124, 99, 65, 255}:   color.RGBA{124, 99, 65, 255},
		color.RGBA{164, 124, 77, 255}:  color.RGBA{164, 124, 77, 255},
		color.RGBA{192, 143, 87, 255}:  color.RGBA{192, 143, 87, 255},
		color.RGBA{205, 167, 117, 255}: color.RGBA{205, 167, 117, 255},
		color.RGBA{221, 190, 148, 255}: color.RGBA{221, 190, 148, 255},
		// Hair
		color.RGBA{59, 20, 27, 255}:  color.RGBA{59, 20, 27, 255},
		color.RGBA{57, 17, 68, 255}:  color.RGBA{57, 17, 68, 255},
		color.RGBA{57, 48, 136, 255}: color.RGBA{57, 48, 136, 255},
		// Eyes
		color.RGBA{39, 87, 7, 255}: color.RGBA{39, 87, 7, 255},
		// Pants
		color.RGBA{116, 116, 116, 255}: color.RGBA{116, 116, 116, 255},
		color.RGBA{117, 116, 114, 255}: color.RGBA{117, 116, 114, 255},
		color.RGBA{147, 146, 144, 255}: color.RGBA{147, 146, 144, 255},
		color.RGBA{169, 169, 169, 255}: color.RGBA{169, 169, 169, 255},
		color.RGBA{170, 169, 167, 255}: color.RGBA{170, 169, 167, 255},
		color.RGBA{200, 199, 197, 255}: color.RGBA{200, 199, 197, 255},
		color.RGBA{226, 225, 223, 255}: color.RGBA{226, 225, 223, 255},
		// Shoes
		color.RGBA{127, 89, 18, 255}:  color.RGBA{127, 89, 18, 255},
		color.RGBA{151, 131, 42, 255}: color.RGBA{151, 131, 42, 255},
		color.RGBA{169, 131, 60, 255}: color.RGBA{169, 131, 60, 255},
		color.RGBA{208, 170, 99, 255}: color.RGBA{208, 170, 99, 255},
		// Outline
		color.RGBA{0, 0, 0, 255}: color.RGBA{0, 0, 0, 255},
	}
	MageNames = map[string]color.RGBA{
		// Armor
		"Armor1": color.RGBA{104, 20, 11, 255},
		"Armor2": color.RGBA{152, 38, 28, 255},
		"Armor3": color.RGBA{181, 90, 64, 255},
		"Armor4": color.RGBA{230, 199, 160, 255},
		// Skin
		"Skin1": color.RGBA{124, 99, 65, 255},
		"Skin2": color.RGBA{164, 124, 77, 255},
		"Skin3": color.RGBA{192, 143, 87, 255},
		"Skin4": color.RGBA{205, 167, 117, 255},
		"Skin5": color.RGBA{221, 190, 148, 255},
		// Hair
		"Hair1": color.RGBA{59, 20, 27, 255},
		"Hair2": color.RGBA{57, 17, 68, 255},
		"Hair3": color.RGBA{57, 48, 136, 255},
		// Eyes
		"Eyes1": color.RGBA{39, 87, 7, 255},
		// Pants
		"Pants1": color.RGBA{116, 116, 116, 255},
		"Pants2": color.RGBA{117, 116, 114, 255},
		"Pants3": color.RGBA{147, 146, 144, 255},
		"Pants4": color.RGBA{169, 169, 169, 255},
		"Pants5": color.RGBA{170, 169, 167, 255},
		"Pants6": color.RGBA{200, 199, 197, 255},
		"Pants7": color.RGBA{226, 225, 223, 255},
		// Shoes
		"Shoes1": color.RGBA{127, 89, 18, 255},
		"Shoes2": color.RGBA{151, 131, 42, 255},
		"Shoes3": color.RGBA{169, 131, 60, 255},
		"Shoes4": color.RGBA{208, 170, 99, 255},
		// Outline
		"Outline1": color.RGBA{0, 0, 0, 255},
	}
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
		color.RGBA{221, 190, 148, 255}: color.RGBA{221, 190, 148, 255},
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
		// Outline
		color.RGBA{0, 0, 0, 255}: color.RGBA{0, 0, 0, 255},
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
		// Outline
		"Outline1": color.RGBA{0, 0, 0, 255},
	}

	// ---------------------------- Specific Recolors ------------------------//

	WarriorTestWhite = map[color.RGBA]color.RGBA{
		// Armor
		color.RGBA{78, 96, 117, 255}:   color.RGBA{255, 255, 255, 255},
		color.RGBA{103, 129, 139, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{121, 138, 157, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{145, 162, 178, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{150, 182, 196, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{174, 211, 221, 255}: color.RGBA{255, 255, 255, 255},
		// Skin
		color.RGBA{124, 99, 65, 255}:   color.RGBA{255, 255, 255, 255},
		color.RGBA{164, 124, 77, 255}:  color.RGBA{255, 255, 255, 255},
		color.RGBA{192, 143, 87, 255}:  color.RGBA{255, 255, 255, 255},
		color.RGBA{205, 167, 117, 255}: color.RGBA{255, 255, 255, 255},
		// Eyes
		color.RGBA{39, 87, 7, 255}: color.RGBA{255, 255, 255, 255},
		// Pants
		color.RGBA{118, 117, 101, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{156, 146, 119, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{183, 168, 135, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{201, 202, 181, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{201, 202, 181, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{226, 227, 214, 255}: color.RGBA{255, 255, 255, 255},
		// Belt
		color.RGBA{171, 171, 171, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{220, 228, 229, 255}: color.RGBA{255, 255, 255, 255},
		color.RGBA{191, 204, 206, 255}: color.RGBA{255, 255, 255, 255},
	}

	WarriorSwordsman = map[color.RGBA]color.RGBA{
		// Armor
		color.RGBA{78, 96, 117, 255}:   color.RGBA{78, 96, 117, 255},
		color.RGBA{103, 129, 139, 255}: color.RGBA{103, 129, 139, 255},
		color.RGBA{121, 138, 157, 255}: color.RGBA{121, 138, 157, 255},
		color.RGBA{145, 162, 178, 255}: color.RGBA{145, 162, 178, 255},
		color.RGBA{150, 182, 196, 255}: color.RGBA{150, 182, 196, 255},
		color.RGBA{174, 211, 221, 255}: color.RGBA{174, 211, 221, 255},
		// Skin
		color.RGBA{124, 99, 65, 255}:   color.RGBA{80, 99, 65, 255},
		color.RGBA{164, 124, 77, 255}:  color.RGBA{104, 124, 77, 255},
		color.RGBA{192, 143, 87, 255}:  color.RGBA{102, 143, 87, 255},
		color.RGBA{205, 167, 117, 255}: color.RGBA{105, 167, 117, 255},
		color.RGBA{221, 190, 148, 255}: color.RGBA{121, 190, 148, 255},
		// Eyes
		color.RGBA{39, 87, 7, 255}: color.RGBA{39, 87, 7, 255},
		// Pants
		color.RGBA{118, 117, 101, 255}: color.RGBA{255, 117, 101, 255},
		color.RGBA{156, 146, 119, 255}: color.RGBA{255, 146, 119, 255},
		color.RGBA{183, 168, 135, 255}: color.RGBA{255, 168, 135, 255},
		color.RGBA{201, 202, 181, 255}: color.RGBA{255, 202, 181, 255},
		color.RGBA{201, 202, 181, 255}: color.RGBA{255, 202, 181, 255},
		color.RGBA{226, 227, 214, 255}: color.RGBA{226, 227, 214, 255},
		// Belt
		color.RGBA{171, 171, 171, 255}: color.RGBA{171, 171, 171, 255},
		color.RGBA{220, 228, 229, 255}: color.RGBA{90, 228, 70, 255},
		color.RGBA{191, 204, 206, 255}: color.RGBA{90, 204, 90, 255},
		// Outline
		color.RGBA{0, 0, 0, 255}: color.RGBA{0, 0, 0, 255},
	}

	WhiteMage = map[color.RGBA]color.RGBA{
		// Armor
		color.RGBA{104, 20, 11, 255}:   color.RGBA{190, 200, 220, 255},
		color.RGBA{152, 38, 28, 255}:   color.RGBA{110, 138, 160, 255},
		color.RGBA{181, 90, 64, 255}:   color.RGBA{190, 190, 200, 255},
		color.RGBA{230, 199, 160, 255}: color.RGBA{200, 199, 200, 255},
		// Skin
		color.RGBA{124, 99, 65, 255}:   color.RGBA{124, 99, 65, 255},
		color.RGBA{164, 124, 77, 255}:  color.RGBA{164, 124, 77, 255},
		color.RGBA{192, 143, 87, 255}:  color.RGBA{192, 143, 87, 255},
		color.RGBA{205, 167, 117, 255}: color.RGBA{130, 167, 117, 255},
		color.RGBA{221, 190, 148, 255}: color.RGBA{130, 190, 148, 255},
		// Hair
		color.RGBA{59, 20, 27, 255}:  color.RGBA{59, 20, 27, 255},
		color.RGBA{57, 17, 68, 255}:  color.RGBA{57, 17, 68, 255},
		color.RGBA{57, 48, 136, 255}: color.RGBA{57, 48, 136, 255},
		// Eyes
		color.RGBA{39, 87, 7, 255}: color.RGBA{39, 87, 7, 255},
		// Pants
		color.RGBA{116, 116, 116, 255}: color.RGBA{116, 116, 116, 255},
		color.RGBA{117, 116, 114, 255}: color.RGBA{117, 116, 114, 255},
		color.RGBA{147, 146, 144, 255}: color.RGBA{147, 146, 144, 255},
		color.RGBA{169, 169, 169, 255}: color.RGBA{169, 169, 169, 255},
		color.RGBA{170, 169, 167, 255}: color.RGBA{170, 169, 167, 255},
		color.RGBA{200, 199, 197, 255}: color.RGBA{200, 199, 197, 255},
		color.RGBA{226, 225, 223, 255}: color.RGBA{226, 225, 223, 255},
		// Shoes
		color.RGBA{127, 89, 18, 255}:  color.RGBA{127, 89, 18, 255},
		color.RGBA{151, 131, 42, 255}: color.RGBA{151, 131, 42, 255},
		color.RGBA{169, 131, 60, 255}: color.RGBA{169, 131, 60, 255},
		color.RGBA{208, 170, 99, 255}: color.RGBA{108, 170, 209, 255},
		// Outline
		color.RGBA{0, 0, 0, 255}: color.RGBA{0, 0, 0, 255},
	}
)
