package send

// Flusher is used to flush any buffers that may need cleared before exiting the [EmailEvent] function. This is useful
// in the context of lambdas like [GCP Cloud Functions] where it is not reliable that there will be CPU or memory
// allocated after the [EmailEvent] function ends.
//
// [GCP Cloud Functions]: https://cloud.google.com/functions
type Flusher interface {
	Flush() error
}

type noopFlusher struct{}

func (n noopFlusher) Flush() error {
	return nil
}
