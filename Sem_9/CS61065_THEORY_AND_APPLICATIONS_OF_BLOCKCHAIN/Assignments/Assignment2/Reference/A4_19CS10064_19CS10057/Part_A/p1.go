//  Assignment No. - 4 Part - A
//  Hyperledger Fabric
//  CS61065 - Theory and Applications of Blockchain
//  Semester - 7 (Autumn 2022-23)
//  Group Members - Vanshita Garg (19CS10064) & Shristi Singh (19CS10057)

// Import dependencies and define smart contract
package main

import (
    "fmt"
    "log"
    "encoding/json"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}

// declare a structure to store data of a student in the ledger
type Student struct {
    RollNo string `json:"rollNo"`
    Name string `json:"name"`
}

// createStudent function to create a new student in the ledger
func (s *SmartContract) CreateStudent(ctx contractapi.TransactionContextInterface, roll string, name string) error {
    // it checks( using StudentExists() ) if this ‘roll’ is already present in the ledger or not.
    // If not present, it creates a new student with the given roll and name.
    // If present, it returns an error.
    exists, err := s.StudentExists(ctx, roll)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("the student %s already exists", roll)
    }
    // create a new student object
    student := Student{
        RollNo: roll,
        Name: name,
    }
    // convert the student object into JSON format
    studentJSON, err := json.Marshal(student)
    if err != nil {
        return err
    }
    // put the student object in the ledger
    return ctx.GetStub().PutState(roll, studentJSON)
}

// write a func StudentExists(ctx contractapi.TransactionContextInterface, roll string)(bool, error)
// to check if a student with the given roll number exists in the ledger or not.
func (s *SmartContract) StudentExists(ctx contractapi.TransactionContextInterface, roll string) (bool, error) {
    // get the student object from the ledger
    studentJSON, err := ctx.GetStub().GetState(roll)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }
    // if the student object is not present in the ledger, return false
    return studentJSON != nil, nil
}

// write a func ReadStudent(ctx contractapi.TransactionContextInterface, roll string)(*Student, error)
// to read the details of a student from the ledger.
func (s *SmartContract) ReadStudent(ctx contractapi.TransactionContextInterface, roll string) (string, error) {
    // return student's name corresponding to the given roll number and error if any
    exists, err := s.StudentExists(ctx, roll)
    if err != nil {
        return "", err
    }
    if !exists {
        return "", fmt.Errorf("the student %s does not exist", roll)
    }
    studentJSON, err := ctx.GetStub().GetState(roll)
    if err != nil {
        return "", fmt.Errorf("failed to read from world state: %v", err)
    }
    // return student's name corresponding to the given roll number and error if any
    return string(studentJSON), nil
}

// ReadAllStudents function to read all the students from the ledger
// It will return a list with each item having two values,
// ■ the Student’s name corresponding to ‘roll’
// ■ the error if any
// Then it uses ctx.GetStub().GetState() to access the Students’ data and
// returns the roll and name of all students.

func (s *SmartContract) ReadAllStudents(ctx contractapi.TransactionContextInterface) ([]string, error) {
    // get all the students from the ledger
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()
    // iterate through the results and add them to the list and  return the list with roll and name of all students
    var allStudents []string
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }
        var student Student
        err = json.Unmarshal(queryResponse.Value, &student)
        if err != nil {
            return nil, err
        }
        // append the student's name and roll number to the list in the form {Roll No: Name}
        allStudents = append(allStudents, fmt.Sprintf("{Roll No: %s, Name: %s}", student.RollNo, student.Name))
    }
    return allStudents, nil
}       

// function main
func main() {
    // create a new smart contract
    smartContract, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        log.Panicf("Error creating smart contract: %v", err)
    }
    // start the smart contract
    if err := smartContract.Start(); err != nil {
        log.Panicf("Error starting smart contract: %v", err)
    }
}