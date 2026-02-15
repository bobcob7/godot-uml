package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bobcob7/go-uml/internal/server"
	"github.com/bobcob7/go-uml/pkg/gouml"
)

// version is set at build time via ldflags.
var version = "dev"

const (
	exitSuccess    = 0
	exitValidation = 1
	exitSystem     = 2
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(exitSystem)
	}
	switch os.Args[1] {
	case "render":
		os.Exit(cmdRender(os.Args[2:]))
	case "validate":
		os.Exit(cmdValidate(os.Args[2:]))
	case "serve":
		os.Exit(cmdServe(os.Args[2:]))
	case "version":
		fmt.Printf("go-uml %s\n", version)
		os.Exit(exitSuccess)
	case "help", "--help", "-h":
		printUsage()
		os.Exit(exitSuccess)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(exitSystem)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage: go-uml <command> [options]

Commands:
  render    Render a PlantUML file to SVG
  validate  Validate a PlantUML file
  serve     Start the HTTP server with live editor
  version   Print version information
  help      Show this help

Run 'go-uml <command> --help' for command-specific help.`)
}

func cmdRender(args []string) int {
	outputFile, inputPath := parseRenderArgs(args)
	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "Usage: go-uml render <file.puml|-> [-o output.svg]")
		return exitSystem
	}
	var input *os.File
	if inputPath == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return exitSystem
		}
		defer func() { _ = f.Close() }()
		input = f
	}
	var out *os.File
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return exitSystem
		}
		defer func() { _ = f.Close() }()
		out = f
	} else {
		out = os.Stdout
	}
	if err := gouml.Render(input, out); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		if isValidationError(err) {
			return exitValidation
		}
		return exitSystem
	}
	return exitSuccess
}

func cmdValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return exitSystem
	}
	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: go-uml validate <file.puml>")
		return exitSystem
	}
	inputPath := remaining[0]
	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return exitSystem
	}
	defer func() { _ = f.Close() }()
	errs := gouml.Validate(f)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", inputPath, e.Line, e.Column, e.Message)
		}
		return exitValidation
	}
	fmt.Println("OK")
	return exitSuccess
}

func cmdServe(args []string) int {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	port := fs.Int("port", 8080, "port to listen on")
	host := fs.String("host", "localhost", "host to bind to")
	if err := fs.Parse(args); err != nil {
		return exitSystem
	}
	cfg := server.DefaultConfig()
	cfg.Port = *port
	cfg.Host = *host
	srv := server.New(cfg)
	fmt.Fprintf(os.Stderr, "go-uml server listening on http://%s:%d\n", cfg.Host, cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return exitSystem
	}
	return exitSuccess
}

func parseRenderArgs(args []string) (outputFile, inputPath string) {
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "-o" && i+1 < len(args):
			outputFile = args[i+1]
			i++
		case args[i] == "--help" || args[i] == "-h":
			return "", ""
		case args[i] == "-" || !strings.HasPrefix(args[i], "-"):
			inputPath = args[i]
		}
	}
	return outputFile, inputPath
}

func isValidationError(err error) bool {
	return strings.Contains(err.Error(), ":")
}
