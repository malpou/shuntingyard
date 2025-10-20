# shuntingyard

A Go package implementing the Shunting Yard algorithm with floating-point support for parsing and evaluating mathematical expressions.

## Features

- Full floating-point number support (float64)
- Operators: `+`, `-`, `*`, `/`
- Proper operator precedence and left-associativity
- Parentheses support
- Comprehensive error handling
- Zero dependencies, thread-safe

## Installation

```bash
go get github.com/malpou/shuntingyard
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/malpou/shuntingyard"
)

func main() {
    expression := "2 + 3 * 4"

    tokens, err := shuntingyard.Scan(expression)
    if err != nil {
        panic(err)
    }

    postfix, err := shuntingyard.Parse(tokens)
    if err != nil {
        panic(err)
    }

    result, err := shuntingyard.Evaluate(postfix)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s = %v\n", expression, result)
}
```

## API

### `Scan(expression string) ([]string, error)`
Tokenizes a mathematical expression into tokens. Supports integers, floats, operators (`+`, `-`, `*`, `/`), and parentheses.

### `Parse(tokens []string) ([]string, error)`
Converts infix notation to postfix (RPN) using the Shunting Yard algorithm. Handles operator precedence and parentheses.

### `Evaluate(postfixTokens []string) (float64, error)`
Evaluates a postfix expression and returns the float64 result.

## Testing

```bash
go test -v
go test -bench=.
```

## License

MIT
