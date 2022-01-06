package main

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	//bytes, _ := base64.StdEncoding.DecodeString("QgtIZWxsbyBXb3JsZA==")
	// Hello World

	bytes, _ := base64.StdEncoding.DecodeString("eyJ1c2VyIjp7InVzZXJfaWQiOiJqMmdnMHMiLCJjaGFubmVsIjoiZXhhbXBsZSIsImRldmljZSI6ImoyZ2cwcydzIE1hY2Jvb2tQcm8ifX0=")
	// {"user":{"user_id":"j2gg0s","channel":"example","device":"j2gg0s's MacbookPro"}}
	fmt.Println(string(bytes))
}
