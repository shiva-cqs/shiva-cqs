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
	ID   string `json:"ID"`
	Name  string `json:"Name"`
	Description string `json:"Description"`
	Stage  string `json:"Stage"`
	TargetRevenue  string `json:"TargetRevenue"`
	Oppurtunity  string `json:"Oppurtunity"`
	Partners  string `json:"Partners"`
	StartDate  string `json:"StartDate"`
	TeamMembers  string `json:"TeamMembers"`
}


// CreateVik adds a new Vik to the world state with given details
func (s *SmartContract) CreateVik(ctx contractapi.TransactionContextInterface, VikNumber string, ID string, Name string, Description string, Stage string, TargetRevenue string, Oppurtunity string, Partners string, StartDate string, TeamMembers string) error {
	Vik := Vik{
		ID: ID,
		Name: Name,
		Description: Description,
		Stage: Stage,
		TargetRevenue: TargetRevenue,
		Oppurtunity: Oppurtunity,
		Partners: Partners,
		StartDate: StartDate,
		TeamMembers: TeamMembers,
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

//-----------------------------------------------------------------------------------------------
// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Vik
}

// InitLedger adds a base set of Viks to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	Viks := []Vik{
		Vik{ID: "1", Name: "Siva", Description: "Solution 1", Stage: "One", TargetRevenue: "$15,000", Oppurtunity: "$3000", Partners: "", StartDate: "01/01/2021", TeamMembers: ""},
		Vik{ID: "2", Name: "Smith", Description: "Solution 2", Stage: "Two", TargetRevenue: "$15,000", Oppurtunity: "$3000", Partners: "", StartDate: "01/01/2021", TeamMembers: ""},
		Vik{ID: "3", Name: "Christo", Description: "Solution 3", Stage: "Three", TargetRevenue: "$15,000", Oppurtunity: "$3000", Partners: "", StartDate: "01/01/2021", TeamMembers: ""},
		Vik{ID: "4", Name: "Melanie", Description: "Solution 4", Stage: "Four", TargetRevenue: "$15,000", Oppurtunity: "$3000", Partners: "", StartDate: "01/01/2021", TeamMembers: ""},
		Vik{ID: "5", Name: "Yuvesh", Description: "Solution 5", Stage: "Three", TargetRevenue: "$15,000", Oppurtunity: "$3000", Partners: "", StartDate: "01/01/2021", TeamMembers: ""},
		Vik{ID: "6", Name: "Lymina", Description: "Solution 6", Stage: "Four", TargetRevenue: "$15,000", Oppurtunity: "$3000", Partners: "", StartDate: "01/01/2021", TeamMembers: ""}
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

//**********************************************************************************************
//Submit Smart Contract the Chaincode
peer lifecycle chaincode querycommitted --channelID mychannel --name fabcar --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

//**********************************************************************************************
//invoke the Chaincode

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls true --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n fabcar --peerAddresses
localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt --isInit -c '{"function":"InitLedger","Args":[]}'

//**********************************************************************************************
//Query the Chaincode
peer chaincode query -C mychannel -n fabcar -c '{"Args":["queryAllCars"]}'
//**********************************************************************************************
