package ui

import (
	"log"

	"github.com/BrunoTulio/goscanner/pkg/slices"
	"github.com/BrunoTulio/goscanner/scanner"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/BrunoTulio/goscanner/bootstrap"
	"github.com/BrunoTulio/goscanner/server"
)

type MyApp struct {
	myApp                     fyne.App
	myWindow                  fyne.Window
	myDialogServerConfig      *dialog.CustomDialog
	server                    server.Server
	scanner                   scanner.Scanner
	selectedDeviceID          string
	selectDevice              *widget.Select
	startButton               *widget.Button
	stopButton                *widget.Button
	loadButton                *widget.Button
	portEntry                 *widget.Entry
	serverIsRunning           bool
	visibleDialogServerConfig bool
	devicesScan               scanner.Devices
	loadingDevice             bool
}

func (app *MyApp) GetDevice() scanner.Device {
	de, exist := slices.ContainsFn(app.devicesScan, func(d scanner.Device) bool {
		return app.selectedDeviceID == d.Name()
	})

	if !exist {
		return nil
	}

	return de

}

func (app *MyApp) Run() {
	app.myWindow = app.myApp.NewWindow("Scanner App")

	app.selectDevice = widget.NewSelect(nil, func(id string) {
		app.selectedDeviceID = id
		log.Println("Dispositivo selecionado:", id)
	})
	app.selectDevice.PlaceHolder = "Selecione dispositivo"

	app.startButton = widget.NewButton("Iniciar Digitalização", func() {
		if app.serverIsRunning {
			return
		}
		app.openDialogFormServer()
	})

	app.stopButton = widget.NewButton("Parar Digitalização", func() {
		if !app.serverIsRunning {
			return
		}
		err := app.server.Stop()
		if err != nil {
			app.serverIsRunning = true
			app.showError(err)
			app.changeBotoesStatus()
		}
		app.serverIsRunning = false
		app.changeBotoesStatus()
	})

	app.loadButton = widget.NewButton("Recarregar dispositivos", func() {
		if app.serverIsRunning {
			return
		}
		app.loadDeviceScans()
	})

	content := container.NewVBox(
		app.selectDevice,
		app.startButton,
		app.stopButton,
		app.loadButton,
	)

	app.myWindow.SetContent(content)
	app.myWindow.Resize(fyne.NewSize(300, 200))
	app.myWindow.Show()
	app.changeBotoesStatus()
	app.portEntry = widget.NewEntry()
	app.portEntry.SetPlaceHolder("Porta (opcional)")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Porta:", Widget: app.portEntry},
		},
		OnCancel: app.closeDialogFormServer,
		OnSubmit: app.submitConfigServer,
	}
	app.myDialogServerConfig = dialog.NewCustom("Configurações do Servidor", "Fechar", form, app.myWindow)
	app.myDialogServerConfig.SetOnClosed(func() {
		app.visibleDialogServerConfig = false
	})
	app.loadDeviceScans()
	app.myApp.Run()

}

func New(server server.Server, scanner scanner.Scanner) *MyApp {
	return &MyApp{
		myApp:   app.New(),
		server:  server,
		scanner: scanner,
	}
}

func (app *MyApp) changeBotoesStatus() {
	if app.startButton == nil || app.stopButton == nil || app.selectDevice == nil {
		return
	}

	if app.serverIsRunning {
		app.startButton.Disable()
		app.stopButton.Enable()
		app.selectDevice.Disable()
		app.loadButton.Disable()
		return
	}

	app.startButton.Enable()
	app.stopButton.Disable()
	app.selectDevice.Enable()
	app.loadButton.Enable()

	if app.loadingDevice {
		app.loadButton.Disable()
	}

}

func (app *MyApp) showError(err error) {
	dialog.ShowError(err, app.myWindow)
}

func (app *MyApp) showInformation(title, description string) {
	dialog.ShowInformation(title, description, app.myWindow)
}

func (app *MyApp) isValidDeviceId() bool {
	return app.selectedDeviceID != ""
}

func (app *MyApp) isVisibleDialogFormServer() bool {
	return app.myDialogServerConfig != nil && app.visibleDialogServerConfig
}

func (app *MyApp) closeDialogFormServer() {
	app.portEntry.Text = ""
	app.visibleDialogServerConfig = false
	app.myDialogServerConfig.Hide()
}

func (app *MyApp) openDialogFormServer() {

	if app.isVisibleDialogFormServer() {
		return
	}

	if !app.isValidDeviceId() {
		app.showInformation("Selecionar Dispositivo", "Por favor, selecione um dispositivo de scanner.")
		return
	}

	app.visibleDialogServerConfig = true
	app.myDialogServerConfig.Show()
}

func (app *MyApp) devicesScanNames() []string {
	names := make([]string, len(app.devicesScan))

	for i, device := range app.devicesScan {
		names[i] = device.Name()
	}
	return names
}

func (app *MyApp) loadDeviceScans() {
	if app.loadingDevice {
		return
	}
	app.selectDevice.PlaceHolder = "Aguarde carregando..."
	app.loadingDevice = true
	app.selectedDeviceID = ""
	app.selectDevice.Selected = ""
	app.selectDevice.Refresh()
	app.changeBotoesStatus()
	go func() {
		devices, err := app.scanner.List()
		app.loadingDevice = false
		app.selectDevice.PlaceHolder = "Selecione dispositivo"
		app.selectDevice.Refresh()
		app.changeBotoesStatus()

		if err != nil {
			app.showError(err)
			return
		}
		app.devicesScan = devices
		app.selectDevice.PlaceHolder = "Selecione dispositivo"
		app.selectDevice.Options = app.devicesScanNames()
		app.selectDevice.Refresh()
		app.changeBotoesStatus()
	}()
}

func (app *MyApp) submitConfigServer() {
	if app.serverIsRunning {
		app.showInformation("Digitalização", "Servidor de digitalização já está rodando")
	}

	port := app.portEntry.Text
	if port == "" {
		port = bootstrap.PortDefault
	}

	app.server.SetPort(port)

	err := app.server.IsValid()
	if err != nil {
		app.closeDialogFormServer()
		app.serverIsRunning = false
		app.showError(err)
		app.changeBotoesStatus()

		return
	}
	app.server.StartAsync()

	go func() {
		select {
		case err := <-app.server.GetStartError():
			if err != nil {
				app.closeDialogFormServer()
				app.serverIsRunning = false
				app.showError(err)
				app.changeBotoesStatus()
			}
		}
	}()

	app.serverIsRunning = true
	app.changeBotoesStatus()
	app.closeDialogFormServer()
}
