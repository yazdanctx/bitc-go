package compress

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func CompressAvifQ12(input, output string) error {
	cmd := exec.Command("avifenc", "--min", "12", "--max", "12", input, output)
	return cmd.Run()
}

func OutputPath(input string, format FormatID, outDir string) string {
	ext := extensionForFormat(format)
	base := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
	return filepath.Join(outDir, fmt.Sprintf("%s-%s.%s", base, string(format), ext))
}

func extensionForFormat(f FormatID) string {
	switch f {
	case FormatAvifQ12:
		return "avif"
	default:
		return "bin"
	}
}

func CompressorForFormat(f FormatID) func(string, string) error {
	switch f {
	case FormatAvifQ12:
		return CompressAvifQ12
	default:
		return func(_, _ string) error {
			return fmt.Errorf("unknown format: %s", f)
		}
	}
}
