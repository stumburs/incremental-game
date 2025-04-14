package state

import (
	"fmt"
	"incremental-game/src/db"
	"incremental-game/src/theme"
	"math"
	"reflect"
	"time"

	"fyne.io/fyne/v2/canvas"
)

type UIState struct {
	CountText      *canvas.Text
	TechDebtText   *canvas.Text
	BugsText       *canvas.Text
	ClickPowerText *canvas.Text
}

type State struct {
	Count      int
	TechDebt   int
	Bugs       int
	ClickPower int
	ClickDelay time.Duration
	UI         *UIState
}

func NewState(db *db.Database, UI *UIState) *State {

	clickPowerUpgrade, err := db.FetchUpgrade(fmt.Sprintf("%s%s", db.Paths.Upgrades, db.UpgradePaths.ClickPower))
	if err != nil {
		panic("Failed to fetch clickPower from database")
	}
	clickPower := clickPowerUpgrade.Level

	decreaseDelayUpgrade, err := db.FetchUpgrade(fmt.Sprintf("%s%s", db.Paths.Upgrades, db.UpgradePaths.DecreaseDelay))
	if err != nil {
		panic("Failed to fetch decreaseDelay from database")
	}

	return &State{
		Count:      db.GetValue(db.Paths.Count),
		TechDebt:   db.GetValue(db.Paths.Debt),
		Bugs:       db.GetValue(db.Paths.Bugs),
		ClickPower: clickPower,
		ClickDelay: (time.Millisecond * 3000) - (time.Millisecond * (10 * time.Duration(decreaseDelayUpgrade.Level))), // 3000ms default
		UI:         UI,
	}

}

func (state *State) UpdateFromDB(db *db.Database) {
	state.Count = db.GetValue(db.Paths.Count)
	state.TechDebt = db.GetValue(db.Paths.Debt)
	state.Bugs = db.GetValue(db.Paths.Bugs)

	clickPowerUpgrade, err := db.FetchUpgrade(fmt.Sprintf("%s%s", db.Paths.Upgrades, db.UpgradePaths.ClickPower))
	if err != nil {
		panic("Failed to fetch clickPower from database")
	}
	clickPower := clickPowerUpgrade.Level

	// TODO: Simplify
	decreaseDelayUpgrade, err := db.FetchUpgrade(fmt.Sprintf("%s%s", db.Paths.Upgrades, db.UpgradePaths.DecreaseDelay))
	if err != nil {
		panic("Failed to fetch decreaseDelay from database")
	}
	state.ClickDelay = (time.Millisecond * 3000) - (time.Millisecond * (10 * time.Duration(decreaseDelayUpgrade.Level)))

	state.ClickPower = clickPower

	state.UpdateTextFromDB(db)
}

func (state *State) UpdateTextFromDB(db *db.Database) {

	state.UI.ClickPowerText.Text = fmt.Sprintf("Click Power: %d", state.ClickPower)

	state.UI.CountText.Text = (fmt.Sprintf("Count: %d", state.Count))

	state.UI.TechDebtText.Text = (fmt.Sprintf("Tech Debt: %d", state.TechDebt))
	if state.TechDebt > 30 {
		state.UI.TechDebtText.Color = theme.TEXT_YELLOW_400
	} else {
		state.UI.TechDebtText.Color = theme.TEXT_STONE_400
	}

	state.UI.BugsText.Text = (fmt.Sprintf("Bugs: %d", state.Bugs))
	if state.Bugs > 30 {
		state.UI.BugsText.Color = theme.TEXT_RED_400
	} else {
		state.UI.BugsText.Color = theme.TEXT_STONE_400
	}

}

// TODO: Refactor this god-awful way of upgrading
func (state *State) PurchaseUpgrade(db *db.Database, fieldName string) {
	v := reflect.ValueOf(db.UpgradePaths)

	field := v.FieldByName(fieldName)

	if !field.IsValid() {
		fmt.Println("Invalid field name:", fieldName)
		return
	}

	path := field.String()
	fullPath := db.Paths.Upgrades + path

	upgrade, _ := db.FetchUpgrade(fullPath)
	state.UpdateFromDB(db)

	switch field.String() {
	case db.UpgradePaths.ClickPower:
		state.PurchaseClickPower(db, upgrade)
	case db.UpgradePaths.DecreaseDelay:
		state.PurchaseDecreaseDelay(db, upgrade)
	}

	// Refresh state + UI
	state.UpdateFromDB(db)
}

func (state *State) PurchaseClickPower(db *db.Database, upgrade *db.Upgrade) {
	if state.Count < upgrade.Cost {
		fmt.Println("Not enough CP to purchase ClickPower")
		fmt.Printf("Local state: %d Remote state: %d\n", state.Count, db.GetValue(db.Paths.Count))
		return
	}

	// Remove CP
	db.DecrementValueBy(db.Paths.Count, upgrade.Cost)

	// Increase Cost
	db.IncrementValueBy(db.Paths.Upgrades+db.UpgradePaths.ClickPower+db.UpgradeSubpaths.Cost, int(math.Floor(float64(upgrade.Cost)*1.5)))

	// Increase Level
	db.IncrementValueBy(db.Paths.Upgrades+db.UpgradePaths.ClickPower+db.UpgradeSubpaths.Level, 1)
}

func (state *State) PurchaseDecreaseDelay(db *db.Database, upgrade *db.Upgrade) {
	fmt.Println("TODO: PurchaseDecreaseDelay()")
}
