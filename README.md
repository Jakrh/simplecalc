# Simple calculator

Parse and calculate numbers from strings like `2 * (12.4 / (7 + -2))`.

Based on Pratt parsing and [jdvillal/parser](https://github.com/jdvillal/parser/).

## Supported operators:

* `+`
* `-`
* `*`
* `/`
* `(` and `)`

## Supported expressions like:

* `-.25 + 2`
* `-1 + -2 * -3`
* `-5 * (2 + -.3)`
* `2 * (12.4 / (7 + -2))`
* `-34 * (2 + -.23)`
* `4 * (272 + 6) - 324 / 8`

---

## Test

```bash
go test -v ./...
```

## Run

```bash
go run .
```
