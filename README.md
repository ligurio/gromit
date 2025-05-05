## Gromit

[![Testing](https://github.com/ligurio/gromit/actions/workflows/test.yml/badge.svg)](https://github.com/ligurio/gromit/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/ligurio/gromit)](https://goreportcard.com/report/github.com/ligurio/gromit)

is a random text generator based on context-free grammars; it uses
an EBNF for grammar definitions. EBNF is an Extended Backus-Naur
Form. It is the standard format for the specification and
documentation of programming languages; it is defined in the
ISO/IEC 14977 standard.

The input is text satisfying the following grammar (represented
itself in EBNF):

```
Production  = name "=" [ Expression ] "." .
Expression  = Alternative { "|" Alternative } .
Alternative = Term { Term } .
Term        = name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .
```

## Usage

```
~$ go build -o gromit -v cmd/main.go
~$ ./gromit -filename ebnf/palindrome.ebnf -start palindrome
khbhk
```

## See also

- [grammarinator](https://github.com/renatahodovan/grammarinator)
  is an ANTLR v4 grammar-based test generator (Python).
- [bnfgen](https://baturin.org/tools/bnfgen/) is a random text
  generator based on context-free grammars, it uses a DSL for
  grammar definitions that is similar to the familiar BNF,
  with two extensions: weighted random selection and deterministic
  repetition (OCaml).
