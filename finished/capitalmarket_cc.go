/*
Copyright 2017 IBM, Infosys Ltd.

Licensed under the Apache License, Version 2.0 (the "License")
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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var counterID = 10000

//TODO: how to generate the IDs

func generateID() (string, error) {
	counterID = counterID + 1
	return strconv.Itoa(counterID), nil
}

const (
	millisPerSecond     = int64(time.Second / time.Millisecond)
	nanosPerMillisecond = int64(time.Millisecond / time.Nanosecond)
)

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(msInt/millisPerSecond,
		(msInt%millisPerSecond)*nanosPerMillisecond), nil
}

//FIOrder is created for trade requests received by FI
type FIOrder struct {
	FIOrderID       string    `json:"fiOrderID"`       // auto-generated unique ID for the FI Order
	FIID            string    `json:"fiID"`            // Unique ID of the FI
	CustodianBankID string    `json:"custodianBankID"` // Unique ID of the Custodian Bank
	BrokerID        string    `json:"brokerID"`        // Unique ID of the broker
	AccountID       string    `json:"accountID"`       // Account ID of the FI
	Product         string    `json:"product"`         // name of the Product
	Status          string    `json:"status"`          // status of the Product
	CreationDate    time.Time `json:"creationDate"`    // date of creation of FIOrder
	StockID         string    `json:"stockID"`         // name of the Stock
	Quantity        int       `json:"quantity"`        // quantity of stock to be bought/sold
	Exchange        string    `json:"exchange"`        // name of exchange
	OrderValidity   string    `json:"orderValidity"`   // validity of the order
	OrderType       string    `json:"orderType"`       // type of Order
	LimitPrice      float32   `json:"limitPrice"`      // limit price
}

// TradeObject Details
type TradeObject struct {
	TradeObjectID    string    `json:"tradeObjectID"`   // auto-generated unique ID for the Trade TradeObject
	SettlementStatus string    `json:"settlemetStatus"` // status of the settlement
	OderTradeNumber  string    `json:"oderTradeNumber"` // OderTradeNumber
	SettlementDate   time.Time `json:"creationDate"`    // date of settlement
}

// Transaction details
type Transaction struct {
	TransactionID    string    `json:"transactionID"` // auto-generated unique ID for the Transaction
	AccountID        string    `json:"accountID"`     // account id of the FI
	StockID          string    `json:"stockID"`       // id of the stock
	Quanity          int       `json:"quantity"`      // quantity of stocks traded
	TransactionDate  time.Time `json:"txnDate"`       // date of Transaction
	TransactionType  string    `json:"txnType"`       // type of txn - debit/credit
	EffectiveBalance int       `json:"balance"`       // effective balance of stocks post transaction
}

// AllFIOrders has a list of all orders ==> AllFIOrders[FIOrderID] = FIOrder
var AllFIOrders map[string]FIOrder

// AllOrdersForFI stores the list of all orders for a FI ==> AllOrdersForFI[FIID] = []FIOrderID
var AllOrdersForFI map[string][]string

// AllOrdersForBroker has a list of all orders for a Broker ==> AllOrdersForBroker[BrokerID] = []FIOrderID
var AllOrdersForBroker map[string][]string

// AllTradeObjects has a list of trade objects ==> TradeObject[TradeObjectID] = TradeObject
var AllTradeObjects map[string]TradeObject

// ConfirmedToFIOrder ==> ConfirmedToFIOrder[ConfirmedOrdererdId] = FIOrderID
var ConfirmedToFIOrder map[string]string

// matched orders array
//var matchedOrderedArray []string

// TradeSettlementMap has a lits  ==>  TradeSettlementMap[TradeObjectID]=[]ConfirmedOrdererdId *** TO CHECK ****
var TradeSettlementMap map[string][]string

// ListOfTransactions ==> Transaction[TransactionID]=Transaction
var ListOfTransactions map[string]Transaction

// ListOfTransactionsForFI ==> ListOfTransactionsForFI[FIID]=(ListOfStocks[StockID]=[]TransactionID)  *** TO CONFIRM ***
var ListOfTransactionsForFI map[string]map[string][]string //or[]TransactionID

// CapitalMarketChainCode defined the chaincode for global mobile wallet
type CapitalMarketChainCode struct {
}

var err error
var bytesArray []byte

// Initialize the Trade Object Map
func initAllTradeObjects(stub shim.ChaincodeStubInterface) ([]byte, error) {

	bytesArray, err = stub.GetState("AllTradeObjects")
	if err != nil {
		fmt.Printf("Failed to initialize the AllTradeObjects for block chain :%v\n", err)
		return nil, err
	}
	if len(bytesArray) != 0 {
		fmt.Printf("All Trade Objects map exists.\n")
		err = json.Unmarshal(bytesArray, &AllTradeObjects)
		if err != nil {
			fmt.Printf("Failed to initialize the AllTradeObjects for block chain :%v\n", err)
			return nil, err
		}
	} else { // create a new map for AllTradeObjects
		fmt.Printf("All Trade Objects map does not exist. To be created. \n")
		AllTradeObjects = make(map[string]TradeObject)
		bytesArray, err = json.Marshal(&AllTradeObjects)
		if err != nil {
			fmt.Printf("Failed to initialize the AllTradeObjects for block chain :%v\n", err)
			return nil, err
		}
		err = stub.PutState("AllTradeObjects", bytesArray)
		if err != nil {
			fmt.Printf("Failed to initialize the AllTradeObjects for block chain :%v\n", err)
			return nil, err
		}
	}
	fmt.Printf("Initiliazed AllTradeObjects : %v\n", AllTradeObjects)
	return nil, err
}

// Initialize the AllOrdersForBroker Map
func initAllOrdersForBroker(stub shim.ChaincodeStubInterface) ([]byte, error) {

	bytesArray, err = stub.GetState("AllOrdersForBroker")

	if err != nil {
		fmt.Printf("Failed to initialize the AllOrdersForBroker for block chain :%v\n", err)
		return nil, err
	}
	if len(bytesArray) != 0 {
		fmt.Printf("AllOrdersForBroker map exists.\n")
		err = json.Unmarshal(bytesArray, &AllOrdersForBroker)
		if err != nil {
			fmt.Printf("Failed to initialize the AllOrdersForBroker for block chain :%v\n", err)
			return nil, err
		}
	} else { // create a new map for AllOrdersForBroker
		fmt.Printf("AllOrdersForBroker map does not exist. To be created.\n")
		AllOrdersForBroker = make(map[string][]string)
		bytesArray, err = json.Marshal(&AllOrdersForBroker)
		if err != nil {
			fmt.Printf("Failed to initialize the AllOrdersForBroker for block chain :%v\n", err)
			return nil, err
		}
		err = stub.PutState("AllOrdersForBroker", bytesArray)
		if err != nil {
			fmt.Printf("Failed to initialize the AllOrdersForBroker for block chain :%v\n", err)
			return nil, err
		}
	}
	fmt.Printf("Initiliazed AllOrdersForBroker : %v\n", AllOrdersForBroker)
	return nil, err
}

// Initialize the AllOrdersForFI Map
func initAllOrdersForFI(stub shim.ChaincodeStubInterface) ([]byte, error) {

	bytesArray, err = stub.GetState("AllOrdersForFI")
	if err != nil {
		fmt.Printf("Failed to initialize the AllOrdersForFI for block chain :%v\n", err)
		return nil, err
	}
	if len(bytesArray) != 0 {
		fmt.Printf("AllOrdersForFI map exists.\n")
		err = json.Unmarshal(bytesArray, &AllOrdersForFI)
		if err != nil {
			fmt.Printf("Failed to initialize the AllOrdersForFI for block chain :%v\n", err)
			return nil, err
		}
	} else { // create a new map for AllOrdersForFI
		fmt.Printf("AllOrdersForFI map does not exist. To be created")
		AllOrdersForFI = make(map[string][]string)
		bytesArray, err = json.Marshal(&AllOrdersForFI)
		if err != nil {
			fmt.Printf("Failed to initialize the AllOrdersForFI for block chain :%v\n", err)
			return nil, err
		}
		err = stub.PutState("AllOrdersForFI", bytesArray)
		if err != nil {
			fmt.Printf("Failed to initialize the AllOrdersForFI for block chain :%v\n", err)
			return nil, err
		}
	}
	fmt.Printf("Initiliazed AllOrdersForFI : %v\n", AllOrdersForFI)
	return nil, err
}

// Initialize the AllFIOrders Map
func initAllFIOrders(stub shim.ChaincodeStubInterface) ([]byte, error) {

	bytesArray, err = stub.GetState("AllFIOrders")
	if err != nil {
		fmt.Printf("Failed to initialize the AllFIOrders for block chain :%v\n", err)
		return nil, err
	}
	if len(bytesArray) != 0 {
		fmt.Printf("AllFIOrders map exists.\n")
		err = json.Unmarshal(bytesArray, &AllFIOrders)
		if err != nil {
			fmt.Printf("Failed to initialize the AllFIOrders for block chain :%v\n", err)
			return nil, err
		}
	} else { // create a new map for AllFIOrders
		fmt.Printf("AllFIOrders map does not exist. To be created\n")
		AllFIOrders = make(map[string]FIOrder)
		bytesArray, err = json.Marshal(&AllFIOrders)
		if err != nil {
			fmt.Printf("Failed to initialize the AllFIOrders for block chain :%v\n", err)
			return nil, err
		}
		err = stub.PutState("AllFIOrders", bytesArray)
		if err != nil {
			fmt.Printf("Failed to initialize the AllFIOrders for block chain :%v\n", err)
			return nil, err
		}
	}
	fmt.Printf("Initiliazed AllFIOrders : %v\n", AllFIOrders)
	return nil, err
}

// Init function
func (t *CapitalMarketChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	initAllFIOrders(stub)
	initAllOrdersForFI(stub)
	initAllOrdersForBroker(stub)
	initAllTradeObjects(stub)
	fmt.Println("Initialization complete")

	return nil, err
}

// add orders created by the FI
func (t *CapitalMarketChainCode) createOrdersByFI(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating all orders by FI")
	fmt.Printf("len args: %d\n", len(args))
	if len(args) != 2 {
		fmt.Printf("Incorrect number of arguments.\n")
		return nil, errors.New("Incorrect number of arguments")
	}
	fmt.Printf("args[0]: %v\n", args[0])
	fmt.Printf("args[1]: %v\n", args[1])

	var fiOrders []FIOrder
	var err error

	err = json.Unmarshal([]byte(args[1]), &fiOrders)
	if err != nil {
		fmt.Printf("Error unmarshalling fi orders data : %v\n", err)
		return nil, errors.New("Failed to create fi orders")
	}
	fmt.Printf("fi orders after unmarshal: %v\n", fiOrders)

	if len(fiOrders) > 0 {
		for _, fiOrder := range fiOrders {
			fiOrder.FIOrderID, _ = generateID()
			AllFIOrders[fiOrder.FIOrderID] = fiOrder
			AllOrdersForBroker[fiOrder.BrokerID] = append(AllOrdersForBroker[fiOrder.BrokerID], fiOrder.FIOrderID)
			AllOrdersForFI[fiOrder.FIID] = append(AllOrdersForFI[fiOrder.FIID], fiOrder.FIOrderID)
		}
		fmt.Printf("Orders created successfully \n")
		return nil, nil
	}
	return nil, errors.New("There are no orders available for the FI")

}

/*
	Returns the list of FIOrders for a FI based on status
*/
func getAllOrdersForFIBasedOnStatus(FIID string, Status string, stub shim.ChaincodeStubInterface) ([]FIOrder, error) {
	var fiOrderIDs []string
	var fiOrder FIOrder
	var ok bool
	var fiOrdersByStatus []FIOrder

	if fiOrderIDs, ok = AllOrdersForFI[FIID]; ok {
		fmt.Printf("fiOrders : %v\n", fiOrderIDs)
		for _, id := range fiOrderIDs {
			// get details of each FI Orders
			fmt.Printf("fiOrders ids : %v\n", id)
			if fiOrder, ok = AllFIOrders[id]; ok {
				if len(Status) > 0 {
					if fiOrder.Status == Status {
						fiOrdersByStatus = append(fiOrdersByStatus, fiOrder)
					}
				} else {
					fiOrdersByStatus = append(fiOrdersByStatus, fiOrder)
				}
			}
		}
		fmt.Printf("List Of Orders by FI %s : %v \n", FIID, fiOrdersByStatus)
		return fiOrdersByStatus, nil
	}
	return nil, errors.New("Unable to find any orders for FI")

}

