package compress

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type ProgressMsg struct {
	Current  string
	Done     int
	Total    int
	Result   *CompressResult
	Finished bool
	Summary  *CompressionSummary
}

func RunCompression(
	images []ImageFile,
	formats []FormatOption,
	outDir string,
	results chan<- ProgressMsg,
) {
	var enabledFormats []FormatID
	for _, f := range formats {
		if f.Enabled {
			enabledFormats = append(enabledFormats, f.ID)
		}
	}

	totalJobs := len(images) * len(enabledFormats)
	if totalJobs == 0 {
		results <- ProgressMsg{Finished: true, Summary: &CompressionSummary{}}
		return
	}

	os.MkdirAll(outDir, 0755)

	jobs := make(chan CompressJob, totalJobs)
	allResults := make([]CompressResult, 0, totalJobs)

	for _, img := range images {
		for _, fid := range enabledFormats {
			jobs <- CompressJob{Image: img, Format: fid}
		}
	}
	close(jobs)

	workers := runtime.NumCPU()
	if workers > len(images) {
		workers = len(images)
	}
	if workers < 1 {
		workers = 1
	}

	done := make(chan CompressResult, totalJobs)

	for i := 0; i < workers; i++ {
		go func() {
			for job := range jobs {
				outPath := OutputPath(job.Image.Path, job.Format, outDir)
				compressor := CompressorForFormat(job.Format)

				start := time.Now()
				err := compressor(job.Image.Path, outPath)
				duration := time.Since(start)

				var compSize int64
				if err == nil {
					info, statErr := os.Stat(outPath)
					if statErr == nil {
						compSize = info.Size()
					}
				}

				savings := 0.0
				if job.Image.OrigSize > 0 && err == nil {
					savings = float64(job.Image.OrigSize-compSize) / float64(job.Image.OrigSize) * 100
				}

				result := CompressResult{
					Image:          job.Image,
					Format:         job.Format,
					OutputPath:     outPath,
					OriginalSize:   job.Image.OrigSize,
					CompressedSize: compSize,
					Savings:        savings,
					Duration:       duration,
					Err:            err,
				}
				done <- result
			}
		}()
	}

	go func() {
		for result := range done {
			allResults = append(allResults, result)
			results <- ProgressMsg{
				Current: result.Image.Name,
				Done:    len(allResults),
				Total:   totalJobs,
				Result:  &result,
			}
		}

		summary := buildSummary(allResults, images)
		results <- ProgressMsg{Finished: true, Summary: summary}
	}()
}

func buildSummary(results []CompressResult, images []ImageFile) *CompressionSummary {
	var totalOrig, totalSaved int64
	formatStats := make(map[FormatID]struct {
		total   int64
		saved   int64
		count   int
	})

	for _, r := range results {
		if r.Err != nil {
			continue
		}
		totalOrig += r.OriginalSize
		totalSaved += r.OriginalSize - r.CompressedSize

		fs := formatStats[r.Format]
		fs.total += r.OriginalSize
		fs.saved += r.OriginalSize - r.CompressedSize
		fs.count++
		formatStats[r.Format] = fs
	}

	bestFormat := FormatID("")
	bestAvg := 0.0
	for fid, fs := range formatStats {
		if fs.count == 0 {
			continue
		}
		avg := float64(fs.saved) / float64(fs.total) * 100
		if avg > bestAvg {
			bestAvg = avg
			bestFormat = fid
		}
	}

	return &CompressionSummary{
		Results:       results,
		TotalOriginal: totalOrig,
		TotalSaved:    totalSaved,
		BestFormat:    bestFormat,
		BestFormatAvg: bestAvg,
		Timestamp:     time.Now().Format("2006-01-02-150405"),
	}
}

func CopyBestResults(summary *CompressionSummary, outputDir string) error {
	os.MkdirAll(outputDir, 0755)

	seen := make(map[string]bool)
	for _, r := range summary.Results {
		if r.Err != nil {
			continue
		}
		key := r.Image.Path
		if seen[key] {
			continue
		}

		best := BestResultForImage(summary.Results, r.Image.Path)
		if best == nil {
			continue
		}

		src, err := os.Open(best.OutputPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", best.OutputPath, err)
		}
		defer src.Close()

		dstPath := filepath.Join(outputDir, r.Image.Name)
		dst, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("create %s: %w", dstPath, err)
		}
		defer dst.Close()

		_, err = dst.ReadFrom(src)
		if err != nil {
			return fmt.Errorf("copy to %s: %w", dstPath, err)
		}

		seen[key] = true
	}
	return nil
}

func BestResultForImage(results []CompressResult, imagePath string) *CompressResult {
	var best *CompressResult
	for i := range results {
		r := &results[i]
		if r.Err != nil {
			continue
		}
		if r.Image.Path != imagePath {
			continue
		}
		if best == nil || r.CompressedSize < best.CompressedSize {
			best = r
		}
	}
	return best
}
