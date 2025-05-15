[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Jakrh/simplecalc)

# Simple calculator

Parse and calculate numbers from strings like `2 * (12.4 / (7 + -2))`, or multi-line inputs that separated by semicolons like `x = 2; y = 5.25; z = x * (3 + -y); z`.

Based on Pratt parsing and [jdvillal/parser](https://github.com/jdvillal/parser/).

## Supported operators:

* `+`
* `-`
* `*`
* `/`
* `**`
* `(` and `)`

## Supported expressions like:

* `-.25 + 2`
* `-1 + -2 * -3`
* `-5 * (2 + -.3)`
* `2 * (12.4 / (7 + -2))`
* `-34 * (2 + -.23)`
* `4 * (272 + 6) - 324 / 8`
* `2 ** 10`
* `x = 2; y = 5.25; z = x * (3 + -y); z`
* `x = 1.6; y = .25; -((2.5 * x) ** 6) ** y / .5 ** 3`

---

## How it works

The calculator processes an input in three stages:

1. **Lex**: splits the input string into tokens.
2. **Parse**: builds an `Expression` AST (Abstract Syntax Tree) from tokens with binding power of each operator.
3. **Evaluate**: traverses the AST to compute the final result.

**Example**

0. Input

    `-16 ** (.25 + (7 - -2) / 4)`

1. Input -> Lex -> Tokens

    `["-", "16", "**", "(", ".25", "+", "(", "7", "-", "-", "2", ")", "/", "4", ")"]`

2. Tokens -> Parse -> Expression (binding powers: `**` > `-` negative sign,  `/` > `+` and `-`)

    `(- 0 (** 16 (+ 0.25 (/ (- 7 (- 0 2)) 4))))`

3. Evaluate -> result

    `-1024`

---

## Run the calculator

### Download the `golang.org/x/term` package

```bash
go mod tidy
```

### Test

```bash
go test -v ./...
```

### Run

```bash
go run .
```

### Run with debug mode to observe each stage of processing

```bash
DEBUG=1 go run .
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
