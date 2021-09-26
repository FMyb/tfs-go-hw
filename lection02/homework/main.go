package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"
)

type Operation struct {
	Company   string      `json:"company"`
	ID        interface{} `json:"id"`
	Type      string      `json:"type"`
	Value     *int        `json:"value"`
	CreatedAt *time.Time  `json:"created_at"`
}

type CompanyOperation struct {
	Company           string        `json:"company"`
	OperationsCount   int           `json:"valid_operations_count"`
	Balance           int           `json:"balance"`
	InvalidOperations []interface{} `json:"invalid_operations,omitempty"`
}

func (op *Operation) UnmarshalJSON(data []byte) error {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	m := f.(map[string]interface{})
	operationMap, ok := m["operation"]
	if ok {
		v := operationMap.(map[string]interface{})
		id, ok := v["id"]
		if ok {
			switch a := id.(type) {
			case string:
				op.ID = a
			case int:
				op.ID = a
			case float64:
				if a == math.Trunc(a) {
					op.ID = int(a)
				}
			}
		}
		value, ok := v["value"]
		if ok {
			switch a := value.(type) {
			case string:
				parseVal, err := strconv.Atoi(a)
				if err == nil {
					op.Value = &parseVal
				}
			case float64:
				if a == math.Trunc(a) {
					parseVal := int(a)
					op.Value = &parseVal
				}
			case int:
				parseVal := value.(int)
				op.Value = &parseVal
			}
		}
		creationAt, ok := v["created_at"]
		if ok {
			create, _ := time.Parse(time.RFC3339, creationAt.(string))
			op.CreatedAt = &create
		}
		t, ok := v["type"]
		if ok {
			parseT, ok := t.(string)
			if ok {
				switch parseT {
				case "income", "outcome", "+", "-":
					op.Type = parseT
				}
			}
		}
	}
	id, ok := m["id"]
	if ok {
		switch a := id.(type) {
		case string:
			op.ID = a
		case int:
			op.ID = a
		case float64:
			if a == math.Trunc(a) {
				op.ID = int(a)
			}
		}
	}
	value, ok := m["value"]
	if ok {
		switch a := value.(type) {
		case string:
			parseVal, err := strconv.Atoi(a)
			if err == nil {
				op.Value = &parseVal
			}
		case float64:
			if a == math.Trunc(a) {
				parseVal := int(a)
				op.Value = &parseVal
			}
		case int:
			parseVal := a
			op.Value = &parseVal
		}
	}
	t, ok := m["type"]
	if ok {
		parseT, ok := t.(string)
		if ok {
			switch parseT {
			case "income", "outcome", "+", "-":
				op.Type = parseT
			}
		}
	}
	company, ok := m["company"]
	if ok {
		op.Company = company.(string)
	}
	creationAt, ok := m["created_at"]
	if ok {
		create, _ := time.Parse(time.RFC3339, creationAt.(string))
		op.CreatedAt = &create
	}
	return nil
}

func main() {
	var f2 *os.File
	defer f2.Close()
	var err error
	inputFile := flag.String("file", "", "path to json source")
	flag.Parse()
	if inputFile != nil && *inputFile != "" {
		f2, err = os.Open(*inputFile)
		if err != nil {
			fmt.Println("Can't open file as flag file")
			return
		}
	} else {
		inputFile := os.Getenv("FILE")
		if inputFile != "" {
			f2, err = os.Open(inputFile)
			if err != nil {
				fmt.Println("Can't open file as ENV FILE")
			}
		} else {
			f2 = os.Stdin
		}
	}
	data, err := ioutil.ReadAll(f2)
	if err != nil {
		fmt.Println(err)
		return
	}
	var operations []Operation
	err = json.Unmarshal(data, &operations)
	if err != nil {
		fmt.Println(err)
		return
	}
	operationsMap := make(map[string][]Operation)
	for _, operation := range operations {
		if operation.Company == "" {
			continue
		}
		operationList, ok := operationsMap[operation.Company]
		if !ok {
			operationList = make([]Operation, 0, 5)
		}
		operationsMap[operation.Company] = append(operationList, operation)
	}
	result := make([]CompanyOperation, 0, len(operationsMap))
	for company, operationList := range operationsMap {
		comp := CompanyOperation{Company: company}
		for _, operation := range operationList {
			if operation.CreatedAt == nil || operation.ID == nil {
				continue
			}
			if operation.Type == "" || operation.Value == nil {
				comp.InvalidOperations = append(comp.InvalidOperations, operation.ID)
			} else {
				comp.OperationsCount++
				switch operation.Type {
				case "+", "income":
					comp.Balance += *operation.Value
				case "-", "outcome":
					comp.Balance -= *operation.Value
				}
			}
		}
		result = append(result, comp)
	}
	str, err := json.MarshalIndent(result, "", "	")
	if err != nil {
		fmt.Println(err)
		return
	}
	f, _ := os.Create("lection02/homework/out.json")
	defer f.Close()
	_, err = fmt.Fprint(f, string(str))
	if err != nil {
		fmt.Println(err)
		return
	}
}
