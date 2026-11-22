package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v3"
	"github.com/zojaev/ps_vreme/vreme"
)

func main() {
	cmd := &cli.Command{
		Name:  "Vremenske postaje",
		Usage: "Aplikacija upravlja z vremenskimi postajami in njihovimi meritvami",
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:  "minTemp",
				Usage: "Najmanjša veljavna temperatura",
				Value: -50,
			},
			&cli.Float64Flag{
				Name:  "maxTemp",
				Usage: "Največja veljavna temperatura",
				Value: 50,
			},
			&cli.IntFlag{
				Name:  "stMeritev",
				Usage: "Najmanjše število meritev potrebnih za izračun statistike postaje",
				Value: 3,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			minTemp := cmd.Float64("minTemp")
			maxTemp := cmd.Float64("maxTemp")
			stMeritev := cmd.Int("stMeritev")
			return runTUI(minTemp, maxTemp, stMeritev)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func exampleValues() {
	vreme.DodajPostajo("LJ01", "Bežigrad", "Ljubljana")
	vreme.DodajMeritev("LJ01", 18.1)
	vreme.DodajMeritev("LJ01", 22.3)
	vreme.DodajMeritev("LJ01", 17.8)
	vreme.DodajMeritev("LJ01", 19.0)

	vreme.DodajPostajo("LJ02", "Vič", "Ljubljana")
	vreme.DodajMeritev("LJ02", 10.3)
	vreme.DodajMeritev("LJ02", 9.5)
	vreme.DodajMeritev("LJ02", 11.6)

	vreme.DodajPostajo("MB01", "Pohorje", "Maribor")
	vreme.DodajMeritev("MB01", 5.1)
}

func runTUI(minTemp float64, maxTemp float64, stMeritev int) error {
	vreme.NastaviMax(maxTemp)
	vreme.NastaviMin(minTemp)
	vreme.NastaviStMeritev(stMeritev)

	exampleValues()

	app := tview.NewApplication()
	app.EnableMouse(true)

	mainView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})
	mainView.SetBorder(true)

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[gray]Press [white]q[gray] to quit").
		SetTextAlign(tview.AlignLeft)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()
			return nil
		}
		return event
	})

	buttonMeritve := tview.NewButton("Prikaži meritve")
	buttonPorocilo := tview.NewButton("Prikaži poročilo")

	buttonMeritve.SetSelectedFunc(func() {
		mainView.SetText(vreme.IzpisMeritev())
	})

	buttonPorocilo.SetSelectedFunc(func() {
		mainView.SetText(vreme.IzpisPorocila())
	})

	buttons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(buttonMeritve, 0, 1, false).
		AddItem(buttonPorocilo, 0, 1, false)

	formPostaja := tview.NewForm()
	formPostaja.
		AddInputField("Oznaka:", "", 20, nil, nil).
		AddInputField("Ime:", "", 20, nil, nil).
		AddInputField("Kraj:", "", 20, nil, nil).
		AddButton("Dodaj postajo", func() {
			oznaka := formPostaja.GetFormItemByLabel("Oznaka:").(*tview.InputField).GetText()
			ime := formPostaja.GetFormItemByLabel("Ime:").(*tview.InputField).GetText()
			kraj := formPostaja.GetFormItemByLabel("Kraj:").(*tview.InputField).GetText()

			err := vreme.DodajPostajo(oznaka, ime, kraj)

			if err != nil {
				mainView.SetText(err.Error())
				return
			}

			mainView.SetText("Postaja dodana")

			formPostaja.GetFormItemByLabel("Oznaka:").(*tview.InputField).SetText("")
			formPostaja.GetFormItemByLabel("Ime:").(*tview.InputField).SetText("")
			formPostaja.GetFormItemByLabel("Kraj:").(*tview.InputField).SetText("")
		})

	formMeritev := tview.NewForm()
	formMeritev.
		AddInputField("Oznaka:", "", 20, nil, nil).
		AddInputField("Temperatura:", "", 20, nil, nil).
		AddButton("Dodaj", func() {
			oznaka := formMeritev.GetFormItemByLabel("Oznaka:").(*tview.InputField).GetText()
			temperaturaStr := formMeritev.GetFormItemByLabel("Temperatura:").(*tview.InputField).GetText()
			temperatura, err := strconv.ParseFloat(temperaturaStr, 64)

			if err != nil {
				mainView.SetText("Neveljavna temperatura")
				return
			}

			izpis, err := vreme.DodajMeritev(oznaka, temperatura)

			if err != nil {
				mainView.SetText(err.Error())
				return
			}

			mainView.SetText(izpis)

			formMeritev.GetFormItemByLabel("Oznaka:").(*tview.InputField).SetText("")
			formMeritev.GetFormItemByLabel("Temperatura:").(*tview.InputField).SetText("")

		})

	forms := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(formPostaja, 0, 1, false).
		AddItem(formMeritev, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainView, 0, 1, true).
		AddItem(forms, 10, 0, false).
		AddItem(buttons, 3, 0, false).
		AddItem(footer, 3, 0, false)

	return app.SetRoot(layout, true).Run()
}
