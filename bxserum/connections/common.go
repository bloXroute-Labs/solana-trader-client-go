package connections

type Streamer[T any] func() (T, error)

func (s Streamer[T]) Channel(size int) chan T {
	ch := make(chan T, size)
	s.Into(ch)
	return ch
}

func (s Streamer[T]) Into(ch chan T) {
	go func() {
		for {
			v, err := s()
			if err != nil {
				close(ch)
				return
			}
			ch <- v
		}
	}()
}
