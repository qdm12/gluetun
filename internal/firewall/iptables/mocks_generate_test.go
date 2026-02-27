package iptables

//go:generate mockgen -destination=mocks_test.go -package $GOPACKAGE . CmdRunner,Logger
