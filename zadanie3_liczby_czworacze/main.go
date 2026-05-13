package main

import "fmt"

func main() {
	var n int
	fmt.Print("Podaj liczbe: ")
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
