package fetcher

import (
	"testing"
)

func TestStart(t *testing.T) {
	fetcher := &Fetcher{}
	Process_Url(fetcher, "http://bj.58.com/ershouche/")
}
