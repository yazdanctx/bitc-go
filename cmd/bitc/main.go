package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/yazdanctx/bitc-go/internal/compress"
	"github.com/yazdanctx/bitc-go/internal/scanner"
	"github.com/yazdanctx/bitc-go/internal/version"
)

func main() {
	var (
		showVersion bool
		outputDir   string
		bw          bool
	)

	flag.BoolVar(&showVersion, "version", false, "Print version")
	flag.StringVar(&outputDir, "output", "", "Override output directory")
	flag.BoolVar(&bw, "bw", false, "Convert images to 2-color black & white before compressing")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "bitc — image compression tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  bitc <directory> [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		fmt.Printf("bitc %s\n", version.Version)
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

	fmt.Printf("bitc %s\n\n", version.Version)
	fmt.Printf("Scanning %s ...\n", dir)

	images, err := scanner.ScanDirectory(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	if len(images) == 0 {
		fmt.Println("No images found.")
		os.Exit(0)
	}

	fmt.Printf("Found %d image(s)\n\n", len(images))

	formats := compress.AllFormats()
	resultsCh := make(chan compress.ProgressMsg, 100)

	go compress.RunCompression(images, formats, outputDir, bw, resultsCh)

	var summary *compress.CompressionSummary
	for msg := range resultsCh {
		if msg.Finished {
			summary = msg.Summary
			break
		}
		if msg.Result != nil {
			r := msg.Result
			if r.Err != nil {
				fmt.Printf("  ✗ %s — %v\n", r.Image.Name, r.Err)
			} else {
				fmt.Printf("  %s  %s → %s (%.1f%% saved)\n", r.Image.Name, compress.FormatSize(r.OriginalSize), compress.FormatSize(r.CompressedSize), r.Savings)
			}
		}
	}

	if summary != nil && len(summary.Results) > 0 {
		notifyDone(outputDir)
	}
}

func notifyDone(dir string) {
	msg := fmt.Sprintf("bitc: done — saved to %s", dir)
	exec.Command("osascript", "-e", fmt.Sprintf(`display notification %q with title "bitc"`, msg)).Run()
	fmt.Printf("\n%s\n", msg)
}
