package object

type RGWDriver struct {
	s3 *S3Driver
}

func NewRGWDriver(endpoint, bucket string) *RGWDriver {
	return &RGWDriver{s3: NewS3Driver(endpoint, bucket)}
}

func (r *RGWDriver) Put(key string, data []byte) error {
	return r.s3.Put(key, data)
}

func (r *RGWDriver) Get(key string) ([]byte, error) {
	return r.s3.Get(key)
}

func (r *RGWDriver) Delete(key string) error {
	return r.s3.Delete(key)
}
