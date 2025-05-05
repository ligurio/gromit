// SPDX-License-Identifier: MIT

package gromit

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/exp/ebnf"

	rand "math/rand"
)

var ErrStartNotFound = errors.New("Start production not found")
var ErrBadRange = errors.New("Bad range")

func Random(dst io.Writer, grammar ebnf.Grammar, start string, rng *rand.Rand, maxreps int) error {
	production, err := grammar[start]
	if !err {
		return ErrStartNotFound
	}

	return random(dst, grammar, production.Expr, 0, maxreps, rng)
}

func IsTerminal(expr ebnf.Expression) bool {
	switch expr.(type) {
	case *ebnf.Name:
		name := expr.(*ebnf.Name)
		return !IsCapital(name.String)
	case *ebnf.Range:
		return true
	case *ebnf.Token:
		return true
	default:
		return false
	}
}

func findTerminals(exprs []ebnf.Expression) []ebnf.Expression {
	r := make([]ebnf.Expression, 0, len(exprs))
	for _, expr := range exprs {
		if IsTerminal(expr) {
			r = append(r, expr)
		}
	}
	return r
}

func random(dst io.Writer, grammar ebnf.Grammar, expr ebnf.Expression, depth int, maxreps int, rng *rand.Rand) error {

	var maxdepth int // FIXME
	maxdepth = 100

	if expr == nil {
		return nil
	}

	switch expr.(type) {
	case ebnf.Alternative:
		alt := expr.(ebnf.Alternative)
		var exprs []ebnf.Expression
		if depth > maxdepth {
			exprs = findTerminals(alt)
			if len(exprs) == 0 {
				exprs = alt
			}
		} else {
			exprs = alt
		}
		err := random(dst, grammar, exprs[rng.Intn(len(exprs))], depth+1, maxreps, rng)
		if err != nil {
			return err
		}

	case *ebnf.Group:
		gr := expr.(*ebnf.Group)
		err := random(dst, grammar, gr.Body, depth+1, maxreps, rng)
		if err != nil {
			return err
		}

	case *ebnf.Name:
		name := expr.(*ebnf.Name)
		p := !IsTerminal(expr)
		if p {
			pad(dst)
		}
		err := random(dst, grammar, grammar[name.String], depth+1, maxreps, rng)
		if err != nil {
			return err
		}
		if p {
			pad(dst)
		}

	case *ebnf.Option:
		opt := expr.(*ebnf.Option)
		if depth > maxdepth && !IsTerminal(opt.Body) {
			fmt.Println("non-terminal omitted due to having exceeded recursion depth limit")
		} else if PickBool() {
			err := random(dst, grammar, opt.Body, depth+1, maxreps, rng)
			if err != nil {
				return err
			}
		}

	case *ebnf.Production:
		prod := expr.(*ebnf.Production)
		err := random(dst, grammar, prod.Expr, depth+1, maxreps, rng)
		if err != nil {
			return err
		}

	case *ebnf.Range:
		rang := expr.(*ebnf.Range)
		ch, err := PickString(rang.Begin.String, rang.End.String)
		if err != nil {
			return err
		}
		if _, err := io.WriteString(dst, ch); err != nil {
			return err
		}

	case *ebnf.Repetition:
		rep := expr.(*ebnf.Repetition)
		if depth > maxdepth && !IsTerminal(rep.Body) {
			fmt.Println("Repetition omitted")
		} else {
			reps := rng.Intn(maxreps + 1)
			for i := 0; i < reps; i++ {
				err := random(dst, grammar, rep.Body, depth+1, maxreps, rng)
				if err != nil {
					return err
				}
			}
		}

	case ebnf.Sequence:
		seq := expr.(ebnf.Sequence)
		for _, e := range seq {
			err := random(dst, grammar, e, depth+1, maxreps, rng)
			if err != nil {
				return err
			}
		}

	case *ebnf.Token:
		tok := expr.(*ebnf.Token)
		if _, err := io.WriteString(dst, tok.String+" "); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Bad expression %g", expr)
	}

	return nil
}

func Dict(dst io.Writer, grammar ebnf.Grammar, start string, rng *rand.Rand) error {

	production, err := grammar[start]
	if !err {
		return ErrStartNotFound
	}

	return random1(dst, grammar, production.Expr, 0, rng)
}

func random1(dst io.Writer, grammar ebnf.Grammar, expr ebnf.Expression, depth int, rng *rand.Rand) error {

	switch expr.(type) {
	case ebnf.Alternative:
		alt := expr.(ebnf.Alternative)
		var exprs []ebnf.Expression
		exprs = findTerminals(alt)
		if len(exprs) == 0 {
			exprs = alt
		}
		for _, e := range exprs {
			err := random1(dst, grammar, e, depth+1, rng)
			if err != nil {
				return err
			}
		}

	case *ebnf.Group:
		gr := expr.(*ebnf.Group)
		err := random1(dst, grammar, gr.Body, depth+1, rng)
		if err != nil {
			return err
		}

	case *ebnf.Name:
		name := expr.(*ebnf.Name)
		p := !IsTerminal(expr)
		if p {
			pad(dst)
		}
		err := random1(dst, grammar, grammar[name.String], depth+1, rng)
		if err != nil {
			return err
		}
		if p {
			pad(dst)
		}

	case *ebnf.Option:
		opt := expr.(*ebnf.Option)
		err := random1(dst, grammar, opt.Body, depth+1, rng)
		if err != nil {
			return err
		}

	case *ebnf.Production:
		prod := expr.(*ebnf.Production)
		err := random1(dst, grammar, prod.Expr, depth+1, rng)
		if err != nil {
			return err
		}

	case *ebnf.Range:
		// do nothing

	case *ebnf.Repetition:
		// do nothing

	case ebnf.Sequence:
		seq := expr.(ebnf.Sequence)
		for _, e := range seq {
			err := random1(dst, grammar, e, depth+1, rng)
			if err != nil {
				return err
			}
		}

	case *ebnf.Token:
		tok := expr.(*ebnf.Token)
		if _, err := io.WriteString(dst, dictline(tok.String)+"\n"); err != nil {
			return err
		}

	default:
		log.Fatal("Bad expression", expr)
	}

	return nil
}

func dictline(tok string) string {
	var replacer = strings.NewReplacer(" ", "_", "-", "_")
	str := replacer.Replace(tok)
	return "KEYWORD_" + strings.ToUpper(str) + "=\"" + tok + "\""

}
