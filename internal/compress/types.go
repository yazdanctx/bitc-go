package compress

import (
	"fmt"
	"time"
)

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

func FormatSize(bytes int64) string {
	switch {
	case bytes >= 1<<20:
		return formatFloat(float64(bytes)/(1<<20)) + " MB"
	case bytes >= 1<<10:
		return formatFloat(float64(bytes)/(1<<10)) + " KB"
	default:
		return formatFloat(float64(bytes)) + " B"
	}
}

func formatFloat(f float64) string {
	s := "%.1f"
	if f >= 100 {
		s = "%.0f"
	}
	return fmt.Sprintf(s, f)
}

func AllFormats() []FormatOption {
	return []FormatOption{
		{ID: FormatAvifQ12, Label: "AVIF q12", Enabled: true, Extension: "avif"},
	}
}
