package common

// Exceptionally, these mocks are exported since they are used by all
// provider subpackages tests, and it reduces test code duplication a lot.
// Note mocks.go might need to be removed before re-generating it.
//go:generate mockgen -destination=mocks.go -package $GOPACKAGE . ParallelResolver,Storage,Unzipper,Warner
