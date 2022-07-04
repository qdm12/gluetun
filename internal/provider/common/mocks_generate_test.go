package common

// Exceptionally, these mocks are exported since they are used by all
// provider subpackages tests, and it reduces test code duplication a lot.
//go:generate mockgen -destination=mocks.go -package $GOPACKAGE . ParallelResolver,Storage,Unzipper,Warner
