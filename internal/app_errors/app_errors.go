package app_errors

type BucketExhaustedError struct{}

func (e BucketExhaustedError) Error() string {
	return "Bucket is exhausted"
}
