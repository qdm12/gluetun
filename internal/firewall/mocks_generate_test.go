package firewall

//go:generate mockgen -destination=mocks_test.go -package $GOPACKAGE . CmdRunner,Logger
