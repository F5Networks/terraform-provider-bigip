package bigip

import (
	"reflect"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
)

func TestTfMetadataToBigipMetadata(t *testing.T) {
	input := map[string]interface{}{
		"owner":       "terraform",
		"environment": "dev",
	}

	got, err := tfMetadataToBigipMetadata(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []bigip.ResourceMetadata{
		{Name: "environment", Value: "dev", Persist: "true"},
		{Name: "owner", Value: "terraform", Persist: "true"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("metadata mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestTfMetadataToBigipMetadataRejectsNonString(t *testing.T) {
	input := map[string]interface{}{
		"count": 1,
	}

	_, err := tfMetadataToBigipMetadata(input)
	if err == nil {
		t.Fatal("expected error for non-string metadata value")
	}
}

func TestBigipMetadataToTfMetadata(t *testing.T) {
	input := []bigip.ResourceMetadata{
		{Name: "environment", Value: "dev", Persist: "true"},
		{Name: "owner", Value: "terraform", Persist: "true"},
		{Name: "", Value: "skip", Persist: "true"},
	}

	got := bigipMetadataToTfMetadata(input)
	want := map[string]string{
		"environment": "dev",
		"owner":       "terraform",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("metadata map mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
}
