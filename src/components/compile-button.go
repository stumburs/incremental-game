package components

import (
	"incremental-game/src/db"
	"incremental-game/src/state"
	"math"
	"math/rand/v2"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func StartMainStatRefresher(db *db.Database) {

}

func NewCompileButton(db *db.Database, state *state.State, clickPower int) *fyne.Container {
	var button *widget.Button

	buttonContainer := container.NewVBox()

	button = widget.NewButton("Compile", func() {
		// Change look to progress bar when pressed
		progressBar := widget.NewProgressBar()
		progressBar.Min = 0
		progressBar.Max = 1
		progressBar.Value = 0
		buttonContainer.Objects = []fyne.CanvasObject{progressBar}
		buttonContainer.Refresh()

		// Increment
		go func() {

			var wg sync.WaitGroup
			wg.Add(2)

			// Update progress bar
			go func() {
				defer wg.Done()
				ticker := time.NewTicker(state.ClickDelay / 100) // 100 ticks total
				defer ticker.Stop()

				for i := 0; i <= 100; i++ {
					<-ticker.C
					progressBar.SetValue(float64(i) / 100)
				}
			}()

			go func() {
				defer wg.Done()

				// Bugs
				if rand.Float32() < 0.15 {
					db.IncrementValueBy(db.Paths.Bugs, 1)
					db.SetValue(db.Paths.Count, func() int {
						currentCount := db.GetValue(db.Paths.Count)
						return int(math.Floor(float64(currentCount) * 0.9))
					}())
				}

				// Debt + count
				db.IncrementValueBy(db.Paths.Debt, 1)
				db.IncrementValueBy(db.Paths.Count, state.ClickPower)
			}()

			wg.Wait()

			// TODO: Update text
			state.UpdateFromDB(db)

			// Reset to button
			buttonContainer.Objects = []fyne.CanvasObject{button}
			buttonContainer.Refresh()
		}()
	})

	buttonContainer.Add(button)

	return buttonContainer
}
