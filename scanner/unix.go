//go:build linux
// +build linux

package scanner

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/jung-kurt/gofpdf"

	"github.com/tjgq/sane"
)

type (
	scannerUnix struct {
	}
	deviceUnix struct {
		dd sane.Device
	}
)

// ScanPDF implements Device.
func (d *deviceUnix) ScanPDF() ([]byte, error) {
	conn, err := sane.Open(d.dd.Name)
	if err != nil {
		return nil, fmt.Errorf("Erro ao comunicar com o scanner: %v", err)
	}
	defer conn.Close()

	params, err := conn.Params()
	if err != nil {
		return nil, fmt.Errorf("Erro ao obter parâmetros do scanner: %v", err)
	}

	buffer := make([]byte, params.BytesPerLine*params.Lines)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler do scanner: %v", err)
	}

	// Criar um novo documento PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Converter a imagem para JPEG
	img, _, err := image.Decode(bytes.NewReader(buffer[:n]))
	if err != nil {
		return nil, fmt.Errorf("Erro ao decodificar imagem: %v", err)
	}

	// Inserir imagem no PDF
	var imgBytes bytes.Buffer
	err = jpeg.Encode(&imgBytes, img, nil)
	if err != nil {
		return nil, fmt.Errorf("Erro ao codificar imagem JPEG: %v", err)
	}

	imgData := imgBytes.Bytes()
	pdf.RegisterImageReader("ScanImage", "jpg", bytes.NewReader(imgData))
	// Inserir a imagem no PDF na posição desejada
	pdf.ImageOptions("ScanImage", 10, 10, 190, 0, false, gofpdf.ImageOptions{ImageType: "jpg"}, 0, "")

	// Salvar o PDF em um buffer de bytes
	var pdfBytes bytes.Buffer
	err = pdf.Output(&pdfBytes)
	if err != nil {
		return nil, fmt.Errorf("Erro ao gerar PDF: %v", err)
	}

	return pdfBytes.Bytes(), nil
}

// Name implements Device.
func (d *deviceUnix) Name() string {
	return d.dd.Name
}

// ScanImage implements Device.
func (d *deviceUnix) ScanImage() ([]byte, error) {
	conn, err := sane.Open(d.dd.Name)

	if err != nil {
		return nil, fmt.Errorf("Erro ao comunicar com o scanner: %v", err)
	}

	defer conn.Close()

	params, err := conn.Params()
	if err != nil {
		return nil, fmt.Errorf("Erro para obter parâmetros do scanner: %v", err)
	}

	buffer := make([]byte, params.BytesPerLine*params.Lines)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("Erro ler do scanner: %w", err)
	}

	return buffer[:n], nil

}

// Close implements Device.
func (d *scannerUnix) Close() {

}

// List implements Scanner.
func (s *scannerUnix) List() (Devices, error) {
	devices, err := sane.Devices()

	if err != nil {
		return nil, fmt.Errorf("Erro ao listar dispositivos: %v", err)
	}

	dd := make(Devices, len(devices))
	for i, v := range devices {
		dd[i] = newDeviceUnix(v)
	}

	return dd, nil

}

func newDeviceUnix(d sane.Device) Device {
	return &deviceUnix{d}
}

func NewScanner() (Scanner, error) {
	err := sane.Init()
	if err != nil {
		return nil, fmt.Errorf("Erro ao utilizar sane protocol devices, %v", err)
	}
	return &scannerUnix{}, nil
}
