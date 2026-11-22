// Package vreme omogoča dodajanje in spremljanje temperaturnih meritev različnih vremenskih postaj.

package vreme

import (
	"fmt"
)

// Postaja predstavlja eno vremensko postajo.
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// postaje predstavlja seznam zabeleženih postaj.
var postaje = make(map[string]Postaja)

var minTemp float64 = -50
var maxTemp float64 = 50
var stMeritev int = 3

// NastaviMin nastavi najmanjšo veljavno temperaturo.
func NastaviMin(value float64) {
	minTemp = value
}

// NastaviMax nastavi največjo veljavno temperaturo.
func NastaviMax(value float64) {
	maxTemp = value
}

// NastaviStMeritev nastavi najmanjše število temperaturnih meritev potrebnih za izračun statistike postaje.
func NastaviStMeritev(value int) {
	stMeritev = value
}

// IzpisPorocila vrne poročilo o največji, najmanjši in povprečni temperaturi vsake postaje.
func IzpisPorocila() string {
	var izpis string = ""

	for k, v := range postaje {
		min, max, avg, ok := statistika(k)
		var opis string

		if !ok {
			izpis = izpis + fmt.Sprintf("%s: premalo meritev\n", v.ime)

		} else {
			switch {
			case avg >= 25:
				opis = "Vroče"
			case avg >= 15:
				opis = "Toplo"
			case avg >= 5:
				opis = "Hladno"
			case avg < 5:
				opis = "Mrzlo"
			}

			izpis = izpis + fmt.Sprintf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n", v.ime, min, max, avg, opis)
		}
	}

	if izpis == "" {
		izpis = "Dodana ni bila še nobena postaja\n"
	}

	return izpis
}

// IzpisMeritev vrne izpis, ki vključuje oznako, ime, kraj in zabeležene meritve za vsako postajo.
func IzpisMeritev() string {
	var izpis string = ""

	for k, v := range postaje {
		izpis = izpis + fmt.Sprintf("%s - %s (%s): %v\n", k, v.ime, v.kraj, v.meritve)
	}

	if izpis == "" {
		izpis = "Dodana ni bila še nobena postaja\n"
	}

	return izpis
}

func statistika(oznaka string) (najnizja, najvisja, povprecna float64, ok bool) {
	p, obstaja := postaje[oznaka]
	st_meritev := len(p.meritve)

	if !obstaja || st_meritev < stMeritev {
		ok = false
		return
	}

	ok = true

	najnizja = 51
	najvisja = -51
	vsota := 0.0

	for i := range st_meritev {
		m := p.meritve[i]
		vsota = vsota + m

		if m < najnizja {
			najnizja = m
		}

		if m > najvisja {
			najvisja = m
		}
	}

	povprecna = vsota / float64(st_meritev)
	return
}

// DodajMeritev doda meritev temperatura postaji oznaka.
// Vrne potrdilo o dodani meritvi ali napako, če je oznaka ali temperatura neveljavna.
func DodajMeritev(oznaka string, temperatura float64) (string, error) {
	p, ok := postaje[oznaka]

	if temperatura < minTemp || temperatura > maxTemp {
		return "", fmt.Errorf("Napaka: Neveljavna temperatura")
	}

	if !ok {
		return "", fmt.Errorf("Napaka: Postaja z oznako \"%s\" ne obstaja\n", oznaka)
	}

	p.meritve = append(p.meritve, temperatura)
	postaje[oznaka] = p
	return fmt.Sprintf("%s: Uspešno dodana temperatura\n", oznaka), nil
}

// DodajPostajo v seznam postaj doda postajo z določeno oznako, imenom in krajem.
// Če postaja oznaka že obstaja ali sta ime/kraj neveljavna, vrne napako.
func DodajPostajo(oznaka string, ime string, kraj string) error {
	_, ok := postaje[oznaka]

	if ok {
		return fmt.Errorf("Postaja z oznako %s že obstaja\n", oznaka)
	}

	if ime == "" || kraj == "" {
		return fmt.Errorf("Prazno ime ali kraj postaje\n")
	}

	postaje[oznaka] = Postaja{ime: ime, kraj: kraj, meritve: []float64{}}

	return nil
}
