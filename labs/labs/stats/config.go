package stats

import (
	"encoding/csv"
	"fmt"
	"io"
	"labs/labs/common"
	"math/rand"
	"os"
	"reflect"
	"strconv"
)

const (
	LabID = "5"
)

var (
	Config = common.LabConfig{
		Lab: Metadata,
		Charts: map[string]common.Chart{
			ErrorAnalysisChartID: ErrorAnalysisChart,
		},
	}

	Metadata = common.LabMetadata{
		ID:   LabID,
		Name: "Statistical Analysis",
		Charts: map[string]common.ChartMetadata{
			ErrorAnalysisChartID: ErrorAnalysisMeta,
		},
	}
)

// ...existing code...

func RandomSlice(stdev, mean float64, size int) []float64 {
	slice := make([]float64, size)

	for i := range size {
		slice[i] = stdev*rand.NormFloat64() + mean
	}

	return slice
}

type PositionType string

const (
	Programmer PositionType = "Programmer"
	Tester     PositionType = "Tester"
)

func IsValidPosition(pos string) bool {
	switch pos {
	case string(Programmer), string(Tester):
		return true
	default:
		return false
	}
}

type SalaryRecord struct {
	ID              []string       `csv:"Працівник"`
	Age             []int          `csv:"Вік"`
	Position        []PositionType `csv:"Посада"`
	ExperienceYears []int          `csv:"Досвід (років)"`
	Projects        []int          `csv:"Проектів завершено"`
	Salary          []float64      `csv:"Зарплата (USD)"`
}

func ReadSalaryCSV(filename string) (*SalaryRecord, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	r.Read() // Skip header

	rec := &SalaryRecord{
		ID:              []string{},
		Age:             []int{},
		Position:        []PositionType{},
		ExperienceYears: []int{},
		Projects:        []int{},
		Salary:          []float64{},
	}
	recT := reflect.TypeFor[SalaryRecord]()

	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		for i := 0; i < recT.NumField(); i++ {
			field := recT.Field(i)

			tag := field.Tag.Get("csv")
			if tag == "" {
				continue
			}

			value := row[i]

			switch field.Type.Kind() {
			case reflect.String:
				if IsValidPosition(value) && field.Name == "Position" {
					rec.Position = append(rec.Position, PositionType(value))
				} else if field.Name == "ID" {
					rec.ID = append(rec.ID, value)
				}
			case reflect.Int:
				intValue, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("invalid int value for %s: %v", field.Name, err)
				}
				switch field.Name {
				case "Age":
					rec.Age = append(rec.Age, intValue)
				case "ExperienceYears":
					rec.ExperienceYears = append(rec.ExperienceYears, intValue)
				case "Projects":
					rec.Projects = append(rec.Projects, intValue)
				}
			case reflect.Float64:
				floatValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid float value for %s: %v", field.Name, err)
				}
				rec.Salary = append(rec.Salary, floatValue)
			}
		}
	}

	recV := reflect.ValueOf(rec).Elem()

	for i := 0; i < recV.NumField(); i++ {
		field := recV.Field(i)
		sliceType := field.Type()
		clipped := reflect.MakeSlice(sliceType, field.Len(), field.Len())

		reflect.Copy(clipped, field)
	}

	return rec, nil
}
