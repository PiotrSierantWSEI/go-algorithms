# Zadanie 3 - Liczby czworacze

Program wczytuje od użytkownika jedną liczbę naturalną `n` i wypisuje wszystkie liczby
czworacze mniejsze od `n`. Liczby czworacze to czwórki liczb pierwszych w postaci
`p`, `p+2`, `p+6`, `p+8`.

Do wyznaczania liczb pierwszych program używa sita Eratostenesa, a potem sprawdza,
czy dla danej liczby `p` również `p+2`, `p+6` i `p+8` są pierwsze.

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

 for p := 2; p+8 < n; p++ {
  if !sito[p] && !sito[p+2] && !sito[p+6] && !sito[p+8] {
   fmt.Println(p, p+2, p+6, p+8)
  }
 }
}
```
