package test

import (
	"testing"
	"yunyuim/lib"
)

func TestJxToken(t *testing.T) {
	tokenStr := `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJkZXYudC5xaWRpYW55dS5jb20iLCJhdWQiOiJkZXYudC5xaWRpYW55dS5jb20iLCJpYXQiOjE2NDQzMDQxOTQsIm5iZiI6MTY0NDMwNDE5NCwiZXhwIjoxNjQ1NjAwMTk0LCJ1c2VyX3Rva2VuIjozMjMyNzgwMTAsInNob3BfdG9rZW4iOjgxOTEyMjc0LCJ0eXBlIjoibHBfbWFuYWdlX3VzZXIiLCJsb2dpbl90eXBlIjoiREVfU1lUX0g1IiwiZGV2aWNlIjoiREVfU1lUX0g1In0.utu1YvUz3Yqfmy7b4p3OLyun2PjRZrlsf_RX3Ga1Rcs`

	userInfo, err := lib.ParseToken(tokenStr)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(userInfo.UserTokenList)
}
