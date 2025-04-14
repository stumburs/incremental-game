package components

import (
	"image/color"
	"incremental-game/src/db"
	"incremental-game/src/theme"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func CreateUpgradeCard(upgrade *db.Upgrade, onClick func()) fyne.CanvasObject {

	// Left side
	// Description
	titleText := canvas.NewText(upgrade.Name, color.White)
	titleText.TextSize = 16
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	// Level
	levelText := canvas.NewText("Level: "+strconv.Itoa(upgrade.Level), theme.TEXT_STONE_300)
	levelText.TextSize = 14

	// Right side
	// Price
	priceText := canvas.NewText(strconv.Itoa(upgrade.Cost)+" CP", theme.TEXT_STONE_400)
	priceText.TextSize = 14

	// Left-aligned text group
	textCol := container.NewVBox(titleText, levelText)

	// Horizontal layout
	content := container.NewHBox(
		textCol,
		layout.NewSpacer(),
		priceText,
	)

	bg := canvas.NewRectangle(color.RGBA{R: 40, G: 40, B: 40, A: 255})

	padded := container.NewPadded(content)
	card := container.NewStack(bg, padded)

	// Clickable with invisible button
	button := widget.NewButton("", onClick)
	button.Importance = widget.LowImportance
	buttonContainer := container.NewStack(button, card)

	return buttonContainer
}
