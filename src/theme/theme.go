package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var TEXT_RED_400 = color.RGBA{R: 0xFF, G: 0x64, B: 0x67, A: 0xFF}
var TEXT_STONE_300 = color.RGBA{R: 0xD6, G: 0xD3, B: 0xD1, A: 0xFF}
var TEXT_STONE_400 = color.RGBA{R: 0xA8, G: 0xA2, B: 0x9E, A: 0xFF}
var TEXT_YELLOW_400 = color.RGBA{R: 0xFC, G: 0xC8, B: 0x00, A: 0xFF}

type CustomTheme struct {
	fyne.Theme
}

var _ fyne.Theme = (*CustomTheme)(nil)

func (c CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	// Background color
	if name == theme.ColorNameBackground {
		if variant == theme.VariantDark {
			return color.RGBA{R: 0x0c, G: 0x0a, B: 0x09, A: 0xff}
		}
	}

	// Button color
	if name == theme.ColorNameButton {
		if variant == theme.VariantDark {
			return color.RGBA{R: 0x1f, G: 0x3b, B: 0x8a, A: 0xff}
		}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (c CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (c CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if name == theme.IconNameHome {
		return fyne.NewStaticResource("myHome", nil)
	}

	return theme.DefaultTheme().Icon(name)
}

func (c CustomTheme) Size(name fyne.ThemeSizeName) float32 {

	if name == theme.SizeNameText {
		return 40
	}
	return theme.DefaultTheme().Size(name)
}
