// SPDX-License-Identifier: MIT

package gromit

import (
	"errors"
	"golang.org/x/exp/ebnf"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	rand "math/rand"
)

var ErrStartNotFound = errors.New("Start production not found")
var ErrBadRange = errors.New("Bad range")

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func Random(dst io.Writer, grammar ebnf.Grammar, start string, seed int64, maxreps int) error {
	production, err := grammar[start]
	if !err {
		return ErrStartNotFound
	}

	if seed == -1 {
		rand.Seed(time.Now().UTC().UnixNano())
	}

	return random(dst, grammar, production.Expr, 0, maxreps)
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

func random(dst io.Writer, grammar ebnf.Grammar, expr ebnf.Expression, depth int, maxreps int) error {

	var maxdepth int // FIXME
	maxdepth = 100

	if expr == nil {
		return nil
	}

	switch expr.(type) {
	case ebnf.Alternative:
		log.Debug("Expression is alternative")
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
		err := random(dst, grammar, exprs[rand.Intn(len(exprs))], depth+1, maxreps)
		if err != nil {
			return err
		}

	case *ebnf.Group:
		log.Debug("Expression is group")
		gr := expr.(*ebnf.Group)
		err := random(dst, grammar, gr.Body, depth+1, maxreps)
		if err != nil {
			return err
		}

	case *ebnf.Name:
		log.Debug("Expression is name")
		name := expr.(*ebnf.Name)
		p := !IsTerminal(expr)
		if p {
			pad(dst)
		}
		err := random(dst, grammar, grammar[name.String], depth+1, maxreps)
		if err != nil {
			return err
		}
		if p {
			pad(dst)
		}

	case *ebnf.Option:
		log.Debug("Expression is option")
		opt := expr.(*ebnf.Option)
		if depth > maxdepth && !IsTerminal(opt.Body) {
			log.Debug("non-terminal omitted due to having exceeded recursion depth limit\n")
		} else if PickBool() {
			err := random(dst, grammar, opt.Body, depth+1, maxreps)
			if err != nil {
				return err
			}
		}

	case *ebnf.Production:
		log.Debug("Expression is production")
		prod := expr.(*ebnf.Production)
		err := random(dst, grammar, prod.Expr, depth+1, maxreps)
		if err != nil {
			return err
		}

	case *ebnf.Range:
		log.Debug("Expression is range")
		rang := expr.(*ebnf.Range)
		ch, err := PickString(rang.Begin.String, rang.End.String)
		if err != nil {
			return err
		}
		if _, err := io.WriteString(dst, ch); err != nil {
			return err
		}

	case *ebnf.Repetition:
		log.Debug("Expression is repetition")
		rep := expr.(*ebnf.Repetition)
		if depth > maxdepth && !IsTerminal(rep.Body) {
			log.Debug("Repetition omitted\n")
		} else {
			reps := rand.Intn(maxreps + 1)
			log.Debug(reps, " times")
			for i := 0; i < reps; i++ {
				err := random(dst, grammar, rep.Body, depth+1, maxreps)
				if err != nil {
					return err
				}
			}
		}

	case ebnf.Sequence:
		log.Debug("Expression is sequence")
		seq := expr.(ebnf.Sequence)
		for _, e := range seq {
			err := random(dst, grammar, e, depth+1, maxreps)
			if err != nil {
				return err
			}
		}

	case *ebnf.Token:
		log.Debug("Expression is token")
		tok := expr.(*ebnf.Token)
		if _, err := io.WriteString(dst, tok.String+" "); err != nil {
			return err
		}

	default:
		log.Fatal("Bad expression", expr)
	}

	return nil
}

func Dict(dst io.Writer, grammar ebnf.Grammar, start string, seed int64) error {

	production, err := grammar[start]
	if !err {
		return ErrStartNotFound
	}

	if seed == -1 {
		rand.Seed(time.Now().UTC().UnixNano())
	}

	return random1(dst, grammar, production.Expr, 0)
}

func random1(dst io.Writer, grammar ebnf.Grammar, expr ebnf.Expression, depth int) error {

	switch expr.(type) {
	case ebnf.Alternative:
		log.Debug("Expression is alternative")
		alt := expr.(ebnf.Alternative)
		var exprs []ebnf.Expression
		exprs = findTerminals(alt)
		if len(exprs) == 0 {
			exprs = alt
		}
		for _, e := range exprs {
			err := random1(dst, grammar, e, depth+1)
			if err != nil {
				return err
			}
		}

	case *ebnf.Group:
		log.Debug("Expression is group")
		gr := expr.(*ebnf.Group)
		err := random1(dst, grammar, gr.Body, depth+1)
		if err != nil {
			return err
		}

	case *ebnf.Name:
		log.Debug("Expression is name")
		name := expr.(*ebnf.Name)
		p := !IsTerminal(expr)
		if p {
			pad(dst)
		}
		err := random1(dst, grammar, grammar[name.String], depth+1)
		if err != nil {
			return err
		}
		if p {
			pad(dst)
		}

	case *ebnf.Option:
		log.Debug("Expression is option")
		opt := expr.(*ebnf.Option)
		err := random1(dst, grammar, opt.Body, depth+1)
		if err != nil {
			return err
		}

	case *ebnf.Production:
		log.Debug("Expression is production")
		prod := expr.(*ebnf.Production)
		err := random1(dst, grammar, prod.Expr, depth+1)
		if err != nil {
			return err
		}

	case *ebnf.Range:
		log.Debug("Expression is range")

	case *ebnf.Repetition:
		log.Debug("Expression is repetition")

	case ebnf.Sequence:
		log.Debug("Expression is sequence")
		seq := expr.(ebnf.Sequence)
		for _, e := range seq {
			err := random1(dst, grammar, e, depth+1)
			if err != nil {
				return err
			}
		}

	case *ebnf.Token:
		log.Debug("Expression is token")
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
