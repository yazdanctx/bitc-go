package compress

import "time"

type FormatID string

const (
	FormatAvifQ12 FormatID = "avif-q12"
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
	Image          ImageFile
	Format         FormatID
	OutputPath     string
	OriginalSize   int64
	CompressedSize int64
	Savings        float64
	Duration       time.Duration
	Err            error
}

type CompressionSummary struct {
	Results       []CompressResult
	TotalOriginal int64
	TotalSaved    int64
	Timestamp     string
}

func AllFormats() []FormatOption {
	return []FormatOption{
		{ID: FormatAvifQ12, Label: "AVIF q12", Enabled: true, Extension: "avif"},
	}
}
