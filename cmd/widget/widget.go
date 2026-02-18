package main

import (
	"fmt"
	"hypha/app/pkg"
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {

	a := app.New()
	w := a.NewWindow("Hypha🍄")
	w.Resize(fyne.NewSize(1920, 1080))

	// Header
	title := widget.NewLabelWithStyle("Hypha 0.1",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})

	// Status indicator
	statusLabel := widget.NewLabelWithStyle("● Ready",
		fyne.TextAlignLeading,
		fyne.TextStyle{})

	// Log output as a Label (can be updated even though it looks read-only)
	logText := ""
	logLabel := widget.NewLabel("")
	logLabel.Wrapping = fyne.TextWrapWord

	logScroll := container.NewScroll(logLabel)
	logScroll.SetMinSize(fyne.NewSize(760, 400))

	// Helper to append to log - now updates label
	appendLog := func(text string) {
		logText += text + "\n"
		logLabel.SetText(logText)
		logScroll.ScrollToBottom()
	}

	// Buttons
	unzipBtn := widget.NewButtonWithIcon("Unzip Certificates",
		theme.FolderOpenIcon(), func() {
			statusLabel.SetText("● Unzipping...")
			appendLog("Starting certificate extraction...")

			err := pkg.Unzip(pkg.HOST_PATH, pkg.DESTINATION_FOLDER)
			if err != nil {
				statusLabel.SetText("● Unzip failed")
				appendLog(fmt.Sprintf("ERROR: %v", err))
				log.Printf("Unzip error: %v", err)
				return
			}

			statusLabel.SetText("● Unzip complete")
			appendLog("✓ Certificates extracted successfully")
		})

	startBtn := widget.NewButtonWithIcon("Start Nebula",
		theme.MediaPlayIcon(), func() {
			statusLabel.SetText("● Starting Nebula...")
			appendLog("Starting Nebula service...")

			pipe, err := pkg.NebulaStart(pkg.NEBULA_PATH, pkg.DESTINATION_CERTS)
			if err != nil {
				statusLabel.SetText("● Start failed")
				appendLog(fmt.Sprintf("ERROR: %v", err))
				log.Printf("Failed to start nebula: %v", err)
				return
			}

			statusLabel.SetText("● Nebula running")
			appendLog("✓ Nebula started successfully")
			appendLog("--- Nebula Output ---")

			// Create channel for lines
			linesChan := make(chan string, 100)

			// Read from pipe using raw bytes
			go func() {
				defer close(linesChan)
				buf := make([]byte, 1024)
				for {
					n, err := pipe.Read(buf)
					if n > 0 {
						line := string(buf[:n])
						fmt.Println("RAW READ:", line) // Debug
						linesChan <- line
					}
					if err != nil {
						if err != io.EOF {
							linesChan <- fmt.Sprintf("ERROR: %v", err)
						}
						break
					}
				}
			}()

			// Update UI from channel
			go func() {
				for line := range linesChan {
					currentLine := line
					fyne.Do(func() {
						logText += currentLine + "\n"
						logLabel.SetText(logText)
						logScroll.ScrollToBottom()
					})
				}
			}()
		})

	// Button container
	buttonBox := container.NewGridWithColumns(2,
		unzipBtn,
		startBtn,
	)

	// Main layout
	content := container.NewBorder(
		// Top
		container.NewVBox(
			title,
			widget.NewSeparator(),
			statusLabel,
			widget.NewSeparator(),
		),
		// Bottom
		container.NewVBox(
			widget.NewSeparator(),
			buttonBox,
		),
		// Left, Right
		nil, nil,
		// Center
		logScroll,
	)

	w.SetContent(container.NewPadded(content))
	w.ShowAndRun()
}
