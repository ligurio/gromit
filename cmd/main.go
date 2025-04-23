package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	gromit "github.com/ligurio/gromit"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/ebnf"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

var (
	name    = flag.String("filename", "", "filename with grammar")
	action  = flag.String("action", "fuzz", "action (possible values: fuzz and dict)")
	start   = flag.String("start", "", "start string")
	seed    = flag.Int64("seed", -1, "number used to initialize a pseudorandom number generator")
	maxreps = flag.Int("maxreps", 10, "maximum number of repetitions")
	depth   = flag.Int("depth", 30, "maximum depth")
	padding = flag.String("padding", " ", "non-terminal padding characters")
	debug   = flag.Bool("debug", false, "enable verbosity")
)

func main() {

	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if *name == "" && *start == "" {
		flag.Usage()
		fmt.Println()
		log.Fatal("Filename or start string is not specified.")
		os.Exit(1)
	}

	f, err := os.Open(*name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	grammar, err := ebnf.Parse(*name, bufio.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}

	err = ebnf.Verify(grammar, *start)
	if err != nil {
		log.Fatal(err)
	}

	if *action == "dict" {
		err = gromit.Dict(os.Stdout, grammar, *start, *seed)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}

	log.Debug("Grammar was successfully verified.")
	err = gromit.Random(os.Stdout, grammar, *start, *seed, *maxreps)
	if err != nil {
		log.Fatal(err)
	}
}
