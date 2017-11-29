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
type InformationChainCode struct {

}

//response request from the student to get approved or not by university admin
type User struct{
    ProfilePic string `json:"profilePic"`
    Name string `json:"username"`
    Email string`jsin:"email"`
    DateOfBirth time.Time `json:"dateOfBirth`
    Gender string `json:"gender"`
    Profession string `json:"profession"`
    Education []Education `json:"education"`
    Experience []Profession `json:"experience"`
}
type Education struct{
    Degree  string `json:"degree"`
    Board string `json:"board"`
    Institute string `json:"Institute`
    YearOfPassout int `json:"yearOfPassout"`
    Score   float64 `json:"score"`
}
type Profession struct{
    Organisation string `json:"organisation"`
    Designation string `json:"designation"`
    Location string `json:"location`
    DateOfJoining time.Time `json:"dateOfJoining"`
    DateOfRelieving time.Time `json:"dateOfRelieving"`
    Experience float64 `json:experience`
    StillWork bool `json:stilWorking`
}



var logger = shim.NewLogger("Informationa_Chaincode _CC")

func main() {
    err := shim.Start(new(InformationChainCode))
    if err != nil {
        fmt.Printf("Error while starting Information Chaincode - %s", err)
    }
}

func(t *InformationChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {

		logger.Info("########### Informationa Chaincode Init ###########")
     	fmt.Println("The Network as Been started by Information smartcontract")
    	fmt.Println("Ready To Take Approval and Requests for digital Expernice ")
    	return shim.Success(nil)
}

func (t *InformationChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(" Invoke ");
	function,args := stub.GetFunctionAndParameters()
	if function == "init" {
		return t.Init(stub)	
	}else if function == "register" {
		return register(stub, args)
	}else if function =="AddEducation"{
		return AddEducation(stub,args)
	}else if function=="AddExperience"{
		return AddExperience(stub,args)
	}else if function=="UpdateEducation"{
        return UpdateEducation(stub,args)
    }else if function =="UpdateExperience"{
       return UpdateExperience(stub,args)    
    }
    
	
    logger.Errorf("Received unknown invoke function name -%s",function)
	return shim.Error("Received unknown invoke function name -'" + function + "'")
}



func register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    fmt.Println("Registration of student in process... ")
    var education []Education
    var employee []Profession
    if len(args) !=7{
        return shim.Error("Incorrect number of arguments. Expecting 7. ID followed by details")
    }
    id :=args[0]
    displayPicture :=args[1]
    name :=args[2]
    dob,err := time.Parse("2006-01-02", args[3])
    if err != nil {
        return shim.Error("3st argument must be a date string and should be in 2006-01-02 format")
    }
    gender :=strings.ToLower(args[4])
    if gender!="male" && gender!="female"{
        return shim.Error("Gender must be male or female in 4rd argument")  
    }
    email := args[5]
    profession := args[6]
    
    // To check if Id already exists
    DetailsAsBytes, err := stub.GetState(id)
    if err != nil {
        fmt.Println("Registration Failed")
        return shim.Error("Failed to get ID details: " + err.Error())
    } else if DetailsAsBytes != nil {
        fmt.Println("This Id already exists: "+args[0])
        fmt.Println("Registration Failed")
        return shim.Error("This Id already exists: " + args[0])
    }
    //Assinging to student json
    //education := Education{degree,board,institute,passOut,score}
    registerDetails := User{displayPicture,name,email,dob,gender,profession,education,employee}
    DetailsJSONasBytes, err := json.Marshal(registerDetails)
    if err != nil {
        fmt.Println("error while Json marshal")
        return shim.Error(err.Error())
    }
    //write the variable into the ledger
    err = stub.PutState(id, DetailsJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }
    var str []string
    str = append(str,"1")
    //getHistory(stub,str)
    fmt.Println("Registration Successfull ---> Record created " )
    return shim.Success(nil)
}

