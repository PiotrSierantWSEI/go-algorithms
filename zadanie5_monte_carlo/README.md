# Zadanie 5 - Całkowanie metodą Monte Carlo

Program oblicza wartość całki oznaczonej funkcji:

`f(x) = |sin(x) + sin(2x) + sin(4x) + sin(8x)|`

na przedziale `[0, 2π]` metodą Monte Carlo.

Program prosi użytkownika o:

- `epsilon` - dokładność obliczeń (np. `0.000001`)

## Działanie programu

Metoda Monte Carlo polega na losowaniu punktów wewnątrz prostokąta otaczającego wykres funkcji i szacowaniu pola powierzchni pod krzywą na podstawie proporcji punktów, które się pod nią znalazły.

### Krok 1 - Wyznaczenie prostokąta

Funkcja `f(x)` spełnia warunek `0 ≤ f(x) ≤ 4`, dlatego prostokąt ma wymiary:

- `x ∈ (0, 2π)` - oś pozioma
- `y ∈ (0, 4)`  - oś pionowa
- pole prostokąta = `2π × 4 = 8π`

### Krok 2 - Losowanie punktów

Program losuje punkt `P(x, y)` o współrzędnych:

- `x` - liczba losowa z przedziału `(0, 2π)`
- `y` - liczba losowa z przedziału `(0, 4)`

Następnie sprawdza warunek: czy punkt leży pod wykresem funkcji, tzn. `y ≤ f(x)`.

- Jeśli tak: licznik `L1` jest zwiększany o 1
- Jeśli nie: licznik `L2` jest zwiększany o 1

### Krok 3 - Szacowanie całki

Po każdym wylosowanym punkcie wartość całki jest szacowana wzorem:

```go
∫₀^{2π} f(x)dx ≈ 8π × L1 / (L1 + L2)
```

### Krok 4 - Warunek stopu

Program kończy działanie, gdy różnica między kolejnymi szacowaniami jest mniejsza od epsilon:

```go
|calka(N+1) - calka(N)| < epsilon
```

Minimalnie generowanych jest 1000 punktów, żeby uniknąć przedwczesnego zakończenia.

## Przykładowe uruchomienie

```go
Podaj epsilon: 0.000001
Wynik calki: 7.443840
Liczba punktow: 7443841
```

## Kod źródłowy

```go
package main

import (
 "fmt"
 "math"
 "math/rand"
 "time"
)

func f(x float64) float64 {
 return math.Abs(math.Sin(x) + math.Sin(2*x) + math.Sin(4*x) + math.Sin(8*x))
}

func main() {
 var epsilon float64
 fmt.Print("Podaj epsilon: ")
 fmt.Scan(&epsilon)

 rng := rand.New(rand.NewSource(time.Now().UnixNano()))

 // Prostokąt: x ∈ (0, 2π), y ∈ (0, 4), pole = 8π
 rectangleArea := 8 * math.Pi

 var L1, L2 int64 // L1: pod krzywą, L2: nad krzywą
 var prevIntegral float64
 var currentIntegral float64
 const minPoints = 1000

 for {
  x := rng.Float64() * 2 * math.Pi
  y := rng.Float64() * 4

  if y <= f(x) {
   L1++
  } else {
   L2++
  }

  total := L1 + L2
  currentIntegral = rectangleArea * float64(L1) / float64(total)

  if total >= minPoints && math.Abs(currentIntegral-prevIntegral) < epsilon {
   break
  }

  prevIntegral = currentIntegral
 }

 fmt.Printf("Wynik całki: %.6f\n", currentIntegral)
 fmt.Printf("Liczba punktów: %d\n", L1+L2)
}

```
