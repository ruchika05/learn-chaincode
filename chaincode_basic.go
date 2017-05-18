package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Object details
type Object struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"qty"`
	Price    float32 `json:"price"`
}

// ListOfObjects to store all objects
type ListOfObjects map[string]string

// MyChaincode function
type MyChaincode struct {
}

func getListOfObjects(shim shim.ChaincodeStubInterface) (map[string]string, error) {
	var err error
	var bytesRead []byte
	var list map[string]string

	bytesRead, err = shim.GetState("ListOfObjects")
	if err != nil {
		fmt.Println("Unable to get the list of Objects")
		return nil, err
	}
	if len(bytesRead) > 1 {
		fmt.Println("List of Objects exists, return the same")
		err = json.Unmarshal(bytesRead, &list)
		if err != nil {
			fmt.Println("Unable to get the list of Objects")
			return nil, err
		}
	} else {
		list = make(map[string]string)
		bytesRead, err = json.Marshal(&list)
		if err != nil {
			fmt.Println("Unable to get the list of Objects")
			return nil, err
		}
		err = shim.PutState("ListOfObjects", bytesRead)
		if err != nil {
			fmt.Println("Unable to get the list of Objects")
			return nil, err
		}
	}
	fmt.Println("returning the list of objects")
	return list, nil

}

func setListOfObjects(shim shim.ChaincodeStubInterface, list map[string]string) error {
	var err error
	var bytesRead []byte

	bytesRead, err = json.Marshal(&list)
	if err != nil {
		fmt.Println("Unable to update the list of Objects")
		return err
	}
	err = shim.PutState("ListofObjects", bytesRead)
	if err != nil {
		fmt.Println("Unable to update the list of Objects")
		return err
	}
	fmt.Println("updated the list of objects")
	return nil

}

func addObject(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	//var bytesRead []byte

	var obj Object

	if len(args) != 1 {
		fmt.Println("addObject called with incorrect number of arguments")
		return nil, errors.New("addObject called with incorrect number of arguments")
	}
	fmt.Printf("addObject called with args : %v\n", args[0])

	err = json.Unmarshal([]byte(args[0]), &obj)

	if err != nil {
		fmt.Printf("err : %v\n", err)
	}

	err = stub.PutState(obj.ID, []byte(args[0]))

	if err != nil {
		fmt.Printf("err : %v\n", err)
	}

	fmt.Printf("addObject called with obj : %v\n", obj)

	return nil, nil

}
func removeObject(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		fmt.Println("removeObject called with incorrect number of arguments")
		return nil, errors.New("removeObject called with incorrect number of arguments")
	}
	fmt.Printf("removeObject called with args : %v\n", args[0])
	return nil, nil

}
func updateObject(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		fmt.Println("updateObject called with incorrect number of arguments")
		return nil, errors.New("updateObject called with incorrect number of arguments")
	}
	fmt.Printf("updateObject called with args : %v\n", args[0])
	return nil, nil

}

func getObject(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var bytesRead []byte

	if len(args) != 1 {
		fmt.Println("getObject called with incorrect number of arguments")
		return nil, errors.New("getObject called with incorrect number of arguments")
	}
	fmt.Printf("getObject called with args : %v\n", args[0])

	bytesRead, err = stub.GetState(args[0])
	if err != nil {
		fmt.Printf("err : %v\n", err)
		return nil, err
	}
	return bytesRead, nil

}

func getAllObjects(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var isOk bool
	val, err := stub.ReadCertAttribute("position")
	fmt.Printf("Position => %v error %v \n", string(val), err)
	isOk, _ = stub.VerifyAttribute("position", []byte("Software Engineer"))
	if isOk {
		fmt.Printf("am ok")

	}

	if len(args) != 0 {
		fmt.Println("getAllObjects called with incorrect number of arguments")
		return nil, errors.New("getAllObjects called with incorrect number of arguments")
	}
	fmt.Printf("getAllObjects called\n")
	return nil, nil

}

// Init function
func (t *MyChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Initiliazing the chaincode")
	return nil, nil
}

// Invoke function
func (t *MyChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke called for function: " + function)
	fmt.Printf("args: %s\n", args)
	if function == "addObject" {
		return addObject(stub, args)
	} else if function == "removeObject" {
		return removeObject(stub, args)
	} else if function == "updateObject" {
		return updateObject(stub, args)
	}
	return nil, nil
}

// Query function
func (t *MyChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query called for function: " + function)
	fmt.Printf("args: %s\n", args)
	if function == "getObject" {
		return getObject(stub, args)
	} else if function == "getAllObjects" {
		return getAllObjects(stub, args)
	}
	return nil, nil
}

func main() {
	err := shim.Start(new(MyChaincode))
	if err != nil {
		fmt.Printf("Error starting MyChaincode: %s", err)
	}
}
