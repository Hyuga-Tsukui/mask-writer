package masker

import (
	"bytes"
	"testing"

	"github.com/goccy/go-json"
)

func TestMaskWriter_NoKeyNoChange(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"name": "John Doe",
		"company": map[string]interface{}{
			"city": "New York City",
		},
	}

	expected := map[string]interface{}{
		"name": "John Doe",
		"company": map[string]interface{}{
			"city": "New York City",
		},
	}

	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	keysToMask := []string{"email", "password"}
	buffer := bytes.NewBuffer([]byte{})
	maskWriter := NewMaskWriter(buffer, keysToMask, "******")
	_, err = maskWriter.Write(b)
	if err != nil {
		t.Fatal(err)
	}

	result := buffer.Bytes()
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, expectedBytes) {
		t.Errorf("Execution result differs from expected value. Got %s, wanted %s", result, expectedBytes)
	}
}

func TestMaskWriter_Write_MaskingKeyword(t *testing.T) {
	t.Parallel()
	input := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john.doe@example.com",
		"password": "Passw@ard123",
		"company": map[string]interface{}{
			"city":  "New York City",
			"email": "company@example.com",
		},
		"friends": []map[string]interface{}{
			{
				"name":     "Jane Smith",
				"email":    "jane.smith@example.com",
				"password": "abcdef",
			},
			{
				"name":     "Bob Johnson",
				"email":    "bob.johnson@example.com",
				"password": "abcdef",
			},
		},
	}

	expected := map[string]interface{}{
		"name":     "John Doe",
		"email":    "******",
		"password": "******",
		"company": map[string]interface{}{
			"city":  "New York City",
			"email": "******",
		},
		"friends": []map[string]interface{}{
			{
				"name":     "Jane Smith",
				"email":    "******",
				"password": "******",
			},
			{
				"name":     "Bob Johnson",
				"email":    "******",
				"password": "******",
			},
		},
	}

	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	keysToMask := []string{"email", "password"}
	buffer := bytes.NewBuffer([]byte{})
	maskWriter := NewMaskWriter(buffer, keysToMask, "******")
	_, err = maskWriter.Write(b)
	if err != nil {
		t.Fatal(err)
	}

	result := buffer.Bytes()
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, expectedBytes) {
		t.Errorf("Execution result differs from expected value. Got %s, wanted %s", result, expectedBytes)
	}
}

func TestMaskWriter_Write_Error_InvalidJSON(t *testing.T) {
	keysToMask := []string{"email", "password"}
	mask := "******"
	buffer := bytes.NewBuffer([]byte{})
	maskWriter := NewMaskWriter(buffer, keysToMask, mask)

	invalidJSON := []byte(`{"name": "John Doe", "email": "john.doe@example.com", "password": "Passw@rd123", invalid}`)

	_, err := maskWriter.Write(invalidJSON)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func BenchmarkMaskWriter_Write(b *testing.B) {
	input := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john.doe@example.com",
		"password": "Passw@ard123",
		"company": map[string]interface{}{
			"city":  "New York City",
			"email": "company@example.com",
		},
		"friends": []map[string]interface{}{
			{
				"name":     "Jane Smith",
				"email":    "jane.smith@example.com",
				"password": "abcdef",
			},
			{
				"name":     "Bob Johnson",
				"email":    "bob.johnson@example.com",
				"password": "abcdef",
			},
		},
	}
	keysToMask := []string{"email", "password"}
	mask := "******"
	buffer := bytes.NewBuffer([]byte{})
	maskWriter := NewMaskWriter(buffer, keysToMask, mask)
	inputBytes, _ := json.Marshal(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maskWriter.Write(inputBytes)
		buffer.Reset()
	}
}
