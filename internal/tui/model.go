package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/yazdanctx/bitc-go/internal/compress"
	"github.com/yazdanctx/bitc-go/internal/scanner"
)

type state int

const (
	stateScanning state = iota
	stateConfig
	stateCompressing
	stateResults
)

type progressInfo struct {
	Done  int
	Total int
}

type model struct {
	state        state
	dir          string
	outputDir    string
	images       []compress.ImageFile
	formats      []compress.FormatOption
	cursor       int
	formatCursor int
	spinner      spinner.Model
	currentFile  string
	progress     progressInfo
	results      []compress.CompressResult
	summary      *compress.CompressionSummary
	compressCh   chan compress.ProgressMsg
	err          error
}

type ModelAccessor interface {
	GetSummary() *compress.CompressionSummary
}

func (m model) GetSummary() *compress.CompressionSummary {
	return m.summary
}

func InitialModel(dir string, outputDir string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		state:     stateScanning,
		dir:       dir,
		outputDir: outputDir,
		formats:   compress.AllFormats(),
		spinner:   s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.scanDirectory())
}

func (m model) scanDirectory() tea.Cmd {
	return func() tea.Msg {
		images, err := scanner.ScanDirectory(m.dir)
		if err != nil {
			return scanDoneMsg{err: err}
		}
		return scanDoneMsg{images: images}
	}
}

type scanDoneMsg struct {
	images []compress.ImageFile
	err    error
}

type compressMsg struct {
	msg compress.ProgressMsg
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.state == stateConfig && m.canStart() {
				return m.startCompression()
			}
		case " ":
			if m.state == stateConfig {
				m.formats[m.formatCursor].Enabled = !m.formats[m.formatCursor].Enabled
			}
		case "up", "k":
			if m.state == stateConfig {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if m.state == stateConfig {
				if m.cursor < len(m.images)-1 {
					m.cursor++
				}
			}
		case "left", "h":
			if m.state == stateConfig {
				if m.formatCursor > 0 {
					m.formatCursor--
				}
			}
		case "right", "l":
			if m.state == stateConfig {
				if m.formatCursor < len(m.formats)-1 {
					m.formatCursor++
				}
			}
		}

	case scanDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.images = msg.images
		m.state = stateConfig

	case compressMsg:
		cm := msg.msg
		if cm.Finished {
			m.summary = cm.Summary
			m.state = stateResults
			return m, tea.Quit
		}
		m.currentFile = cm.Current
		m.progress = progressInfo{Done: cm.Done, Total: cm.Total}
		if cm.Result != nil {
			m.results = append(m.results, *cm.Result)
		}
		return m, waitForCompress(m.compressCh)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render("Error: " + m.err.Error()) + "\n"
	}

	switch m.state {
	case stateScanning:
		return m.viewScanning()
	case stateConfig:
		return m.viewConfig()
	case stateCompressing:
		return m.viewCompressing()
	case stateResults:
		return m.viewResults()
	default:
		return ""
	}
}

func (m model) canStart() bool {
	for _, f := range m.formats {
		if f.Enabled {
			return true
		}
	}
	return false
}

func (m model) startCompression() (tea.Model, tea.Cmd) {
	m.state = stateCompressing
	m.compressCh = make(chan compress.ProgressMsg, 100)

	return m, tea.Batch(
		m.spinner.Tick,
		startCompressionCmd(m.images, m.formats, m.outputDir, m.compressCh),
	)
}

func startCompressionCmd(
	images []compress.ImageFile,
	formats []compress.FormatOption,
	outDir string,
	ch chan compress.ProgressMsg,
) tea.Cmd {
	return func() tea.Msg {
		go compress.RunCompression(images, formats, outDir, ch)
		return waitForCompress(ch)
	}
}

func waitForCompress(ch chan compress.ProgressMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return compressMsg{msg: compress.ProgressMsg{Finished: true}}
		}
		return compressMsg{msg: msg}
	}
}
