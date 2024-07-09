package server

type Server interface {
	StartAsync()
	Stop() error
	IsValid() error
	SetPort(port string)
	GetStartError() <-chan error
}
