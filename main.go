package main

import (
	"fmt"
	"image/color"
	"os"
	"log"
	"os/exec"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
)

func main() {

	//env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	iconPath := os.Getenv("ICON_PATH")
	confPath := os.Getenv("CONF_PATH")

	//iconApp
	a := app.New()
	r, err := fyne.LoadResourceFromPath(iconPath)
	if err != nil {
    	    log.Println("Error loading icon:", err)
	} else if r != nil {
	    a.SetIcon(r)
	}


	w := a.NewWindow("Amnezia VPN")
	w.Resize(fyne.NewSize(600, 500))


	hello := widget.NewLabel("Amnezia VPN")
	centeredHello := container.NewCenter(hello)

	resultLabel := widget.NewLabel("Connection status:")

	//squareStatus
	statusIndicator := canvas.NewRectangle(color.RGBA{255, 0, 0, 255}) //red
	statusIndicator.SetMinSize(fyne.NewSize(35, 15))

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
		cmd := exec.Command("awg-quick", "up", confPath)
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
		cmd := exec.Command("awg-quick", "down", confPath)
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
		currentStatusText := resultLabel.Text
		cmd := exec.Command("curl", "ifconfig.me")
		output, err := cmd.Output()
		if err != nil {
			resultLabel.SetText(fmt.Sprintf("Error: %s", err))
			log.Println("Error executing command:", err)
			return
		}
		resultLabel.SetText(fmt.Sprintf("IP address: %s", output))

	        go func() {
                        time.Sleep(3 * time.Second)
			resultLabel.SetText(currentStatusText)
			canvas.Refresh(statusIndicator)
                }()
	})




	//exit
	exitButton := widget.NewButton("Exit", func() {
		a.Quit()
	})

	w.SetContent(container.NewVBox(
		centeredHello,
		statusBox,
		executionButton,
		disconnectButton,
		showIP,
		exitButton,
	))

	w.ShowAndRun()
}
