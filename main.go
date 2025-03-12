package main

import (
    "fmt"
    "log"
    "os/exec"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func main() {
    a := app.New()
    w := a.NewWindow("Amnezia VPN")
    w.Resize(fyne.NewSize(600, 500))

    hello := widget.NewLabel("Amnezia VPN")

    resultLabel := widget.NewLabel("Connection status: ")
    resultLabel.Wrapping = fyne.TextWrapWord


    executionButton := widget.NewButton("Connect", func() {
        cmd := exec.Command("awg-quick", "up", "/home/uliana/Downloads/yliana.conf")
        err := cmd.Run()
        if err != nil {
            resultLabel.SetText(fmt.Sprintf("Error: %s", err))
            log.Println("Error executing command:", err)
            return
        }

        resultLabel.SetText("Connection")

   })

    disconnectButton := widget.NewButton("Disconnect", func() {
	cmd := exec.Command("awg-quick", "down", "/home/uliana/Downloads/yliana.conf")
	err := cmd.Run()
	if err != nil {
            resultLabel.SetText(fmt.Sprintf("Error: %s", err))
            log.Println("Error executing command:", err)
            return
        }
	resultLabel.SetText("Disconnected")

   })

    exitButton := widget.NewButton("Exit", func() {
        a.Quit()
    })

    w.SetContent(container.NewVBox(
	hello,
	resultLabel,
	executionButton,
	disconnectButton,
	showIP,
        exitButton,
    ))

    w.ShowAndRun()
}
