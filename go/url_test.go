package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	cUrl, _ := url.Parse("https://im-images.qiandaoapp.com/videos/snapshots/788ebed32b0100d42d5bc571e32c3dc4.mp4.jpg")
	fmt.Println(cUrl.RawPath)
	fmt.Println(cUrl.User)
	fmt.Println(cUrl.Host)
	fmt.Println(cUrl.Fragment)
	fmt.Println(cUrl.Opaque)
	fmt.Println(cUrl.Path)
	fmt.Println(cUrl.RawFragment)
	fmt.Println(cUrl.Scheme)
}
