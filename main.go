package main

import (
	"fmt"
	"image/color"
	"log"
	"os/exec"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
//	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Amnezia VPN")
	w.Resize(fyne.NewSize(600, 500))

	hello := widget.NewLabel("Amnezia VPN")

	resultLabel := widget.NewLabel("Connection status:")

	statusIndicator := canvas.NewRectangle(color.RGBA{255, 0, 0, 255}) //red
	statusIndicator.SetMinSize(fyne.NewSize(15, 15))

	statusBox := container.NewHBox(
		resultLabel,
		statusIndicator,
	)

	updateStatus := func(connected bool) {
		if connected {
			statusIndicator.FillColor = color.RGBA{0, 255, 0, 255} //green
			resultLabel.SetText("Connection status: Connected")
		} else {
			statusIndicator.FillColor = color.RGBA{255, 0, 0, 255} //red
			resultLabel.SetText("Connection status: Disconnected")
		}
		canvas.Refresh(statusIndicator)
	}

//connect
	executionButton := widget.NewButton("Connect", func() {
		cmd := exec.Command("awg-quick", "up", "/home/uliana/Downloads/yliana.conf")
		err := cmd.Run()
		if err != nil {
			resultLabel.SetText(fmt.Sprintf("Error: %s", err))
			log.Println("Error executing command:", err)
			return
		}
		updateStatus(true)
	})

// disconnect
	disconnectButton := widget.NewButton("Disconnect", func() {
		cmd := exec.Command("awg-quick", "down", "/home/uliana/Downloads/yliana.conf")
		err := cmd.Run()
		if err != nil {
			resultLabel.SetText(fmt.Sprintf("Error: %s", err))
			log.Println("Error executing command:", err)
			return
		}
		updateStatus(false)
	})

//ip
	showIP := widget.NewButton("MyIP", func() {
		cmd := exec.Command("curl", "ifconfig.me")
		output, err := cmd.Output()
		if err != nil {
			resultLabel.SetText(fmt.Sprintf("Error: %s", err))
			log.Println("Error executing command:", err)
			return
		}
		resultLabel.SetText(fmt.Sprintf("IP address: %s", output))
	})

//exit
	exitButton := widget.NewButton("Exit", func() {
		a.Quit()
	})

	w.SetContent(container.NewVBox(
		hello,
		statusBox,
		executionButton,
		disconnectButton,
		showIP,
		exitButton,
	))

	w.ShowAndRun()
}
