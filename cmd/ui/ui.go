package main

import (
	"bufio"
	"fmt"
	"hypha/app/pkg"
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

	title := widget.NewLabelWithStyle("Hypha 0.1",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})

	statusLabel := widget.NewLabelWithStyle("● Ready",
		fyne.TextAlignLeading,
		fyne.TextStyle{})

	logText := ""
	logLabel := widget.NewLabel("")
	logLabel.Wrapping = fyne.TextWrapWord

	logScroll := container.NewScroll(logLabel)
	logScroll.SetMinSize(fyne.NewSize(760, 400))

	appendLog := func(text string) {
		logText += text + "\n"
		logLabel.SetText(logText)
		logScroll.ScrollToBottom()
	}


	installBtn := widget.NewButtonWithIcon("Install Nebula",
		theme.DownloadIcon(), func() {
			statusLabel.SetText("● Installing Nebula...")
			appendLog(fmt.Sprintf("Downloading nebula v%s...", pkg.NEBULA_VERSION))

			go func() {
				err := pkg.InstallNebula()
				fyne.Do(func() {
					if err != nil {
						statusLabel.SetText("● Install failed")
						appendLog(fmt.Sprintf("ERROR: %v", err))
						log.Printf("Failed to install nebula: %v", err)
						return
					}
					statusLabel.SetText("● Nebula installed")
					appendLog(fmt.Sprintf("✓ Nebula v%s installed successfully", pkg.NEBULA_VERSION))
				})
			}()
		})

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

	//TODO: block starting nebula before unzipping
	startBtn := widget.NewButtonWithIcon("Start Nebula",
		theme.MediaPlayIcon(), func() {
			if !pkg.IfNebulaExists() {
				statusLabel.SetText("● Not installed")
				appendLog("ERROR: Nebula is not installed. Please install it first.")
				return
			}

			statusLabel.SetText("● Starting Nebula...")
			appendLog("Starting Nebula service...")

			go func() {
				pipe, err := pkg.NebulaStart(pkg.NEBULA_PATH, pkg.DESTINATION_CERTS)
				if err != nil {
					fyne.Do(func() {
						statusLabel.SetText("● Start failed")
						appendLog(fmt.Sprintf("ERROR: %v", err))
						log.Printf("Failed to start nebula: %v", err)
					})
					return
				}

				fyne.Do(func() {
					statusLabel.SetText("● Nebula running")
					appendLog("✓ Nebula started successfully")
					appendLog("--- Nebula Output ---")
				})

				scanner := bufio.NewScanner(pipe)
				for scanner.Scan() {
					line := scanner.Text()
					fyne.Do(func() {
						logText += line + "\n"
						logLabel.SetText(logText)
						logScroll.ScrollToBottom()
					})
				}

				if err := scanner.Err(); err != nil {
					fyne.Do(func() {
						appendLog(fmt.Sprintf("ERROR reading nebula output: %v", err))
					})
				}

				fyne.Do(func() {
					statusLabel.SetText("● Nebula stopped")
					appendLog("--- Nebula exited ---")
				})
			}()
		})

	// Button container
	buttonBox := container.NewGridWithColumns(3,
		installBtn,
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
