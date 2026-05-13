# Algorytmy w Go

Repozytorium zawiera rozwiązania zadań algorytmicznych w języku Go.

## Struktura projektu

- zadanie1_wyznacznik_macierzy
- zadanie2_liczby_pierwsze
- zadanie3_liczby_czworacze
- zadanie4_bisekcja_wielomian5
- zadanie5_monte_carlo

## Generator sprawozdania PDF

W katalogu głównym znajduje się skrypt [generate_report.go](generate_report.go), który:

- skanuje wszystkie katalogi zawierające `README.md` i plik źródłowy (`main.go`, `main.py`, `main.ts`, `main.js` lub dowolny plik `.go`, `.py`, `.ts`, `.js`),
- zbiera opisy z plików README oraz kod źródłowy,
- generuje spis treści z linkami do poszczególnych zadań,
- zapisuje raport jako plik PDF.

Uruchomienie:

```bash
go run ./generate_report.go -out sprawozdanie_algorytmy.pdf "Tytuł sprawozdania" "Imię Nazwisko" 123456
```

Argumenty:

| Argument | Opis |
| --- | --- |
| `-out` | Ścieżka do pliku wyjściowego (domyślnie: `sprawozdanie_algorytmy.pdf`) |
| `"Tytuł sprawozdania"` | Tytuł wyświetlany na stronie tytułowej |
| `"Imię Nazwisko"` | Imię i nazwisko studenta |
| `123456` | Numer indeksu studenta |

Przykład:

```bash
go run ./generate_report.go -out sprawozdanie_algorytmy.pdf "Sprawozdanie zadan algorytmicznych" "Jan Kowalski" 123456
```

## Uruchamianie programów

### Zadanie 1 - Wyznacznik macierzy

```bash
go run ./zadanie1_wyznacznik_macierzy
```
