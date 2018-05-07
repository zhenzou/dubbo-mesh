package limit

func New(count int) *Limiter {
	return &Limiter{make(chan struct{}, count)}
}

type Limiter struct {
	ch chan struct{}
}

func (this *Limiter) Add() {
	this.ch <- struct{}{}
}

func (this *Limiter) Done() {
	<-this.ch
}
