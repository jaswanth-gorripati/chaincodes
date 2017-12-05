package main


import (
	"encoding/json"
	"fmt"
	"bytes"
	"strings"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type InformationChainCode struct {
}
var logger = shim.NewLogger("InformationChainCode")
func main() {
	err := shim.Start(new(InformationChainCode))
	if err != nil {
		fmt.Printf("Error while starting Information Chaincode - %s", err)
	}
}
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
	Degree	string `json:"degree"`
	Board string `json:"board"`
	Institute string  `json:"Institute`
	YearOfPassout int  `json:"yearOfPassout"`
	Score	float64 `json:"score"`
}
type Profession struct{
	Organisation string `json:"organisation"`
	Designation string `json:"designation"`
	Location string  `json:"location`
	DateOfJoining time.Time  `json:"dateOfJoining"`
	DateOfRelieving	time.Time `json:"dateOfRelieving"`
	Experience float64 `json:experience`
	StillWork bool `json:stilWorking`
}

func(t *InformationChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Student's chaincode is starting up")
	fmt.Println(" - ready for action");
	id, err := cid.GetID(stub);
	fmt.Printf("id obtained %v",id);
	mspid, err := cid.GetMSPID(stub);
	fmt.Printf("mspid obtained %v",mspid);
	val, ok, err := cid.GetAttributeValue(stub, "position");
	if err != nil {
		fmt.Println("There was an error trying to retrieve the attribute");
	}else if !ok {
		fmt.Printf("The client identity does not possess the attribute %v",val);
	}else{
		fmt.Printf("value obtained %v",val);
	}
	return shim.Success(nil)
}
func (t *InformationChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(" Invoke ");
	function,args := stub.GetFunctionAndParameters()
	if function == "init" {
		return t.Init(stub)	
	}else if function == "register" {
		return register(stub, args)
	}else if function =="addEducation"{
		return addEducation(stub,args)
	}else if function=="addExperience"{
		return AddExperience(stub,args)
	}else if function=="updateEducation"{
        return UpdateEducation(stub,args)
    }else if function =="updateExperience"{
       return UpdateExperience(stub,args)  
    }else if function =="getPersonalInfo"{
		return getPersonalInfo(stub,args)    
	}else if function =="getInfoById"{
		return getInfoById(stub,args)    
	}
    
    logger.Errorf("Received unknown invoke function name -%s",function)
	return shim.Error("Received unknown invoke function name -'" + function + "'")
}
func register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Registration of student in process... ")
	var education []Education
	var experience []Profession
	if len(args) !=7{
		return shim.Error("Incorrect number of arguments. Expecting 7. ID followed by details")
	}

	// input sanitation
	//err = sanitize_arguments(args)
	//if err != nil {
	//	return shim.Error(err.Error())
	//}
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
		return shim.Error("Failed to get ID  details: " + err.Error())
	} else if DetailsAsBytes != nil {
		fmt.Println("This  Id  already exists: "+args[0])
		fmt.Println("Registration Failed")
		return shim.Error("This  Id already exists: " + args[0])
	}
	//Assinging to student json
	//education := Education{degree,board,institute,passOut,score}
	registerDetails := User{displayPicture,name,email,dob,gender,profession,education,experience}
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
func  addEducation(stub shim.ChaincodeStubInterface,args []string)pb.Response{
	
	if len(args)<6{        
		return shim.Error("Incorrect number of argument passed for adding eduction")
	}

	id:=args[0]
	degree :=args[1]
	board:=args[2]
	institute:=args[3] 
	yop,err:=strconv.Atoi(args[4]) 
	score,err:=strconv.ParseFloat(args[5],64)
	if err != nil {
		return shim.Error("Cannot convert " + err.Error())
	}
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

	updatedetails:=User{informationdetails.ProfilePic,informationdetails.Name,informationdetails.Email,
	informationdetails.DateOfBirth,informationdetails.Gender,informationdetails.Profession,education,informationdetails.Experience}
			
	detailseducationmarshal,err:=json.Marshal(updatedetails)
	if err!=nil {
		fmt.Println("error occured while converting to json")
		return shim.Error(err.Error())
	}
	err=stub.PutState(id,detailseducationmarshal)
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
	shortDate :="2006-12-15"
	doj,err := time.Parse(shortDate, args[4])
	if err != nil {
		return shim.Error("4st argument must be a date string and should be in 2006-01-02 format")
	}

	dor,err := time.Parse(shortDate, args[5])
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
	fmt.Printf("got request for Personal ID  = %v",val1);
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
	buffer.WriteString(string(info))
	buffer.WriteString("}")
	bArrayMemberAlreadyWritten = true
		//fmt.Printf("value is %s",string(queryResponse.Value))
	
	buffer.WriteString("]")
   	fmt.Printf("- Personal Request:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func getInfoById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!= 1{
		shim.Error("Expected only one arg")
	}
	val,ok,err := cid.GetAttributeValue(stub,"accountType");
	if err != nil {
		shim.Error("There was an error trying to retrieve accountType attribute");
	}
	if !ok {
		shim.Error("The client identity does not possess accountType attribute");
	}
	if val =="student" || val =="employee"{
		shim.Error("This User cannot make this request");
	}
	queryFor := args[0];
	fmt.Printf("account Type : %v",val)
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
	buffer.WriteString(string(info))
	buffer.WriteString("}")
	bArrayMemberAlreadyWritten = true
		//fmt.Printf("value is %s",string(queryResponse.Value))
	
	buffer.WriteString("]")
		fmt.Printf("- Request By ID:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}