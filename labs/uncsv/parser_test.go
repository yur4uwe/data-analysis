package uncsv

import (
	"encoding/csv"
	"strings"
	"testing"
)

// TestNewDecoder tests the NewDecoder constructor
func TestNewDecoder(t *testing.T) {
	csvData := "header1,header2\nvalue1,value2"
	reader := csv.NewReader(strings.NewReader(csvData))

	decoder := NewDecoder(reader)

	if decoder == nil {
		t.Fatal("NewDecoder returned nil")
	}
}

// TestDecodeEmptyData tests decoding with empty CSV data
func TestDecodeEmptyData(t *testing.T) {
	csvData := ""
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result struct {
		Field1 []string
		Field2 []int
	}

	err := decoder.Decode(&result)
	// Should handle empty data gracefully
	if err != nil {
		t.Logf("Decode with empty data returned error: %v", err)
	}
}

// TestDecodeSimpleStruct tests decoding into a simple struct of arrays
func TestDecodeSimpleStruct(t *testing.T) {
	csvData := "name,age\nAlice,30\nBob,25\nCharlie,35"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result struct {
		Name []string `csv:"name"`
		Age  []int    `csv:"age"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}
}

// TestDecodeWithFloats tests decoding struct with float arrays
func TestDecodeWithFloats(t *testing.T) {
	csvData := "id,value,price\n1,10.5,20.99\n2,15.3,30.50\n3,12.7,25.75"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result struct {
		ID    []int     `csv:"id"`
		Value []float64 `csv:"value"`
		Price []float64 `csv:"price"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode with floats failed: %v", err)
	}
}

// TestDecodeSingleRow tests decoding a single row
func TestDecodeSingleRow(t *testing.T) {
	csvData := "name,status\nJohn,active"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result struct {
		Name   []string `csv:"name"`
		Status []string `csv:"status"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode single row failed: %v", err)
	}
}

// TestDecodeMultipleFields tests decoding with many fields
func TestDecodeMultipleFields(t *testing.T) {
	csvData := "a,b,c,d,e\n1,2,3,4,5\n6,7,8,9,10"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result struct {
		A []int `csv:"a"`
		B []int `csv:"b"`
		C []int `csv:"c"`
		D []int `csv:"d"`
		E []int `csv:"e"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode multiple fields failed: %v", err)
	}
}

// TestDecodeWithMissingValues tests decoding with missing/empty values
func TestDecodeWithMissingValues(t *testing.T) {
	csvData := "name,value\nAlice,100\nBob,\nCharlie,200"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result struct {
		Name  []string `csv:"name"`
		Value []string `csv:"value"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Logf("Decode with missing values returned error: %v", err)
	}
}

// TestDecodeIntoNil tests Decode with nil pointer
func TestDecodeIntoNil(t *testing.T) {
	csvData := "a,b\n1,2"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	err := decoder.Decode(nil)
	if err == nil {
		t.Error("Decode(nil) should return an error")
	}
}

// TestDecodeIntoNonStruct tests Decode with non-struct type
func TestDecodeIntoNonStruct(t *testing.T) {
	csvData := "a,b\n1,2"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result int
	err := decoder.Decode(&result)
	if err == nil {
		t.Error("Decode into non-struct should return an error")
	}
}

// TestDecodeIntoNonArrayFields tests Decode where struct fields are not arrays
func TestDecodeIntoNonArrayFields(t *testing.T) {
	csvData := "name,age\nAlice,30"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result struct {
		Name string `csv:"name"` // Not an array
		Age  int    `csv:"age"`  // Not an array
	}

	err := decoder.Decode(&result)
	if err == nil {
		t.Error("Decode into non-array fields should return an error")
	}
}

// TestDecodeFieldDecoderInterface tests FieldDecoder interface implementation
func TestDecodeFieldDecoderInterface(t *testing.T) {
	csvData := "value\n100\n200\n300"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result struct {
		Value []CustomType `csv:"value"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Logf("Decode with FieldDecoder implementation: %v", err)
	}
}

// CustomType implements FieldDecoder for testing
type CustomType struct {
	Value int
}

func (c *CustomType) DecodeCSV(s string) error {
	// Simple implementation for testing
	_, err := csv.NewReader(strings.NewReader(s)).Read()
	return err
}

// TestDecodeWithSpecialCharacters tests decoding CSV with special characters
func TestDecodeWithSpecialCharacters(t *testing.T) {
	csvData := "description\n\"Hello, World\"\n\"Line1\nLine2\"\n\"Quote\"\"Test\""
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result struct {
		Description []string `csv:"description"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Logf("Decode with special characters: %v", err)
	}
}

// TestDecodeMultipleTimes tests calling Decode multiple times
func TestDecodeMultipleTimes(t *testing.T) {
	csvData := "name,value\nAlice,100\nBob,200"
	reader := csv.NewReader(strings.NewReader(csvData))
	decoder := NewDecoder(reader)

	var result1 struct {
		Name  []string `csv:"name"`
		Value []int    `csv:"value"`
	}

	err := decoder.Decode(&result1)
	if err != nil {
		t.Logf("First Decode: %v", err)
	}

	// Try decoding again (might fail or reuse decoder)
	var result2 struct {
		Name  []string `csv:"name"`
		Value []int    `csv:"value"`
	}

	err = decoder.Decode(&result2)
	if err != nil {
		t.Logf("Second Decode: %v", err)
	}
}
