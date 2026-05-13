# Zadanie 2 - Liczby pierwsze

Program wczytuje od użytkownika jedną liczbę naturalną `n` i wypisuje wszystkie liczby
pierwsze nie większe od `n`. Do wyznaczania liczb pierwszych używa sita Eratostenesa.

Jeżeli użytkownik poda błędne dane albo liczbę mniejszą od `1`, program wyświetli komunikat:
`Podaj poprawne dane`.

## Kod źródłowy

```go
package main

import "fmt"

func main() {
 var n int
 if _, err := fmt.Scan(&n); err != nil || n < 1 {
  fmt.Println("Podaj poprawne dane")
  return
 }

 sito := make([]bool, n+1)
 for i := 2; i*i <= n; i++ {
  if !sito[i] {
   for j := i * i; j <= n; j += i {
    sito[j] = true
   }
  }
 }

 for i := 2; i <= n; i++ {
  if !sito[i] {
   fmt.Println(i)
  }
 }
}
```
