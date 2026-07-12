package tui

import (
	"fmt"
	"strings"
	"time"
)

func (m model) viewScanning() string {
	var b strings.Builder
	b.WriteString("\n  " + m.spinner.View() + "  Scanning directory for images...\n")
	return b.String()
}

func (m model) viewCompressing() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" bitc-go ") + "\n\n")

	elapsed := time.Since(m.startTime).Truncate(time.Second)

	if m.progress.Total > 0 {
		pct := float64(m.progress.Done) / float64(m.progress.Total) * 100
		bar := renderProgressBar(pct, 40)
		b.WriteString(fmt.Sprintf("  %s  %d/%d files\n\n", bar, m.progress.Done, m.progress.Total))
	}

	if m.currentFile != "" {
		b.WriteString(fmt.Sprintf("  %s %s", m.spinner.View(), fileDefaultStyle.Render(m.currentFile)))
		if m.currentSize > 0 {
			b.WriteString(subtitleStyle.Render(fmt.Sprintf("  (%s)", FormatSize(m.currentSize))))
		}
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("\n  %s elapsed\n", subtitleStyle.Render(elapsed.String())))
	b.WriteString("\n" + subtitleStyle.Render("  q quit"))

	return b.String()
}

func (m model) viewResults() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" bitc-go ") + "\n\n")

	if m.summary == nil {
		return b.String()
	}

	b.WriteString(fmt.Sprintf("  %-30s %10s %10s %10s\n", "File", "Original", "Compressed", "Savings"))
	b.WriteString(strings.Repeat("─", 64) + "\n")

	for _, r := range m.summary.Results {
		if r.Err != nil {
			b.WriteString(fmt.Sprintf("  %-30s %10s %10s %s\n",
				r.Image.Name,
				FormatSize(r.OriginalSize),
				errorStyle.Render("error"),
				errorStyle.Render(r.Err.Error()),
			))
			continue
		}
		b.WriteString(fmt.Sprintf("  %-30s %10s %10s %s\n",
			r.Image.Name,
			FormatSize(r.OriginalSize),
			FormatSize(r.CompressedSize),
			FormatSavings(r.Savings),
		))
	}

	b.WriteString(strings.Repeat("─", 64) + "\n")
	totalPct := 0.0
	if m.summary.TotalOriginal > 0 {
		totalPct = float64(m.summary.TotalSaved) / float64(m.summary.TotalOriginal) * 100
	}
	b.WriteString(summaryStyle.Render(fmt.Sprintf("  Total: %s → %s saved (%.1f%%)",
		FormatSize(m.summary.TotalOriginal),
		FormatSize(m.summary.TotalSaved),
		totalPct,
	)) + "\n\n")

	b.WriteString(fmt.Sprintf("  Saved to: %s\n", m.outputDir))
	b.WriteString("\n" + subtitleStyle.Render("  q quit"))

	return b.String()
}

func renderProgressBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return progressBarStyle.Render(bar)
}
