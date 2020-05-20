
package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a Vik
type SmartContract struct {
	contractapi.Contract
}

// Vik describes basic details of what makes up a Vik
type Vik struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Vik
}

// InitLedger adds a base set of Viks to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	Viks := []Vik{
		Vik{Name: "Siva", Partner: "HP", Target: "$10,000", Achivement: "One", Budget: "$15,000", Expenses: "$3000", Description: "The proposal has to be completed before this yearend"},
		Vik{Name: "Ajay", Partner: "TCS", Target: "$15,000", Achivement: "Two", Budget: "$18,000", Expenses: "$2000", Description: "The proposal has to be completed before this yearend"},
		Vik{Name: "John", Partner: "Wipro", Target: "$20,000", Achivement: "Three", Budget: "$25,000", Expenses: "$3000", Description: "The proposal has to be completed before this yearend"},
		Vik{Name: "Peter", Partner: "Accenture", Target: "$10,000", Achivement: "Four", Budget: "$15,000", Expenses: "$3000", Description: "The proposal has to be completed before this yearend"},
		Vik{Name: "Rakesh", Partner: "Zentax", Target: "$10,000", Achivement: "Five", Budget: "$15,000", Expenses: "$3000", Description: "The proposal has to be completed before this yearend"},
		Vik{Name: "Ajay", Partner: "Verizon", Target: "$10,000", Achivement: "Six", Budget: "$15,000", Expenses: "$3000", Description: "The proposal has to be completed before this yearend"},
		Vik{Name: "Mahendra", Partner: "Cisco", Target: "$10,000", Achivement: "Seven", Budget: "$15,000", Expenses: "$3000", Description: "The proposal has to be completed before this yearend"},
		Vik{Name: "Smith", Partner: "Gigo", Target: "$10,000", Achivement: "Eight", Budget: "$15,000", Expenses: "$3000", Description: "The proposal has to be completed before this yearend"},
	}

	for i, Vik := range Viks {
		VikAsBytes, _ := json.Marshal(Vik)
		err := ctx.GetStub().PutState("Vik"+strconv.Itoa(i), VikAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateVik adds a new Vik to the world state with given details
func (s *SmartContract) CreateVik(ctx contractapi.TransactionContextInterface, VikNumber string, Name string, Partner string, Target string, Achivement string, Budget string, Expenses string, Description string) error {
	Vik := Vik{
		Name:   Name,
		Partner:  model,
		Target: colour,
		Achivement:  owner,
		Budget: Budget,
		Expenses: Expenses,
		Description: Description	
	}

	VikAsBytes, _ := json.Marshal(Vik)

	return ctx.GetStub().PutState(VikNumber, VikAsBytes)
}

// QueryVik returns the Vik stored in the world state with given id
func (s *SmartContract) QueryVik(ctx contractapi.TransactionContextInterface, VikNumber string) (*Vik, error) {
	VikAsBytes, err := ctx.GetStub().GetState(VikNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if VikAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", VikNumber)
	}

	Vik := new(Vik)
	_ = json.Unmarshal(VikAsBytes, Vik)

	return Vik, nil
}

// QueryAllViks returns all Viks found in world state
func (s *SmartContract) QueryAllViks(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "Vik0"
	endKey := "Vik99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		Vik := new(Vik)
		_ = json.Unmarshal(queryResponse.Value, Vik)

		queryResult := QueryResult{Key: queryResponse.Key, Record: Vik}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeVikOwner updates the owner field of Vik with given id in world state
func (s *SmartContract) ChangeVikOwner(ctx contractapi.TransactionContextInterface, VikNumber string, newOwner string) error {
	Vik, err := s.QueryVik(ctx, VikNumber)

	if err != nil {
		return err
	}

	Vik.Owner = newOwner

	VikAsBytes, _ := json.Marshal(Vik)

	return ctx.GetStub().PutState(VikNumber, VikAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabVik chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabVik chaincode: %s", err.Error())
	}
}
