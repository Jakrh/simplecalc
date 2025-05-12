# Simple calculator

Parse and calculate numbers from strings like `2 * (12.4 / (7 + -2))`, or multi-line inputs that separated by semicolons like `x = 2; y = 5.25; z = x * (3 + -y); z`.

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
* `x = 2; y = 5.25; z = x * (3 + -y); z`

---

## Download the `golang.org/x/term` package

```bash
go mod tidy
```

## Test

```bash
go test -v ./...
```

## Run

```bash
go run .
```

## Usage

After starting the calculator, you can type expressions directly into the interactive terminal, assign variables, or use built-in commands.

For example:

```
Enter an expression (or 'exit' to quit):
>>> 2 * (12.4 / (7 + -2))
4.96
>>> x = 10
>>> x * 3
30
```

Input `help` to see a list of available commands, supported operators, and syntax information.
