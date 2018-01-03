## Gromit

[![Build Status](https://travis-ci.org/ligurio/gromit.png?branch=master)](https://travis-ci.org/ligurio/gromit) [![Go Report Card](https://goreportcard.com/badge/github.com/ligurio/gromit)](https://goreportcard.com/report/github.com/ligurio/gromit)

is a grammar fuzzer that is ideally suited for complex text and binary
grammars. Gromit uses EBNF format for grammar specification. EBNF is an
Extended Backus-Naur Form (also known as Context-Free Grammars). It is the
standard format for the specification and documentation of programming
languages. Extended BNF is defined in the [ISO/IEC 14977
standard](http://www.iso.ch/cate/d26153.html).

## How-To Use

```
~$ cat ebnf/palindrome.ebnf
palindrome = "a" | "b" | "aa" | "bb" | "a" palindrome "a" | "b" palindrome "b" .
~$ ./gromit -file ebnf/palindrome.ebnf -start palindrome
abba
~$ ./gromit -file ebnf/palindrome.ebnf -start palindrome
baab
```

See [screencast]().
