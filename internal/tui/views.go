package tui

import (
	"fmt"
	"strings"

	"github.com/yazdun/bitc-go/internal/compress"
)

func (m model) viewScanning() string {
	return "\n  " + m.spinner.View() + "  Scanning directory...\n"
}

func (m model) viewConfig() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" bitc-go ") + "\n\n")
	b.WriteString(subtitleStyle.Render(fmt.Sprintf("Found %d images in %s", len(m.images), m.dir)) + "\n\n")

	b.WriteString("  Files:\n")
	for i, img := range m.images {
		cursor := "  "
		if i == m.cursor {
			cursor = fileSelectedStyle.Render("▸ ")
		}
		name := fileDefaultStyle.Render(img.Name)
		size := subtitleStyle.Render(FormatSize(img.OrigSize))
		b.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, name, size))
	}

	b.WriteString("\n  Formats:\n")
	for i, f := range m.formats {
		cursor := "  "
		if i == m.formatCursor {
			cursor = fileSelectedStyle.Render("▸ ")
		}
		checkbox := checkboxOffStyle.Render("[ ]")
		if f.Enabled {
			checkbox = checkboxOnStyle.Render("[✓]")
		}
		label := fileDefaultStyle.Render(f.Label)
		b.WriteString(fmt.Sprintf("  %s %s %s\n", cursor, checkbox, label))
	}

	b.WriteString("\n")
	if m.canStart() {
		b.WriteString("  " + startButtonStyle.Render(" Start ") + "  (enter)\n")
	} else {
		b.WriteString("  " + startButtonInactiveStyle.Render(" Start ") + "  enable at least one format\n")
	}

	b.WriteString("\n" + subtitleStyle.Render("  ↑/↓ navigate  space toggle format  enter start  q quit"))

	return b.String()
}

func (m model) viewCompressing() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" bitc-go ") + "\n\n")

	if m.progress.Total > 0 {
		pct := float64(m.progress.Done) / float64(m.progress.Total) * 100
		bar := renderProgressBar(pct, 40)
		b.WriteString(fmt.Sprintf("  %s  %d/%d (%.0f%%)\n\n", bar, m.progress.Done, m.progress.Total, pct))
	}

	if m.currentFile != "" {
		b.WriteString(fmt.Sprintf("  %s Compressing %s...\n", m.spinner.View(), m.currentFile))
	}

	b.WriteString("\n" + subtitleStyle.Render("  q quit"))
	return b.String()
}

func (m model) viewResults() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" bitc-go ") + "\n\n")

	if m.summary == nil {
		return b.String()
	}

	b.WriteString(fmt.Sprintf("  %-30s %10s %12s %10s %10s\n", "File", "Original", "Format", "Compressed", "Savings"))
	b.WriteString(strings.Repeat("─", 76) + "\n")

	perImage := make(map[string][]compress.CompressResult)
	for _, r := range m.summary.Results {
		perImage[r.Image.Name] = append(perImage[r.Image.Name], r)
	}

	for _, img := range m.images {
		results := perImage[img.Name]
		for i, r := range results {
			if r.Err != nil {
				continue
			}
			name := img.Name
			if i > 0 {
				name = ""
			}
			b.WriteString(fmt.Sprintf("  %-30s %10s %12s %10s %s\n",
				name,
				FormatSize(r.OriginalSize),
				string(r.Format),
				FormatSize(r.CompressedSize),
				FormatSavings(r.Savings),
			))
		}
	}

	b.WriteString(strings.Repeat("─", 76) + "\n")
	totalPct := 0.0
	if m.summary.TotalOriginal > 0 {
		totalPct = float64(m.summary.TotalSaved) / float64(m.summary.TotalOriginal) * 100
	}
	b.WriteString(summaryStyle.Render(fmt.Sprintf("  Total: %s → saved %s (%.1f%%) — Best: %s",
		FormatSize(m.summary.TotalOriginal),
		FormatSize(m.summary.TotalSaved),
		totalPct,
		string(m.summary.BestFormat),
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
