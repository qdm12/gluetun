package shadowsocks

import "context"

type noopService struct{}

func (n *noopService) Start(_ context.Context) (
	runError <-chan error, err error) {
	return nil, nil
}

func (n *noopService) Stop() (err error) {
	return nil
}
