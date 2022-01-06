package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	token, err := getAccessToken(GrantType, AppId, Secret)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(token)
}

func TestGetUnionId(t *testing.T) {
	token, err := getAccessToken(GrantType, AppId, Secret)
	if err != nil {
		t.Fatal(err)
	}

	union, err := getUnionId(token, "omh-u6g_BUesXCwhRyM1bBfzCHN4")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(union)
}

func TestSha256(t *testing.T) {
	hash := sha256.New()
	hash.Write([]byte("hello1"))
	fmt.Printf("%x\n", hash.Sum(nil))
	encoding := hex.EncodeToString(hash.Sum(nil))
	fmt.Println(encoding)
	// hello 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
	// hello 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
	// hello1 91e9240f415223982edc345532630710e94a7f52cd5f48f5ee1afc555078f0ab
}

func TestRand(t *testing.T) {
	b := make([]byte, 32)
	n, err := rand.Read(b)
	if err == nil {
		fmt.Println(b[:n], " ", n)
		fmt.Println(base64.StdEncoding.EncodeToString(b))
	}

	b1 := make([]byte, 32)
	n1, err := rand.Read(b1)
	if err == nil {
		fmt.Println(b1[:n1], " ", n1)
		fmt.Println(base64.StdEncoding.EncodeToString(b1))
	}
}