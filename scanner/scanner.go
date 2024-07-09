package scanner

type (
	Scanner interface {
		List() (Devices, error)
		Close()
	}
	Devices []Device

	Device interface {
		Name() string
		ScanImage() ([]byte, error)
		ScanPDF() ([]byte, error)
	}
)
