package common

// Exceptionally, the storage mock is exported since it is used by all
// provider subpackages tests, and it reduces test code duplication a lot.
//go:generate mockgen -destination=mocks.go -package $GOPACKAGE . Storage
