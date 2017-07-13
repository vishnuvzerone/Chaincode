
/*
Copyright IBM Corp 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const _propertyID = "property1"
const _registerBlockID = "willRegister"
const _depatmentBlockID = "property1"
const _adminUserName = "admin"
const _adminPassword = "admin123"
const _departmentUserName = "department1"
const _departmentPassword = "department123"

type User struct {
	Name            string `json:"name"`
	Password            string `json:"password"`
	PropertyID         string `json:"propertyid"`
	Info         string `json:"info"`
}

type WillPaper struct {
	ID            string `json:"id"`
	VisibleInfo            string `json:"visibleinfo"`
	HiddenInfo         string `json:"hiddeninfo"`
	IsLocked         bool `json:"islocked"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState(_registerBlockID, []byte(args[0]))
	err = stub.PutState(_depatmentBlockID, []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "adduser" {
		return t.AddUser(stub, args)
	} else if function == "createwillpaper" {
		return t.CreateWillPaper(stub, args)
	} else if function == "unlockthewillbyadmin" {
		return t.UnlockTheWillByAdmin(stub, args)
	} 

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}else if function == "readusername" {
		return t.ReadUserName(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	var valAsbytes []byte
	fmt.Println("running write()")
	valAsbytes, err = json.Marshal(args[1])
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
	return valAsbytes, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// 
// AddUser - Add New User
// 
func (t *SimpleChaincode) AddUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, name, password, propertyid, info string
	var err error
	var valAsbytes []byte
	fmt.Println("running write()")
	
	var usr User

	name         = "\"Name\":\""+args[1]+"\", "							// Variables to define the JSON
	password         = "\"Password\":\""+args[2]+"\", "
	propertyid      = "\"PropertyID\":\""+_propertyID+"\", "
	info	= "\"PropertyID\":\""+args[3]+"\" "

	user_json := "{"+name+password+propertyid+info+"}" 
	err = json.Unmarshal([]byte(user_json), &usr)	
	valAsbytes, err = json.Marshal(usr)
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	//value = args[1]
	err = stub.PutState(key, valAsbytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
	//return valAsbytes, nil
}

// 
// ReadUserName - Read the username from the block
// 
func (t *SimpleChaincode) ReadUserName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	
	var usr User
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	
	err = json.Unmarshal(valAsbytes, &usr)

	return []byte(usr.Name), nil
}

// 
// CreateWillPaper - Creating new will paper
// 
func (t *SimpleChaincode) CreateWillPaper(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var id, hiddeninfo, visibleinfo string
	var islocked bool
	var err error
	var valAsbytes []byte
  
  if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 arguments")
	}
  
	fmt.Println("start to create a will paper")
	
	var wPaper WillPaper

	id         = "\"ID\":\""+args[0]+"\", "							// Variables to define the JSON
	visibleinfo      = "\"VisibleInfo\":\""+args[1]+"\", "
	hiddeninfo         = "\"HiddenInfo\":\""+args[2]+"\", "
	islocked         = "\"IsLocked\": true"

	wPaper_json := "{"+id+hiddeninfo+visibleinfo+"}" 
	err = json.Unmarshal([]byte(wPaper_json), &wPaper)	
	valAsbytes, err = json.Marshal(wPaper)	

	//key = args[0] //rename for funsies
	//value = args[1]
	err = stub.PutState(_registerBlockID, valAsbytes) //write the will to iniitial register
  err = stub.PutState(_depatmentBlockID, valAsbytes) //Write the will to department register
	if err != nil {
		return nil, err
	}
	return nil, nil
	//return valAsbytes, nil
}


// 
// UnlockTheWillByAdmin - Creating new will paper
// 
func (t *SimpleChaincode) UnlockTheWillByAdmin(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var willID, adminUserName, adminPassword string
	var err error
	var valAsbytes []byte
  
  if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 arguments")
	}
  
	fmt.Println("start to create a will paper")
	
	var wPaper WillPaper

	willID = args[0]
	adminUserName = args[1]
	adminPassword = args[2]
	
	if adminUserName == _adminUserName && 
	adminPassword == _adminPassword		{
		valAsbytes, err := stub.GetState(_registerBlockID)
		err = json.Unmarshal(valAsbytes, &wPaper)
		if willID == wPaper.ID && wPaper.IsLocked == true{
			wPaper.IsLocked = false
			valAsbytes, err = json.Marshal(wPaper)	
			err = stub.PutState(_registerBlockID, valAsbytes) //Updating the will to iniitial register
  			err = stub.PutState(_depatmentBlockID, valAsbytes) //Updating the will to department register
		}
	}
	
	if err != nil {
		return nil, err
	}
	return []byte("successfully Unlocked"), nil
	//return valAsbytes, nil
}


