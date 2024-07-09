//go:build windows
// +build windows

package scanner

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type (
	scannerWindows struct {
		m       *ole.IDispatch
		unknown *ole.IUnknown
	}
	deviceWindows struct {
		dd   *ole.IDispatch
		name string
	}
)

func newDeviceWindows(dd *ole.IDispatch, name string) Device {
	return &deviceWindows{dd, name}
}

func (d *deviceWindows) Name() string {
	return d.name
}

func (d *scannerWindows) Close() {
	d.unknown.Release()
	d.m.Release()
	ole.CoUninitialize()
}

func (d *deviceWindows) ScanPDF() ([]byte, error) {
	// Configurar a digitalização para produzir um documento PDF
	_, err := oleutil.PutProperty(d.dd, "Format", "PDF")
	if err != nil {
		return nil, fmt.Errorf("Erro ao configurar o formato PDF: %v", err)
	}

	// Executar a digitalização
	image, err := oleutil.CallMethod(d.dd, "Items", 1)
	if err != nil {
		return nil, fmt.Errorf("Erro ao digitalizar o documento PDF: %v", err)
	}

	// Obter os bytes do documento PDF
	pdfBytes, err := oleutil.GetProperty(image.ToIDispatch(), "FileData")
	if err != nil {
		return nil, fmt.Errorf("Erro ao obter os bytes do PDF: %v", err)
	}

	return pdfBytes.Value().([]byte), nil
}

// ScanImage implements Device.
func (d *deviceWindows) ScanImage() ([]byte, error) {
	image, err := oleutil.CallMethod(d.dd, "Items", 1)
	if err != nil {
		return nil, fmt.Errorf("Erro ao digitalizar a imagem: %v", err)

	}

	// Get image bytes
	imageBytes, err := oleutil.GetProperty(image.ToIDispatch(), "FileData")
	if err != nil {
		return nil, fmt.Errorf("Erro ao obter os bytes da imagem: %v", err)

	}

	return imageBytes.Value().([]byte), nil
}

func (s *scannerWindows) List() (Devices, error) {
	devices, err := oleutil.CallMethod(s.m, "DeviceInfos", 1) // 1 means WIA_DEVICE_TYPE.Scanner
	if err != nil {
		return nil, fmt.Errorf("Erro ao obter os dispositivos de scanner: %v", err)

	}
	deviceCollection := devices.ToIDispatch()

	count, err := oleutil.GetProperty(deviceCollection, "Count")
	if err != nil {
		return nil, fmt.Errorf("Erro ao obter os dispositivos de scanner: %v", err)

	}
	numDevices := int(count.Val)
	dd := Devices{}

	for i := 1; i <= numDevices; i++ {
		device, err := oleutil.CallMethod(deviceCollection, "Item", i)
		if err != nil {
			fmt.Println("Erro ao obter o dispositivo:", err)
			continue
		}
		name, err := oleutil.GetProperty(device.ToIDispatch(), "Name")
		if err != nil {
			fmt.Println("Erro ao obter o nome do dispositivo:", err)
			continue
		}

		dd = append(dd, newDeviceWindows(device.ToIDispatch(), name.ToString()))

	}

	return dd, nil
}

func NewScanner() (Scanner, error) {
	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	// defer ole.CoUninitialize()

	// Create WIA Device Manager object
	unknown, err := oleutil.CreateObject("WIA.DeviceManager")
	if err != nil {
		return nil, fmt.Errorf("Erro ao criar o objeto WIA.DeviceManager: %v", err)
	}
	// defer unknown.Release()

	// Query for IDispatch interface
	manager, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("Erro ao obter a interface IDispatch: %v", err)

	}
	// defer manager.Release()

	return &scannerWindows{manager, unknown}, nil
}
