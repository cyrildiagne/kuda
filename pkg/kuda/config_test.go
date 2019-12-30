package kuda

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseURL(t *testing.T) {
	urlA := "https://test-url.default.example.com"
	expectedA := URLConfig{
		Scheme:    "https",
		Name:      "test-url",
		Namespace: "default",
		Domain:    "example.com",
	}
	resultA, err := ParseURL(urlA)
	if err != nil {
		t.Errorf("err")
	}
	if diff := cmp.Diff(expectedA, *resultA); diff != "" {
		t.Errorf("TestParseURL() mismatch (-want +got):\n%s", diff)
	}

	urlB := "http://test-url2.default.1.2.3.4.xip.io/run"
	expectedB := URLConfig{
		Scheme:    "http",
		Name:      "test-url2",
		Namespace: "default",
		Domain:    "1.2.3.4.xip.io",
	}
	resultB, err := ParseURL(urlB)
	if err != nil {
		t.Errorf("err")
	}
	// if !reflect.DeepEqual(*resultB, expectedB) {
	// 	t.Errorf("Result B error. Got: \n%v, \nExpected: \n%v", resultB, expectedB)
	// }
	if diff := cmp.Diff(expectedB, *resultB); diff != "" {
		t.Errorf("TestParseURL() mismatch (-want +got):\n%s", diff)
	}
}
