package main
import (
	"fmt"
	"log"
	"hypha/app/pkg"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	
	a := app.New()
	w := a.NewWindow("Hypha🍄")
	
	statusLabel := widget.NewLabel("Ready")

	unzipBtn := widget.NewButton("Unzip Certificates", func() {
		statusLabel.SetText("Unzipping...")
		err := pkg.Unzip(pkg.HOST_PATH, pkg.DESTINATION_FOLDER)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Unzip failed: %v", err))
			log.Printf("Unzip error: %v", err)
			return
		}
		statusLabel.SetText("Unzip complete")
	})

	parseBtn := widget.NewButton("Parse Certificates", func() {
		statusLabel.SetText("Parsing...")
		_, err := pkg.ParseCertFolder(pkg.DESTINATION_CERT_PATH)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Parse failed: %v", err))
			log.Printf("Parse error: %v", err)
			return
		}
		statusLabel.SetText("Parse complete")
	})

	startBtn := widget.NewButton("Start Nebula", func() {
		statusLabel.SetText("Starting Nebula...")
		err := pkg.NebulaStart()
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Start failed: %v", err))
			log.Printf("NebulaStart error: %v", err)
			return
		}
		statusLabel.SetText("Nebula started")
	})

	w.SetContent(container.NewVBox(
		statusLabel,
		unzipBtn,
		parseBtn,
		startBtn,
	))
	w.ShowAndRun()
}
