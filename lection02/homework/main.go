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
	Company   string
	ID        interface{}
	Type      string
	Value     *int
	CreatedAt *time.Time
}

type OperationBody struct {
	ID        interface{} `json:"id"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	CreatedAt string      `json:"created_at"`
}

type OperationParse struct {
	Company       string `json:"company"`
	OperationBody `json:""`
	Body          OperationBody `json:"operation"`
}

type CompanyOperation struct {
	Company           string        `json:"company"`
	OperationsCount   int           `json:"valid_operations_count"`
	Balance           int           `json:"balance"`
	InvalidOperations []interface{} `json:"invalid_operations,omitempty"`
}

func convertToOperations(opParse []OperationParse) []Operation {
	res := make([]Operation, len(opParse))
	for i, op := range opParse {
		res[i] = op.convertToOperation()
	}
	return res
}

func (opParse OperationParse) convertToOperation() Operation {
	op := Operation{}
	op.Company = opParse.Company
	id := opParse.ID
	if id == nil {
		id = opParse.Body.ID
	}
	if id != nil {
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
	value := opParse.Value
	if value == nil {
		value = opParse.Body.Value
	}
	if value != nil {
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
	createdAt := opParse.CreatedAt
	if createdAt == "" {
		createdAt = opParse.Body.CreatedAt
	}
	if createdAt != "" {
		create, _ := time.Parse(time.RFC3339, createdAt)
		op.CreatedAt = &create
	}
	t := opParse.Type
	if t == "" {
		t = opParse.Body.Type
	}
	if t != "" {
		switch t {
		case "income", "outcome", "+", "-":
			op.Type = t
		}
	}
	return op
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
	var operationsParse []OperationParse
	err = json.Unmarshal(data, &operationsParse)
	if err != nil {
		fmt.Println(err)
		return
	}
	operations := convertToOperations(operationsParse)
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