////////////Adding Eduction details function for  particular student  ////////////////


 func  AddEducation(stub shim.ChaincodeStubInterface,args []string)pb.Response{

      
       if len(args) <6{        
          return shim.Error("Incorrect number of argument passed for adding eduction")
                      }
      
              id:=args[0]
             degree :=args[1]
             board:=args[2]
             institute:=args[3] 
             yop,err:=strconv.Atoi(args[4]) 
            score,err:=strconv.ParseFloat(args[5],64)
            eductionaldeatils:=Education{degree,board,institute,yop,score}
            DetailsAsBytes, err := stub.GetState(id)
            if err != nil {
              return shim.Error("Failed to get ID details: " + err.Error())
                        }



            informationdetails:=User{}
            json.Unmarshal([]byte(DetailsAsBytes),&informationdetails)
            var education []Education
            initialDetails :=informationdetails


            for i:=0; i<len(informationdetails.Education);i++{
              
              education = append(education,initialDetails.Education[i])
              
                         }
                

            education=append(education,eductionaldeatils)
        

        updatedeatils:=User{informationdetails.ProfilePic,informationdetails.Name,informationdetails.Email,
            informationdetails.DateOfBirth,informationdetails.Gender,informationdetails.Profession,education,informationdetails.Experience}
                  
            detailseducationmarhal,err:=json.Marshal(updatedeatils)
         if err!=nil {
         fmt.Println("error occured while converting to json")
         return shim.Error(err.Error())
                   }
              err=stub.PutState(id,detailseducationmarhal)
                   if err!=nil{
               return shim.Error(err.Error())
                 }
         fmt.Println("you Successfull  add eduction detals  for ....  :"+args[0])
          return shim.Success(nil)
 }

func  UpdateEducation(stub shim.ChaincodeStubInterface,args []string)pb.Response{

      
       if len(args) <6{        
          return shim.Error("Incorrect number of argument passed for adding eduction")
                      }
      
              id:=args[0]
             degree :=args[1]
             board:=args[2]
             institute:=args[3] 
             yop,err:=strconv.Atoi(args[4]) 
            score,err:=strconv.ParseFloat(args[5],64)
            eductionaldeatils:=Education{degree,board,institute,yop,score}
            DetailsAsBytes, err := stub.GetState(id)
            if err != nil {
              return shim.Error("Failed to get ID details: " + err.Error())
                        }



            informationdetails:=User{}
            json.Unmarshal([]byte(DetailsAsBytes),&informationdetails)
            var education []Education
            initialDetails :=informationdetails


            for i:=0; i<len(informationdetails.Education);i++{
              
                   if initialDetails.Education[i].Degree == degree{
                 education=append(education,eductionaldeatils)
                  }else{
                education = append(education,initialDetails.Education[i])
                           }

              
                         }        

        updatedeatils:=User{informationdetails.ProfilePic,informationdetails.Name,informationdetails.Email,
            informationdetails.DateOfBirth,informationdetails.Gender,informationdetails.Profession,education,informationdetails.Experience}
                  
            detailseducationmarhal,err:=json.Marshal(updatedeatils)
         if err!=nil {
         fmt.Println("error occured while converting to json")
         return shim.Error(err.Error())
                   }
              err=stub.PutState(id,detailseducationmarhal)
                   if err!=nil{
               return shim.Error(err.Error())
                 }
         fmt.Println("you Successfull  add eduction detals  for ....  :"+args[0])
          return shim.Success(nil)
 }


func  AddExperience(stub shim.ChaincodeStubInterface,args []string)pb.Response{

      
       if len(args) <8{        
          return shim.Error("Incorrect number of argument passed for adding eduction")
                      }
         id:=args[0]
        organisation:=args[1] 
        designation:=args[2] 
        location:=args[3] 
        doj,err := time.Parse("2006-01-02", args[4])
        if err != nil {
        return shim.Error("4st argument must be a date string and should be in 2006-01-02 format")
           }

        dor,err := time.Parse("2006-01-02", args[5])
        if err != nil {
         return shim.Error("5st argument must be a date string and should be in 2006-01-02 format")
               }
          experience,err:=strconv.ParseFloat(args[6],64)
           stillwork, err := strconv.ParseBool(args[7]); 

            if err != nil {
        
            return shim.Error("parse a bool value is gone error")
         
                              }
          Employee:=Profession{organisation,designation,location,doj,dor,experience,stillwork}
            DetailsAsBytes, err := stub.GetState(id)
            if err != nil {
              return shim.Error("Failed to get ID details: " + err.Error())
                        }



            informationdetails:=User{}
            json.Unmarshal([]byte(DetailsAsBytes),&informationdetails)
            var employee []Profession
            initialDetails :=informationdetails


            for i:=0; i<len(informationdetails.Experience);i++{
              
              employee = append(employee,initialDetails.Experience[i])
              
                         }
                

            employee=append(employee,Employee)
        

        updatedeatils:=User{informationdetails.ProfilePic,informationdetails.Name,informationdetails.Email,
            informationdetails.DateOfBirth,informationdetails.Gender,informationdetails.Profession,informationdetails.Education,employee}
                  
            detailsexperiencemarhal,err:=json.Marshal(updatedeatils)
         if err!=nil {
         fmt.Println("error occured while converting to json")
         return shim.Error(err.Error())
                   }
              err=stub.PutState(id,detailsexperiencemarhal)
                   if err!=nil{
               return shim.Error(err.Error())
                 }
         fmt.Println("you Successfull  add eduction detals  for ....  :"+args[0])
          return shim.Success(nil)
 }

