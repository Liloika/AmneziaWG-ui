package main

import (
    "fmt"
    "image/color"
    "os"
    "log"
    "os/exec"
    "time"
    "bytes"
    "path/filepath"
    "strings"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func checkInterfaceStatus(interfaceName string) bool {
    cmd := exec.Command("ip", "link", "show", interfaceName)
    output, err := cmd.Output()
    if err != nil {
        return false
    }
    
    return bytes.Contains(output, []byte("UP"))
}

func main() {
    isConnected := false    

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
        isConnected = connected
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
        if isConnected {
            resultLabel.SetText("Already connected!")
            
            go func() {
                time.Sleep(3 * time.Second)
                if isConnected {
                    resultLabel.SetText("Connection status: Connected")
                } else {
                    resultLabel.SetText("Connection status: Disconnected")
                }
                canvas.Refresh(statusIndicator)
            }()
            return
        }

        interfaceName := strings.TrimSuffix(filepath.Base(confPath), ".conf")
        if checkInterfaceStatus(interfaceName) {
            updateStatus(true)
            resultLabel.SetText("Already connected!")
            
            go func() {
                time.Sleep(3 * time.Second)
                resultLabel.SetText("Connection status: Connected")
                canvas.Refresh(statusIndicator)
            }()
            return
        }

        cmd := exec.Command("sudo", "awg-quick", "up", confPath)
        var out bytes.Buffer
        var stderr bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &stderr
        err := cmd.Run()
        if err != nil {
            fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
            resultLabel.SetText(fmt.Sprintf("Error: %s", err))
            log.Println("Error executing command:", err)
            return
        }
        updateStatus(true)
    })

    // disconnect
    disconnectButton := widget.NewButton("Disconnect", func() {
        if !isConnected {
            resultLabel.SetText("Already disconnected!")
            
            go func() {
                time.Sleep(3 * time.Second)
                if isConnected {
                    resultLabel.SetText("Connection status: Connected")
                } else {
                    resultLabel.SetText("Connection status: Disconnected")
                }
                canvas.Refresh(statusIndicator)
            }()
            return
        }

        cmd := exec.Command("sudo", "awg-quick", "down", confPath)
        var out bytes.Buffer
        var stderr bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &stderr
        err := cmd.Run()
        if err != nil {
            fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
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
