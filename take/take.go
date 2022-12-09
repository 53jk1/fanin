package take

// Take is a function that takes a channel and a number and returns a channel
func Take(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {

	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 1; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:

			}
		}
	}()
	return takeStream
}
