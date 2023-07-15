package tests

import(
	"time"
)


// Await can wait for some time until 
type Await struct {
	maxWait time.Duration // Max duration to wait before giving up
	poll	time.Duration // How often check for the predicated
	waited  time.Duration // time waited before stopping to wait
}

// Wait sleeping until the predicate function returns true or the maximum waiting time is reached
func (aw *Await) Wait( predicate func() bool ) bool {

	ticker := time.NewTicker(aw.poll)

	// Done channel, indicates the wait is over when it receives a value
	// From the timer if predicate function returns true, or if the maximum waiting time is over
	done := make(chan bool)

	start := time.Now()

	for {
		select {
			case <-ticker.C:
				go func() {
					if predicate() {
						aw.waited = time.Since(start)
						done <- true
					}
				}()

			case <-time.After(aw.maxWait):
				aw.waited = time.Since(start)
				// Predicate never raised before
				done <- false

			case result := <-done:
				return result
		}

	}
}

// TimeWaited returns time waited before stopping to wait
// Interpretation differs if Wait() returned true (predicate function returned true) or false (timeout)
func (aw *Await) TimeWaited() time.Duration {
	return aw.waited	
}

// NewAwait creates a new Await instance 
func NewAwait(maxWait time.Duration, pollTime time.Duration) *Await {
	return &Await{maxWait: maxWait, poll: pollTime}
}