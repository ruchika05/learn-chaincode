
package main

import (
	"encoding/json"
	"errors"
	"fmt"


	"github.com/hyperledger/fabric/core/chaincode/shim"
)


// Loan structure
type Loan struct {
	LoanID       string  `json:"loanID"`
	BuyerName    string  `json:"buyerName"`
	PropertyName string  `json:"propertyName"`
	BankName	 string  `json:"bankName"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


func (t *SimpleChaincode) createLoan(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating loan document")
	var loan Loan
	
	loan = Loan{LoanID:"1234",BuyerName:"ruchika",PropertyName:"HDIL",BankName:"HSBC"}
	loanBytes, err := json.Marshal(&loan)
	
	if err != nil {
			fmt.Println("error creating loan doc" + loan.LoanID)
			//return nil, errors.New("Error creating loan doc " + loan.LoanID)
		}
	err = stub.PutState(loan.LoanID, loanBytes)
	//return nil, nil
	return nil, nil

}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Init firing. Function will be ignored: " + function)

	// Initialize the collection of commercial paper keys
	/*fmt.Println("Initializing paper keys collection")
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	err := stub.PutState("PaperKeys", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize paper key collection")
	}*/

	fmt.Println("Initialization complete")
	return nil, nil
}




func GetLoan(loanid string, stub shim.ChaincodeStubInterface) (Loan, error) {
	var loan Loan

	loanBytes, err := stub.GetState(loanid)
	if err != nil {
		fmt.Println("Error retrieving cp " + loanid)
		return loan, errors.New("Error retrieving cp " + loanid)
	}

	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + loanid)
		return loan, errors.New("Error unmarshalling cp " + loanid)
	}

	return loan, nil
}


func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query running. Function: " + function)

    if function == "GetLoan" {
		fmt.Println("Getting particular cp")
		cp, err := GetLoan(args[0], stub)
		if err != nil {
			fmt.Println("Error Getting particular cp")
			return nil, err
		} else {
			cpBytes, err1 := json.Marshal(&cp)
			if err1 != nil {
				fmt.Println("Error marshalling the cp")
				return nil, err1
			}
			fmt.Println("All success, returning the cp")
			return cpBytes, nil
		}
	} else {
		fmt.Println("Generic Query call")
		bytes, err := stub.GetState(args[0])

		if err != nil {
			fmt.Println("Some error happenend: " + err.Error())
			return nil, err
		}

		fmt.Println("All success, returning from generic")
		return bytes, nil
	}
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke running. Function: " + function)

    if function == "createLoan" {
		return t.createLoan(stub, args)
	}

	return nil, errors.New("Received unknown function invocation: " + function)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting Simple chaincode: %s", err)
	}
}

