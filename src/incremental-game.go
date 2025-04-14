package main

import (
	"fmt"
	"incremental-game/src/components"
	"incremental-game/src/db"
	"incremental-game/src/state"
	"incremental-game/src/theme"
	"reflect"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func buildAllUpgradeCards(db *db.Database, basePath string, state state.State) []fyne.CanvasObject {
	var cards []fyne.CanvasObject

	val := reflect.ValueOf(db.UpgradePaths)
	typ := reflect.TypeOf(db.UpgradePaths)

	for i := range val.NumField() {
		fieldName := typ.Field(i).Name
		subPath := val.Field(i).String()

		fullPath := fmt.Sprintf("%s%s", basePath, subPath)

		upgrade, err := db.FetchUpgrade(fullPath)
		if err != nil {
			fmt.Printf("Failed to fetch upgrade %s: %v\n", fieldName, err)
			continue
		}

		card := components.CreateUpgradeCard(upgrade, func() {
			state.PurchaseUpgrade(db, fieldName)
		})

		cards = append(cards, card)

	}

	return cards
}

// TODO: Refactor this probably
var updateChan = make(chan []fyne.CanvasObject)

func StartUpgradeCardRefresher(db *db.Database, state state.State) {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			newCards := buildAllUpgradeCards(db, db.Paths.Upgrades, state)
			updateChan <- newCards
		}
	}()
}

func main() {

	db := db.NewDatabase()

	a := app.New()
	a.Settings().SetTheme(&theme.CustomTheme{})

	w := a.NewWindow("Incremental Game")
	w.Resize(fyne.NewSize(1280, 800))

	countText := canvas.NewText("Count: Loading...", theme.TEXT_STONE_300)
	countText.Alignment = fyne.TextAlignCenter

	techDebtText := canvas.NewText("Tech Debt: Loading...", theme.TEXT_STONE_400)
	techDebtText.Alignment = fyne.TextAlignCenter
	techDebtText.TextSize = 22

	bugsText := canvas.NewText("Bugs: Loading...", theme.TEXT_STONE_400)
	bugsText.Alignment = fyne.TextAlignCenter
	bugsText.TextSize = 22

	clickPowerText := canvas.NewText("Click Power: Loading...", theme.TEXT_STONE_400)
	clickPowerText.Alignment = fyne.TextAlignCenter
	clickPowerText.TextSize = 22

	ui := &state.UIState{
		CountText:      countText,
		TechDebtText:   techDebtText,
		BugsText:       bugsText,
		ClickPowerText: clickPowerText,
	}

	state := state.NewState(&db, ui)

	compileButtonContainer := components.NewCompileButton(&db, state, state.ClickPower)

	topContainer := container.NewCenter(container.NewVBox(state.UI.CountText, state.UI.TechDebtText, state.UI.BugsText, compileButtonContainer, state.UI.ClickPowerText))

	upgradesText := canvas.NewText("Upgrades", theme.TEXT_STONE_300)

	upgradeCards := buildAllUpgradeCards(&db, db.Paths.Upgrades, *state)

	buttonGrid := container.NewGridWithColumns(3, upgradeCards...)

	bottomContainer := container.NewVBox(upgradesText, buttonGrid)

	scrollableUpgrades := container.NewVScroll(bottomContainer)

	verticalSplit := container.NewVSplit(topContainer, scrollableUpgrades)
	verticalSplit.Offset = 0.6

	w.SetContent(verticalSplit)

	// Continuously update upgrade buttons every 5 seconds (time set in StartUpgradeCardRefresher using channel)
	StartUpgradeCardRefresher(&db, *state)
	go func() {
		for newCards := range updateChan {
			buttonGrid.Objects = newCards
			buttonGrid.Refresh()
		}
	}()

	// Continuously update state every 5 seconds
	go func() {
		for {
			state.UpdateFromDB(&db)
			time.Sleep(5 * time.Second)
		}
	}()

	w.ShowAndRun()
}
