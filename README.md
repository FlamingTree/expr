# Expr 
[![Build Status](https://travis-ci.org/antonmedv/expr.svg?branch=master)](https://travis-ci.org/antonmedv/expr) 
[![Go Report Card](https://goreportcard.com/badge/github.com/FlamingTree/expr)](https://goreportcard.com/report/github.com/FlamingTree/expr) 
[![Code Coverage](https://scrutinizer-ci.com/g/antonmedv/expr/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/antonmedv/expr/?branch=master) 
[![GoDoc](https://godoc.org/github.com/FlamingTree/expr?status.svg)](https://godoc.org/github.com/FlamingTree/expr)

<img src="docs/images/logo-small.png" width="150" alt="expr logo" align="right">

**Expr** package provides an engine that can compile and evaluate expressions. 
An expression is a one-liner that returns a value (mostly, but not limited to, booleans).
It is designed for simplicity, speed and safety.

The purpose of the package is to allow users to use expressions inside configuration for more complex logic. 
It is a perfect candidate for the foundation of a _business rule engine_. 
The idea is to let configure things in a dynamic way without recompile of a program:

```coffeescript
# Get the special price if
user.Group in ["good_customers", "collaborator"]

# Promote article to the homepage when
len(article.Comments) > 100 and article.Category not in ["misc"]

# Send an alert when
product.Stock < 15
```

## Features

* Seamless integration with Go (no need to redefine types)
* Static typing ([example](https://godoc.org/github.com/FlamingTree/expr#example-Env)).
  ```go
  out, err := expr.Compile(`name + age`)
  // err: invalid operation + (mismatched types string and int)
  // | name + age
  // | .....^
  ```
* User-friendly error messages.
* Reasonable set of basic operators.
* Builtins `all`, `none`, `any`, `one`, `filter`, `map`.
  ```coffeescript
  all(Tweets, {.Size < 140})
  ```
* Fast ([benchmarks](https://github.com/antonmedv/golang-expression-evaluation-comparison#readme)): uses bytecode virtual machine and optimizing compiler.

## Install

```
go get github.com/FlamingTree/expr
```

<a href="https://www.patreon.com/antonmedv">
	<img src="https://c5.patreon.com/external/logo/become_a_patron_button@2x.png" width="160">
</a>

## Documentation

* See [Getting Started](docs/Getting-Started.md) page for developer documentation.
* See [Language Definition](docs/Language-Definition.md) page to learn the syntax.

## Input widget

<img src="https://user-images.githubusercontent.com/141232/61127899-3e249280-a4b9-11e9-8e88-385e2d769cce.png" align="right" alt="input widget screenshot" width="297">

Also, I have an input widget for react, vue and angular which allows editing expressions with syntax highlighting and autocomplete based on your types declaration of environment. 

[Contact me to buy ($499)](mailto:anton@medv.io)

## Examples

Executing arbitrary expressions.

```go
env := map[string]interface{}{
    "foo": 1,
    "bar": struct{Value int}{1},
}

out, err := expr.Eval("foo + bar.Value", env)
```

Static type checker with struct as environment.

```go
type Env struct {
	Foo int
	Bar *Bar
}

type Bar struct {
	Value int
}

program, err := expr.Compile("Foo + Bar.Value", expr.Env(&Env{}))

out, err := expr.Run(program, &Env{1, &Bar{2}})
```

Using env's methods as functions inside expressions.

```go
type Env struct {
	Name string
}

func (e *Env) Title() string {
	return strings.Title(e.Name)
}

program, err := expr.Compile(`"Hello " + Title()`, expr.Env(&Env{}))

out, err := expr.Run(program, &Env{"world"})
```

## Contributing

**Expr** consist of a few packages for parsing source code to AST, type checking AST, compiling to bytecode and VM for running bytecode program.

Also expr provides powerful tool [exe](cmd/exe) for debugging. It has interactive terminal debugger for our bytecode virtual machine.

<p align="center">
    <img src="docs/images/debug.gif" alt="debugger" width="605">
</p>
    

## Who is using Expr?

* [Aviasales](https://aviasales.ru) are actively using Expr for different parts of the search engine.
* [Mystery Minds](https://www.mysteryminds.com/en/) uses Expr to allow easy yet powerful customization of its matching algorithm.

## License

[MIT](LICENSE)
