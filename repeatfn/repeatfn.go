package repeatfn

// RepeatFn is a function that takes a channel and a function and returns a channel
func RepeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():

			}
		}
	}()
	return valueStream
}
