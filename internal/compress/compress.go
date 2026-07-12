package compress

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
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
	bw bool,
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

	results <- ProgressMsg{Total: totalJobs, Done: 0, Current: "starting..."}

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
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				outPath := OutputPath(job.Image.Path, job.Format, outDir)
				compressor := CompressorForFormat(job.Format)

				start := time.Now()

				var inputPath string
				if bw {
					tmpFile := filepath.Join(outDir, fmt.Sprintf(".tmp-%d%s", time.Now().UnixNano(), filepath.Ext(job.Image.Path)))
					if preErr := PreprocessToBW(job.Image.Path, tmpFile); preErr != nil {
						done <- CompressResult{
							Image:      job.Image,
							Format:     job.Format,
							OutputPath: outPath,
							Duration:   time.Since(start),
							Err:        fmt.Errorf("preprocess: %w", preErr),
						}
						continue
					}
					inputPath = tmpFile
				} else {
					inputPath = job.Image.Path
				}

				err := compressor(inputPath, outPath)
				if bw {
					os.Remove(inputPath)
				}
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
		wg.Wait()
		close(done)
	}()

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

	for _, r := range results {
		if r.Err != nil {
			continue
		}
		totalOrig += r.OriginalSize
		totalSaved += r.OriginalSize - r.CompressedSize
	}

	return &CompressionSummary{
		Results:       results,
		TotalOriginal: totalOrig,
		TotalSaved:    totalSaved,
		Timestamp:     time.Now().Format("2006-01-02-150405"),
	}
}


