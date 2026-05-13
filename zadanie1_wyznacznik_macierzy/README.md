# Zadanie 1 - Wyznacznik macierzy

Program służy do obliczania wyznacznika macierzy kwadratowej po podaniu jej rozmiaru
oraz wszystkich elementów. Użytkownik wpisuje jedną liczbę określającą wielkość macierzy,
np. `2` dla macierzy 2x2, `3` dla macierzy 3x3 itd.

Program sprawdza poprawność danych wejściowych. Jeśli użytkownik wpisze literę albo
naciśnie Enter bez podania liczby, program wyświetli komunikat błędu i poprosi
o ponowne wpisanie wartości. Rozmiar macierzy musi mieścić się w zakresie od 1 do 10.

Po wprowadzeniu wszystkich elementów program wyświetla wpisaną macierz oraz obliczony
wyznacznik. Dla macierzy 1x1 i 2x2 wynik jest liczony bezpośrednio, a dla większych
macierzy program wykorzystuje rekurencyjne rozwinięcie Laplace'a, tworząc mniejsze
macierze pomocnicze.

## Kod źródłowy

```Go
package main

import (
 "bufio"
 "fmt"
 "os"
 "strconv"
)

const maksRozmiar = 10

func main() {
 scanner := bufio.NewScanner(os.Stdin)

 fmt.Println("=== Wyznacznik macierzy kwadratowej ===")
 fmt.Printf("Podaj rozmiar od 1 do %d, a nastepnie elementy macierzy.\n\n", maksRozmiar)

 size := wczytajRozmiar(scanner)
 matrix := make([][]float64, size)

 fmt.Println("\n--- Uzupelnianie macierzy ---")
 for row := range matrix {
  matrix[row] = make([]float64, size)
  for col := range matrix[row] {
   matrix[row][col] = wczytajLiczbe(scanner, fmt.Sprintf("A[%d][%d]: ", row+1, col+1))
  }
 }

 fmt.Println("\n--- Wprowadzona macierz ---")
 wypiszMacierz(matrix)
 fmt.Printf("\nWyznacznik macierzy wynosi: %.10g\n", wyznacznik(matrix))
}

func wczytajRozmiar(scanner *bufio.Scanner) int {
 for {
  size, err := strconv.Atoi(wczytajTekst(scanner, "Podaj rozmiar macierzy: "))
  if err != nil {
   fmt.Println("Blad: wpisz liczbe calkowita, np. 2 dla macierzy 2x2.")
  } else if size < 1 || size > maksRozmiar {
   fmt.Printf("Blad: rozmiar musi byc w zakresie od 1 do %d.\n", maksRozmiar)
  } else {
   return size
  }
 }
}

func wczytajLiczbe(scanner *bufio.Scanner, prompt string) float64 {
 for {
  number, err := strconv.ParseFloat(wczytajTekst(scanner, prompt), 64)
  if err == nil {
   return number
  }
  fmt.Println("Blad: wpisz liczbe, np. 4 albo 3.5.")
 }
}

func wczytajTekst(scanner *bufio.Scanner, prompt string) string {
 for {
  fmt.Print(prompt)
  if !scanner.Scan() {
   fmt.Println("\nBlad: zakonczono wprowadzanie danych przed uzupelnieniem programu.")
   os.Exit(1)
  }
  if text := scanner.Text(); text != "" {
   return text
  }
  fmt.Println("Blad: wartosc nie moze byc pusta.")
 }
}

func wyznacznik(matrix [][]float64) float64 {
 size := len(matrix)
 if size == 1 {
  return matrix[0][0]
 }
 if size == 2 {
  return matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
 }

 result, sign := 0.0, 1.0
 for col, value := range matrix[0] {
  result += sign * value * wyznacznik(mniejsza(matrix, col))
  sign *= -1
 }
 return result
}

func mniejsza(matrix [][]float64, skippedCol int) [][]float64 {
 result := make([][]float64, 0, len(matrix)-1)
 for row := 1; row < len(matrix); row++ {
  minorRow := make([]float64, 0, len(matrix)-1)
  for col, value := range matrix[row] {
   if col != skippedCol {
    minorRow = append(minorRow, value)
   }
  }
  result = append(result, minorRow)
 }
 return result
}

func wypiszMacierz(matrix [][]float64) {
 for _, row := range matrix {
  fmt.Print("|")
  for _, value := range row {
   fmt.Printf(" %10.10g", value)
  }
  fmt.Println(" |")
 }
}

```
