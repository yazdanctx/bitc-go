package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yazdanctx/bitc-go/internal/compress"
	"github.com/yazdanctx/bitc-go/internal/tui"
)

var version = "dev"

func main() {
	var (
		showVersion bool
		outputDir   string
	)

	flag.BoolVar(&showVersion, "version", false, "Print version")
	flag.StringVar(&outputDir, "output", "", "Override output directory")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "bitc-go — image compression tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  bitc-go <directory> [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		fmt.Printf("bitc-go %s\n", version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	dir := flag.Arg(0)
	dir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s is not a directory\n", dir)
		os.Exit(1)
	}

	if outputDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		outputDir = filepath.Join(home, "Downloads", "compressed-"+time.Now().Format("2006-01-02-150405"))
	}

	p := tea.NewProgram(tui.InitialModel(dir, outputDir), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(tui.ModelAccessor); ok && m.GetSummary() != nil {
		if err := compress.CopyBestResults(m.GetSummary(), outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving results: %v\n", err)
			os.Exit(1)
		}
	}
}
