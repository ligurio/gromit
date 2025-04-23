## Gromit

[![Testing](https://github.com/ligurio/gromit/actions/workflows/test.yml/badge.svg)](https://github.com/ligurio/gromit/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/ligurio/gromit)](https://goreportcard.com/report/github.com/ligurio/gromit)

is a grammar fuzzer that is ideally suited for complex text and binary
grammars. Gromit uses EBNF format for grammar specification. EBNF is an
Extended Backus-Naur Form (also known as Context-Free Grammars). It is the
standard format for the specification and documentation of programming
languages. Extended BNF is defined in the [ISO/IEC 14977
standard](http://www.iso.ch/cate/d26153.html).

## How-To Use

```
~$ ./gromit -file ebnf/palindrome.ebnf -start palindrome
khbhk
```

See [screencast](https://asciinema.org/a/155319).
