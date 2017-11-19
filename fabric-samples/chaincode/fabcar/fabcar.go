/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)
var logger = shim.NewLogger("treelogger")


// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Tree structure, with 4 properties.  Structure tags are used by encoding/json library
type Car struct {
	Types   string `json:"types"`
	HouseAddress string `json:"address"`
	OwnerName string `json:"owner"`
	QuantiTREE string `json:"quantity"`
	CreditByGov string `json:"credit"`
	Supplier string `json:"company"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryHouseAddress" { // queryCar used for getting data about specific car based on car #
														// now get based on HouseAddress so calling it queryHouseAddress
		return s.queryHouseAddress(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createEntry" { // createCar used for creating a new entry
												// now we want to createEntry containing all the fields in the Tree struct
		return s.createEntry(APIstub, args)
	} else if function == "queryAllOwners" { // queryAllCars is now queryAllOwners
		return s.queryAllOwners(APIstub)
	} else if function == "changeTreeQuantity" { // changeCarOwner is now changeTreeQuantity
		return s.changeTreeQuantity(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryHouseAddress(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	entryAsBytes, err := APIstub.GetState(args[0])
		if err != nil {
			return shim.Error(err.Error())
		}
	fmt.Println(entryAsBytes)
	return shim.Success(entryAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response { // hard code data here!!
	trees := []Car{
		Car{Types: "Blue Spruce", HouseAddress: "2084 Mountbatten Place", OwnerName: "xr45h2", QuantiTREE: "6", CreditByGov: "$22.33", Supplier: "Tree Inc."},
		Car{Types: "Red Maple", HouseAddress: "3456 Oxford Street", OwnerName: "3yh654", QuantiTREE: "4", CreditByGov: "$17.66", Supplier: "Tree Inc."},
		Car{Types: "Cucumber Tree", HouseAddress: "837 Wharncliffe Road", OwnerName: "ij2c5g", QuantiTREE: "11", CreditByGov: "$37.12", Supplier: "Tree Inc."},
		Car{Types: "Yellow Birch", HouseAddress: "52 Irwin Street", OwnerName: "sw3hg6", QuantiTREE: "54", CreditByGov: "$94.09", Supplier: "Tree Inc."},
	}

	i := 0
	for i < len(trees) {
		// fmt.Println("i is ", i)
		logger.Debugf("i is %d", i)
		entryAsBytes, _ := json.Marshal(trees[i])
		APIstub.PutState("CAR"+strconv.Itoa(i), entryAsBytes) // maybe not right
		fmt.Printf("Added %s", "CAR"+strconv.Itoa(i))
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createEntry(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	// Type: [""], HouseAddress: "", OwnerName: "", QuantiTREE: "", CreditByGov: "", Supplier: ""},
	var entry = Car{Types: args[1], HouseAddress: args[2], OwnerName: args[3], QuantiTREE: args[4], CreditByGov: args[5], Supplier: args[6]}
	logger.Debug("HelloCreating \n")
	entryAsBytes, _ := json.Marshal(entry)
	APIstub.PutState(args[0], entryAsBytes)
	logger.Debugf("put state %s \n", args[0])
	return shim.Success(nil)
}

func (s *SmartContract) queryAllOwners(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CAR0" // specify the range 			// worried about this ???? I feel like this is specified somewhere else
	endKey := "CAR999" // of allowed keys we can query
	logger.Debug("Hello \n")

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error(err.Error())
	}

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllOwners:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeTreeQuantity(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	entryAsBytes, _ := APIstub.GetState(args[0])
	entry := Car{}

	json.Unmarshal(entryAsBytes, &entry)
	entry.QuantiTREE = args[1]

	entryAsBytes, _ = json.Marshal(entry)
	APIstub.PutState(args[0], entryAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}

}
