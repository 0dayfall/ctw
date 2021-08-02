package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// We declare a subcommand using the `NewFlagSet`
	// function, and proceed to define new flags specific
	// for this subcommand.
	ruleCmd := flag.NewFlagSet("rules", flag.ExitOnError)
	ruleAdd := ruleCmd.String("add", "", "")
	ruleDelete := ruleCmd.String("delete", "", "")
	ruleList := ruleCmd.String("list", "", "")

	// For a different subcommand we can define different
	// supported flags.
	streamCmd := flag.NewFlagSet("stream", flag.ExitOnError)

	// The subcommand is expected as the first argument
	// to the program.
	if len(os.Args) < 2 {
		fmt.Println("expected subcommands")
		os.Exit(1)
	}

	// Check which subcommand is invoked.
	switch os.Args[1] {

	// For every subcommand, we parse its own flags and
	// have access to trailing positional arguments.
	case "rules":
		err := ruleCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
		if ruleCmd.Parsed() {
			if *ruleAdd == "" {
				fmt.Println("Please supply the rule using -add option.")
				return
			}
			if *ruleDelete == "" {
				fmt.Println("Please supply the rule using -delete option.")
				return
			}
			if *ruleList == "" {

				return
			}
			fmt.Printf("You have added: %q\n", *ruleAdd)
		}
	case "stream":
		err := streamCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
		//request.CreateGetRequest(stream.CreateUrl())
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}
