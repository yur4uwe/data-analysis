package uncsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type FieldDecoder interface {
	DecodeCSV(string) error
}

type FieldEncoder interface {
	EncodeCSV() (string, error)
}

type Decoder struct {
	r *csv.Reader
}

func NewDecoder(r *csv.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

type Encoder struct {
	w *csv.Writer
}

// Assumes struct of arrays
func (p *Decoder) Decode(v any) error {
	header, err := p.r.Read()
	if err != nil {
		return err
	}

	// Strip BOM from first header column if present
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\uFEFF")
	}

	destT := reflect.TypeOf(v)
	if destT == nil {
		return fmt.Errorf("cannot decode nil")
	}
	if destT.Kind() != reflect.Pointer || destT.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("type mismatch, expected pointer to struct of slices|arrays got: %s", destT.Kind())
	}
	destT = destT.Elem()
	destV := reflect.ValueOf(v).Elem()

	columnNameToField := make(map[string]int)
	for i, name := range header {
		columnNameToField[name] = i
	}

	for rowIdx := 0; ; rowIdx++ {
		row, err := p.r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		for i := range destT.NumField() {
			fieldT := destT.Field(i)

			fieldKind := fieldT.Type.Kind()
			if fieldKind != reflect.Slice && fieldKind != reflect.Array {
				return fmt.Errorf("expected fields to be arrays|slices got %s", fieldT.Type.Kind())
			}

			fieldV := destV.Field(i)
			if rowIdx >= fieldV.Len() && fieldKind == reflect.Array {
				return fmt.Errorf("array field %s doesn't have enough space for row %d (length: %d)",
					fieldT.Name, rowIdx, fieldV.Len())
			}
			if fieldV.Kind() == reflect.Slice && fieldV.Cap() == 0 {
				fieldV.Set(reflect.MakeSlice(fieldV.Type(), 0, 128))
			}

			tag := fieldT.Tag.Get("csv")
			if tag == "" {
				continue
			}

			colIdx, ok := columnNameToField[tag]
			if !ok {
				return fmt.Errorf("field %s: column %s not found in CSV header", fieldT.Name, tag)
			}

			elemType := fieldT.Type.Elem()
			if reflect.PointerTo(elemType).Implements(reflect.TypeFor[FieldDecoder]()) {
				newElem := reflect.New(elemType)
				if decoder, ok := newElem.Interface().(FieldDecoder); ok {
					if err := decoder.DecodeCSV(row[colIdx]); err != nil {
						return fmt.Errorf(
							"field %s row %d: custom decode failed: %w",
							fieldT.Name, rowIdx, err,
						)
					}

					if err := setValueAtIndex(fieldV, rowIdx, newElem.Elem().Interface()); err != nil {
						return fmt.Errorf("field %s row %d: %w", fieldT.Name, rowIdx, err)
					}
					continue
				}
			}
			elemKind := elemType.Kind()

			elemKindSizeBits := getBitSizeFromKind(elemKind)
			switch elemKind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intVal, err := strconv.ParseInt(row[colIdx], 0, elemKindSizeBits)
				if err != nil {
					return fmt.Errorf(
						"failed to parse value %s as %d-bit integer for field %s: %w",
						row[colIdx], elemKindSizeBits, fieldT.Name, err,
					)
				}
				if err := setValueAtIndex(fieldV, rowIdx, intVal); err != nil {
					return fmt.Errorf("field %s row %d: %w", fieldT.Name, rowIdx, err)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				uintVal, err := strconv.ParseUint(row[colIdx], 0, elemKindSizeBits)
				if err != nil {
					return fmt.Errorf(
						"failed to parse value %s as %d-bit unsigned integer for field %s: %w",
						row[colIdx], elemKindSizeBits, fieldT.Name, err,
					)
				}
				if err := setValueAtIndex(fieldV, rowIdx, uintVal); err != nil {
					return fmt.Errorf("field %s row %d: %w", fieldT.Name, rowIdx, err)
				}
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(row[colIdx])
				if err != nil {
					return fmt.Errorf(
						"failed to parse value %s as a boolean for field %s: %w",
						row[colIdx], fieldT.Name, err,
					)
				}
				if err := setValueAtIndex(fieldV, rowIdx, boolVal); err != nil {
					return fmt.Errorf("field %s row %d: %w", fieldT.Name, rowIdx, err)
				}
			case reflect.Float32, reflect.Float64:
				floatVal, err := strconv.ParseFloat(row[colIdx], elemKindSizeBits)
				if err != nil {
					return fmt.Errorf(
						"failed to parse value %s as float%d for field %s: %w",
						row[colIdx], elemKindSizeBits, fieldT.Name, err,
					)
				}
				if err := setValueAtIndex(fieldV, rowIdx, floatVal); err != nil {
					return fmt.Errorf("field %s row %d: %w", fieldT.Name, rowIdx, err)
				}
			case reflect.String:
				if err := setValueAtIndex(fieldV, rowIdx, row[colIdx]); err != nil {
					return fmt.Errorf("field %s row %d: %w", fieldT.Name, rowIdx, err)
				}
			default:
				return fmt.Errorf("expected kind of element to be of simple type")
			}
		}

	}

	return nil
}

func (p *Encoder) Encode(v any) error {
	return nil
}

// Rename and return bits directly
func getBitSizeFromKind(kind reflect.Kind) int {
	switch kind {
	case reflect.Int8, reflect.Uint8:
		return 8
	case reflect.Int16, reflect.Uint16:
		return 16
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 32
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 64
	case reflect.Int, reflect.Uint:
		return int(unsafe.Sizeof(int(0)) * 8)
	default:
		return 0
	}
}

func setValueAtIndex(fieldValue reflect.Value, index int, value any) error {
	kind := fieldValue.Kind()

	switch kind {
	case reflect.Array:
		if index >= fieldValue.Len() {
			return fmt.Errorf("index %d out of bounds for array of length %d", index, fieldValue.Len())
		}
	case reflect.Slice:
		if index >= fieldValue.Cap() {
			newSlice := reflect.MakeSlice(fieldValue.Type(), index+1, (index+1)*2)
			reflect.Copy(newSlice, fieldValue)
			fieldValue.Set(newSlice)
		} else if index >= fieldValue.Len() {
			fieldValue.SetLen(index + 1)
		}
	default:
		return fmt.Errorf("cannot set index on kind %v", kind)
	}

	elem := fieldValue.Index(index)

	switch v := value.(type) {
	case int64:
		elem.SetInt(v)
	case uint64:
		elem.SetUint(v)
	case float64:
		elem.SetFloat(v)
	case bool:
		elem.SetBool(v)
	case string:
		elem.SetString(v)
	default:
		elem.Set(reflect.ValueOf(value))
	}

	return nil
}