func  UpdateExperience(stub shim.ChaincodeStubInterface,args []string)pb.Response{

      
       if len(args) <8{        
          return shim.Error("Incorrect number of argument passed for adding eduction")
                      }
         id:=args[0]
        organisation:=args[1] 
        designation:=args[2] 
        location:=args[3] 
        doj,err := time.Parse("2006-01-02", args[4])
        if err != nil {
        return shim.Error("4st argument must be a date string and should be in 2006-01-02 format")
           }

        dor,err := time.Parse("2006-01-02", args[5])
        if err != nil {
         return shim.Error("5st argument must be a date string and should be in 2006-01-02 format")
               }
          experience,err:=strconv.ParseFloat(args[6],64)
           stillwork, err := strconv.ParseBool(args[7]); 

            if err != nil {
        
            return shim.Error("parse a bool value is gone error")
         
                              }
          Employee:=Profession{organisation,designation,location,doj,dor,experience,stillwork}
            DetailsAsBytes, err := stub.GetState(id)
            if err != nil {
              return shim.Error("Failed to get ID details: " + err.Error())
                        }



            informationdetails:=User{}
            json.Unmarshal([]byte(DetailsAsBytes),&informationdetails)
            var employee []Profession
            initialDetails :=informationdetails


            for i:=0; i<len(informationdetails.Experience);i++{
                        
             if initialDetails.Experience[i].Organisation== organisation{
                    employee=append(employee,Employee)
                    }else{
              employee = append(employee,initialDetails.Experience[i])
                           }

              
                         }
                

          
        

        updatedeatils:=User{informationdetails.ProfilePic,informationdetails.Name,informationdetails.Email,
            informationdetails.DateOfBirth,informationdetails.Gender,informationdetails.Profession,informationdetails.Education,employee}
                  
            detailsexperiencemarhal,err:=json.Marshal(updatedeatils)
         if err!=nil {
         fmt.Println("error occured while converting to json")
         return shim.Error(err.Error())
                   }
              err=stub.PutState(id,detailsexperiencemarhal)
                   if err!=nil{
               return shim.Error(err.Error())
                 }
         fmt.Println("you Successfull  add eduction detals  for ....  :"+args[0])
          return shim.Success(nil)
 }



func getPersonalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    val,ok,err := cid.GetAttributeValue(stub,"accountType");
    if err != nil {
        shim.Error("There was an error trying to retrieve accountType attribute");
    }
    if !ok {
    shim.Error("The client identity does not possess accountType attribute");
    }
    queryFor := "";
    fmt.Printf("account Type : %v",val)
    
    val1,ok,err:= cid.GetAttributeValue(stub,"id");
    if err != nil {
        shim.Error("There was an error trying to retrieve accountType attribute");
    }
    if !ok {
        shim.Error("The client identity does not possess accountType attribute");
    }
    fmt.Printf("got request for Personal ID = %v",val1);
    queryFor = val1
    
    fmt.Printf("queryFor = %v",queryFor)
    info,err:= stub.GetState(queryFor)
    if err != nil {
        return shim.Error(err.Error())
    }
  var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten := false
    var studentrequest1 User
    json.Unmarshal(info, &studentrequest1)
    // Add a comma before array members, suppress it for the first array member
    if bArrayMemberAlreadyWritten == true {
        buffer.WriteString(",")
    }
        
    buffer.WriteString("{\"Details\":")
    // Record is a JSON object, so we write as-is
    buffer.WriteString(",")
    }
        
    buffer.WriteString("{\"Details\":")
    // Record is a JSON object, so we write as-is
    buffer.WriteString(string(info))
    buffer.WriteString("}")
    bArrayMemberAlreadyWritten = true
        //fmt.Printf("value is %s",string(queryResponse.Value))
    
    buffer.WriteString("]")
 fmt.Printf("- Personal Request:\n%s\n", buffer.String())
    return shim.Success(buffer.Bytes())
}


