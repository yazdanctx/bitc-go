package compress

import "time"

type FormatID string

const (
	FormatOxipng       FormatID = "oxipng"
	FormatPngquant     FormatID = "pngquant"
	FormatWebpLossless FormatID = "webp-lossless"
	FormatWebpQ100     FormatID = "webp-q100"
	FormatAvifQ0       FormatID = "avif-q0"
	FormatAvifQ8       FormatID = "avif-q8"
	FormatAvifQ12      FormatID = "avif-q12"
)

type FormatOption struct {
	ID          FormatID
	Label       string
	Enabled     bool
	Extension   string
	Compression func(input, output string) error
}

type ImageFile struct {
	Path     string
	Name     string
	OrigSize int64
}

type CompressJob struct {
	Image  ImageFile
	Format FormatID
}

type CompressResult struct {
	Image        ImageFile
	Format       FormatID
	OutputPath   string
	OriginalSize int64
	CompressedSize int64
	Savings      float64
	Duration     time.Duration
	Err          error
}

type CompressionSummary struct {
	Results       []CompressResult
	TotalOriginal int64
	TotalSaved    int64
	BestFormat    FormatID
	BestFormatAvg float64
	Timestamp     string
}

func AllFormats() []FormatOption {
	return []FormatOption{
		{ID: FormatOxipng, Label: "PNG (oxipng)", Enabled: true, Extension: "png"},
		{ID: FormatPngquant, Label: "PNG (pngquant)", Enabled: true, Extension: "png"},
		{ID: FormatWebpLossless, Label: "WebP lossless", Enabled: true, Extension: "webp"},
		{ID: FormatWebpQ100, Label: "WebP q100", Enabled: true, Extension: "webp"},
		{ID: FormatAvifQ0, Label: "AVIF q0", Enabled: true, Extension: "avif"},
		{ID: FormatAvifQ8, Label: "AVIF q8", Enabled: true, Extension: "avif"},
		{ID: FormatAvifQ12, Label: "AVIF q12", Enabled: true, Extension: "avif"},
	}
}

func BestResult(results []CompressResult) *CompressResult {
	var best *CompressResult
	for i := range results {
		r := &results[i]
		if r.Err != nil {
			continue
		}
		if best == nil || r.CompressedSize < best.CompressedSize {
			best = r
		}
	}
	return best
}
