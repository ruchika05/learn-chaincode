package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// TestMyChaincode test my code
func TestMyChaincode(t *testing.T) {
	testStub := shim.NewMockStub("mock", new(MyChaincode))

	if testStub == nil {
		t.Fatalf("Unable to instantiate mockstub")
	}
	testStub.MockInit("t123", "init", nil)
	var objectBlob = `{"id": "1234", "name":"Pencils","qty":1000,"price":100}`
	testStub.MockInvoke("t123", "addObject", []string{objectBlob})
	testStub.MockInvoke("t123", "addObject", []string{"1", objectBlob})
	testStub.MockQuery("getObject", []string{"123"})
	testStub.MockQuery("getAllObjects", []string{})
}
