# Zadanie 4 - Bisekcja dla wielomianu 5. stopnia

Program wczytuje od użytkownika współczynniki wielomianu:

`f(x) = ax⁵ + bx⁴ + cx³ + dx² + ex + f`

oraz przedział, na którym ma szukać miejsc zerowych.

Program prosi kolejno o:

- `a`
- `b`
- `c`
- `d`
- `e`
- `f`
- przedział w postaci: początek i koniec

Następnie:

- wyznacza punkty krytyczne pochodnej,
- dzieli przedział na mniejsze fragmenty,
- stosuje metodę bisekcji tam, gdzie występuje zmiana znaku,
- wypisuje wszystkie znalezione miejsca zerowe w podanym przedziale.

Jeżeli dane są błędne albo współczynnik `a = 0`, program wyświetli komunikat:
`Podaj poprawne dane`.

Jeżeli w podanym przedziale nie ma miejsc zerowych, program wyświetli komunikat:
`Brak miejsc zerowych w podanym przedziale`.

## Kod źródłowy

```go
package main

import (
 "fmt"
 "math"
 "sort"
)

const (
 scanStep          = 0.001
 epsilon           = 0.000001
 duplicateDistance = 0.0001
 maxIterations     = 1000
)

func main() {
 wspolczynniki, lewy, prawy, err := wczytajDane()
 if err != nil || wspolczynniki[0] == 0 {
  fmt.Println("Podaj poprawne dane")
  return
 }

 if lewy > prawy {
  lewy, prawy = prawy, lewy
 }

 if lewy == prawy {
  fmt.Println("Podaj poprawne dane")
  return
 }

 miejscaZerowe := znajdzMiejscaZeroweWielomianu(wspolczynniki, lewy, prawy)
 if len(miejscaZerowe) == 0 {
  fmt.Println("Brak miejsc zerowych w podanym przedziale")
  return
 }

 fmt.Println("Znalezione miejsca zerowe:")
 for i, miejsceZerowe := range miejscaZerowe {
  fmt.Printf("x%d = %.6f\n", i+1, miejsceZerowe)
 }
}

func wczytajDane() ([]float64, float64, float64, error) {
 wspolczynniki := make([]float64, 6)
 etykiety := []string{"a", "b", "c", "d", "e", "f"}

 for i, etykieta := range etykiety {
  fmt.Printf("Podaj %s: ", etykieta)
  if _, err := fmt.Scan(&wspolczynniki[i]); err != nil {
   return nil, 0, 0, err
  }
 }

 var lewy, prawy float64
 fmt.Print("Podaj przedzial (poczatek koniec): ")
 if _, err := fmt.Scan(&lewy, &prawy); err != nil {
  return nil, 0, 0, err
 }

 return wspolczynniki, lewy, prawy, nil
}

func znajdzMiejscaZeroweWielomianu(wspolczynniki []float64, lewy, prawy float64) []float64 {
 miejscaZerowePochodnej := skanujMiejscaZerowe(obliczPochodna(wspolczynniki), lewy, prawy)
 punkty := append([]float64{lewy, prawy}, miejscaZerowePochodnej...)

 sort.Float64s(punkty)
 punkty = usunPowtorzeniaZPosortowanych(punkty)

 miejscaZerowe := make([]float64, 0)

 for _, punkt := range punkty {
  if math.Abs(obliczWartoscWielomianu(wspolczynniki, punkt)) <= epsilon {
   miejscaZerowe = dolaczJesliUnikalne(miejscaZerowe, punkt)
  }
 }

 for i := 0; i < len(punkty)-1; i++ {
  lewyPunkt := punkty[i]
  prawyPunkt := punkty[i+1]

  if prawyPunkt-lewyPunkt <= epsilon {
   continue
  }

  lewaWartosc := obliczWartoscWielomianu(wspolczynniki, lewyPunkt)
  prawaWartosc := obliczWartoscWielomianu(wspolczynniki, prawyPunkt)

  if lewaWartosc*prawaWartosc < 0 {
   miejsceZerowe := bisekcja(wspolczynniki, lewyPunkt, prawyPunkt)
   miejscaZerowe = dolaczJesliUnikalne(miejscaZerowe, miejsceZerowe)
  }
 }

 sort.Float64s(miejscaZerowe)
 return miejscaZerowe
}

func skanujMiejscaZerowe(wspolczynniki []float64, lewy, prawy float64) []float64 {
 miejscaZerowe := make([]float64, 0)
 dlugoscPrzedzialu := prawy - lewy
 liczbaKrokow := int(math.Ceil(dlugoscPrzedzialu / scanStep))

 for i := 0; i < liczbaKrokow; i++ {
  poczatek := lewy + float64(i)*scanStep
  koniec := lewy + float64(i+1)*scanStep

  if koniec > prawy {
   koniec = prawy
  }

  wartoscPoczatku := obliczWartoscWielomianu(wspolczynniki, poczatek)
  wartoscKonca := obliczWartoscWielomianu(wspolczynniki, koniec)

  if math.Abs(wartoscPoczatku) <= epsilon {
   miejscaZerowe = dolaczJesliUnikalne(miejscaZerowe, poczatek)
  }

  if math.Abs(wartoscKonca) <= epsilon {
   miejscaZerowe = dolaczJesliUnikalne(miejscaZerowe, koniec)
  }

  if wartoscPoczatku*wartoscKonca < 0 {
   miejsceZerowe := bisekcja(wspolczynniki, poczatek, koniec)
   miejscaZerowe = dolaczJesliUnikalne(miejscaZerowe, miejsceZerowe)
  }
 }

 return miejscaZerowe
}

func bisekcja(wspolczynniki []float64, lewy, prawy float64) float64 {
 lewaWartosc := obliczWartoscWielomianu(wspolczynniki, lewy)
 prawaWartosc := obliczWartoscWielomianu(wspolczynniki, prawy)
 poprzedniSrodek := lewy

 for i := 0; i < maxIterations; i++ {
  srodek := (lewy + prawy) / 2
  wartoscSrodka := obliczWartoscWielomianu(wspolczynniki, srodek)

  if math.Abs(wartoscSrodka) <= epsilon || math.Abs(srodek-poprzedniSrodek) <= epsilon {
   return srodek
  }

  if lewaWartosc*wartoscSrodka < 0 {
   prawy = srodek
   prawaWartosc = wartoscSrodka
  } else {
   lewy = srodek
   lewaWartosc = wartoscSrodka
  }

  if math.Abs(prawy-lewy) <= epsilon || math.Abs(prawaWartosc-lewaWartosc) <= epsilon {
   return (lewy + prawy) / 2
  }

  poprzedniSrodek = srodek
 }

 return (lewy + prawy) / 2
}

func obliczWartoscWielomianu(wspolczynniki []float64, x float64) float64 {
 wynik := 0.0
 for _, wspolczynnik := range wspolczynniki {
  wynik = wynik*x + wspolczynnik
 }
 return wynik
}

func obliczPochodna(wspolczynniki []float64) []float64 {
 stopien := len(wspolczynniki) - 1
 pochodna := make([]float64, 0, stopien)

 for i := 0; i < stopien; i++ {
  wykladnik := float64(stopien - i)
  pochodna = append(pochodna, wspolczynniki[i]*wykladnik)
 }

 return pochodna
}

func usunPowtorzeniaZPosortowanych(wartosci []float64) []float64 {
 if len(wartosci) == 0 {
  return wartosci
 }

 unikalne := []float64{wartosci[0]}
 for i := 1; i < len(wartosci); i++ {
  if math.Abs(wartosci[i]-unikalne[len(unikalne)-1]) > duplicateDistance {
   unikalne = append(unikalne, wartosci[i])
  }
 }

 return unikalne
}

func dolaczJesliUnikalne(miejscaZerowe []float64, miejsceZerowe float64) []float64 {
 for _, istniejace := range miejscaZerowe {
  if math.Abs(istniejace-miejsceZerowe) <= duplicateDistance {
   return miejscaZerowe
  }
 }

 return append(miejscaZerowe, miejsceZerowe)
}

```