/*
	Returns the list of FIOrders for a Broker
*/
func getAllOrdersForBrokerBasedOnStatus(BrokerID string, Status string, stub shim.ChaincodeStubInterface) ([]FIOrder, error) {
	var fiOrderIDs []string
	var fiOrdersByStatus []FIOrder
	var fiOrder FIOrder
	var ok bool

	if fiOrderIDs, ok = AllOrdersForBroker[BrokerID]; ok {
		fmt.Printf("fiOrders : %v\n", fiOrderIDs)
		for _, id := range fiOrderIDs {
			// get details of each FI Orders
			fmt.Printf("fiOrders ids : %v\n", id)
			if fiOrder, ok = AllFIOrders[id]; ok {
				if len(Status) > 0 {
					if fiOrder.Status == Status {
						fiOrdersByStatus = append(fiOrdersByStatus, fiOrder)
					}
				} else {
					fiOrdersByStatus = append(fiOrdersByStatus, fiOrder)
				}
			}
		}
		fmt.Printf("List Of Orders by FI %s : %v \n", BrokerID, fiOrdersByStatus)
		return fiOrdersByStatus, nil
	}
	return nil, errors.New("Unable to find any orders for Broker")
}

// Query function
func (t *CapitalMarketChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var allOrders []FIOrder
	var err error
	var allBytes []byte
	if function == "getAllOrdersForFIBasedOnStatus" {
		if len(args) != 2 {
			fmt.Printf("Incorrect number of arguments to call getAllOrdersForFIBasedOnStatus.\n")
			return nil, errors.New("Incorrect number of arguments")
		}
		allOrders, err = getAllOrdersForFIBasedOnStatus(args[0], args[1], stub)
		if err != nil {
			fmt.Printf("Error getting All Orders for FI %s : %v\n", args[0], err)
			return nil, err
		}
		allBytes, err := json.Marshal(&allOrders)
		if err != nil {
			fmt.Printf("Error unmarshalling all orders : %v\n", err)
			return nil, err
		}
		fmt.Printf("All orders for FI %s successfully read\n", args[0])
		return allBytes, nil
	} else if function == "getAllOrdersForBrokerBasedOnStatus" {
		if len(args) != 2 {
			fmt.Printf("Incorrect number of arguments.\n")
			return nil, errors.New("Incorrect number of arguments to call getAllOrdersForBrokerBasedOnStatus ")
		}
		allOrders, err = getAllOrdersForBrokerBasedOnStatus(args[0], args[1], stub)
		if err != nil {
			fmt.Printf("Error getting All Orders for Broker %s : %v", args[0], err)
			return nil, err
		}
		allBytes, err = json.Marshal(&allOrders)
		if err != nil {
			fmt.Printf("Error unmarshalling all orders : %v\n", err)
			return nil, err
		}
		fmt.Printf("All orders for Broker %s successfully read\n", args[0])
		return allBytes, nil

	} //else {
	fmt.Println("received unknown function call: ", function)
	//}
	return nil, nil
}

// Invoke function
func (t *CapitalMarketChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke running. Function: " + function)
	fmt.Printf("args: %s\n", args)

	if function == "createOrdersByFI" {
		return t.createOrdersByFI(stub, args)
	}
	return nil, errors.New("Received unknown function invocation: " + function)
}

func main() {
	err := shim.Start(new(CapitalMarketChainCode))
	if err != nil {
		fmt.Printf("Error starting Capital Market chaincode: %s\n", err)
	}
}
