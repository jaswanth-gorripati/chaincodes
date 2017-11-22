/*
Copyright IBM Corp. 2016 All Rights Reserved.

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
	 "fmt"
	 "bytes"
     "encoding/json"
     "strconv"
     "time"
     "strings"


	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//simple chaincode for univeristy 
type Academicchaincode struct{

}

//response request from the student to get approved or not by university admin
type studentrequest struct{

	StudentId string `json:"studentId"`
	Degree string `json:"degree"`
	Percentage  float64 `json:marks"`
	YearOfPassOut string `json:"yop"`
	CollegeName string `json:"clgname"`
	Location string `json:"location"`
	University_Board string `json:"board"`
	Timeregistred time.Time `json:"timeenrolled"`
	Status string `json:"status"`  
	Remarks string    `json:"remarks"`

}


var logger = shim.NewLogger("university_CC")

func main() {
	   err := shim.Start(new(Academicchaincode))
	   if err != nil {
	   logger.Errorf("Error while Initializing Academic Chaincode - %s",err)
	    }
}

func(t *Academicchaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

		logger.Info("########### Academic Chaincode Init ###########")
     	fmt.Println("The Network as Been started by University smartcontract")
    	fmt.Println("Ready For Approval and Requests of Degree")
    	return shim.Success(nil)
}

func (t *Academicchaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(" Invoke ");
	function,args := stub.GetFunctionAndParameters()
	if function == "init" {
		return t.Init(stub)	
	}else if function == "RequestEnroll" {
		return RequestEnroll(stub, args);
	}else if function =="getdetailsbyboard"{
		return getdetailsbyboard(stub,args)
	}else if function=="hypernymprocess"{
		return hypernymprocess(stub,args)
	}
	
    logger.Errorf("Received unknown invoke function name -%s",function)
	return shim.Error("Received unknown invoke function name -'" + function + "'")
}


func  RequestEnroll(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args)!=7{
		return shim.Error("The Requested  Degree not full fill the all Requriments")
	}
	studentId:=args[0]
	degree  :=args[1]
	percentage,err :=strconv.ParseFloat(args[2],64)
	yop:=args[3]
	clgname:=args[4]
	localArea:= args[5]
	university_board:=args[6]
	stauts:="pending" 
	Current_date:=time.Now().Local()

	studentrequest:=studentrequest{studentId,degree,percentage,yop,clgname,localArea,university_board,Current_date,stauts,""}
  
    studentrequestmarshall,err:=json.Marshal(studentrequest)

	if err !=nil{
		logger.Errorf("error occured while converting to json")
		return shim.Error(err.Error())
	}


    err=stub.PutState(""+studentId+"_"+degree+"",studentrequestmarshall)
    if err!=nil{
    	logger.Errorf("error occured while updating  to ledger")
		return shim.Error(err.Error())
    }


     logger.Info("details are entered")
     return shim.Success(nil);
} 


////////////rich query function for getting particular student  ////////////////

func hypernymprocess(stub shim.ChaincodeStubInterface,args []string) pb.Response{

     
     if len(args)<3{
     	
     	return shim.Error("the arguments which passed or not upto the mark")
     }
    
    studentId:=args[0]
    degree:=args[1]
    status:=args[2]
    remarks:=args[3]

    var variable=""+studentId+"_"+degree+""
    //strconv.string(variable)
      studentrequestjson:=studentrequest{}

    newstudentrequest,err:=stub.GetState(variable)

   if err!=nil{
    	logger.Errorf("error occured while updating  to ledger")
		return shim.Error(err.Error())
    }
    json.Unmarshal([]byte(newstudentrequest),&studentrequestjson)


    Current_date:=time.Now().Local()
    studentupdaterespnose:=studentrequest{studentrequestjson.StudentId,studentrequestjson.Degree,studentrequestjson.Percentage,studentrequestjson.YearOfPassOut,
    	studentrequestjson.CollegeName,studentrequestjson.Location,studentrequestjson.University_Board,
    	Current_date,status,remarks}

    studentresponsemarshall,err:=json.Marshal(studentupdaterespnose)

	if err !=nil{
		logger.Errorf("error occured while converting to json")
		return shim.Error(err.Error())
	}


    err=stub.PutState(""+studentId+"_"+degree+"",studentresponsemarshall)
    if err!=nil{
    	logger.Errorf("error occured while updating  to ledger")
		return shim.Error(err.Error())
    }


     logger.Info("details are updated")
     return shim.Success(nil);
 }

func getdetailsbyboard(stub shim.ChaincodeStubInterface,args []string) pb.Response{

     if len(args)<1{

     	return shim.Error("the details did not get expecting one argument")
     }

    board_university:= strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"board\":\"%s\"}}}",board_university)
	
	resultsIterator,err:= stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
  	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
		return shim.Error(err.Error())
		}
	    var studentrequest1 studentrequest
		json.Unmarshal(queryResponse.Value, &studentrequest1)   
        // Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		 
		buffer.WriteString("{")
     	// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
		//fmt.Printf("value is %s",string(queryResponse.Value))
	}
	buffer.WriteString("]")
   	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes()) 
}


