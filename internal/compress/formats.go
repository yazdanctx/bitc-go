package compress

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func CompressOxipng(input, output string) error {
	cmd := exec.Command("oxipng", "-o", "max", "--out", output, input)
	return cmd.Run()
}

func CompressPngquant(input, output string) error {
	cmd := exec.Command("pngquant", "--force", "--output", output, "--quality=100", "--speed=1", input)
	return cmd.Run()
}

func CompressWebpLossless(input, output string) error {
	cmd := exec.Command("cwebp", "-lossless", input, "-o", output)
	return cmd.Run()
}

func CompressWebpQ100(input, output string) error {
	cmd := exec.Command("cwebp", "-q", "100", input, "-o", output)
	return cmd.Run()
}

func CompressAvifQ0(input, output string) error {
	cmd := exec.Command("avifenc", "--min", "0", "--max", "0", input, output)
	return cmd.Run()
}

func CompressAvifQ8(input, output string) error {
	cmd := exec.Command("avifenc", "--min", "8", "--max", "8", input, output)
	return cmd.Run()
}

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
	case FormatOxipng, FormatPngquant:
		return "png"
	case FormatWebpLossless, FormatWebpQ100:
		return "webp"
	case FormatAvifQ0, FormatAvifQ8, FormatAvifQ12:
		return "avif"
	default:
		return "bin"
	}
}

func CompressorForFormat(f FormatID) func(string, string) error {
	switch f {
	case FormatOxipng:
		return CompressOxipng
	case FormatPngquant:
		return CompressPngquant
	case FormatWebpLossless:
		return CompressWebpLossless
	case FormatWebpQ100:
		return CompressWebpQ100
	case FormatAvifQ0:
		return CompressAvifQ0
	case FormatAvifQ8:
		return CompressAvifQ8
	case FormatAvifQ12:
		return CompressAvifQ12
	default:
		return func(_, _ string) error {
			return fmt.Errorf("unknown format: %s", f)
		}
	}
}
