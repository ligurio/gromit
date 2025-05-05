// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	rand "math/rand"

	gromit "github.com/ligurio/gromit"
	"golang.org/x/exp/ebnf"
)

var (
	name    = flag.String("filename", "", "filename with grammar")
	action  = flag.String("action", "fuzz", "action (possible values: fuzz and dict)")
	start   = flag.String("start", "", "start string")
	seed    = flag.Int64("seed", 0, "random seed; if 0, seed is generated (default)")
	maxreps = flag.Int("maxreps", 10, "maximum number of repetitions")
	depth   = flag.Int("depth", 30, "maximum depth")
	padding = flag.String("padding", " ", "non-terminal padding characters")
)

func main() {
	flag.Parse()

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	source := rand.NewSource(*seed)
	rng := rand.New(source)

	if *name == "" && *start == "" {
		flag.Usage()
		fmt.Println("Filename or start string is not specified.")
		os.Exit(1)
	}

	f, err := os.Open(*name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	grammar, err := ebnf.Parse(*name, bufio.NewReader(f))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ebnf.Verify(grammar, *start)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *action == "dict" {
		err = gromit.Dict(os.Stdout, grammar, *start, rng)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(1)
	}

	err = gromit.Random(os.Stdout, grammar, *start, rng, *maxreps)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
