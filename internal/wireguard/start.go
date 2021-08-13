package wireguard

import "context"

type Starter interface {
	Start(ctx context.Context, done chan<- error)
}

func (w *Wireguard) Start(ctx context.Context, done chan<- error) {
	go func() {
		done <- w.Run(ctx)
	}()
}
