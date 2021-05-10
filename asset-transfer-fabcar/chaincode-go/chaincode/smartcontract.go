package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a Car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a simple car
type Car struct {
	ID     string `json:"ID"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
	Amount string `json:"amount"`
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		Car{ID: "CAR0", Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko", Amount: "1000"},
		Car{ID: "CAR1", Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad", Amount: "1000"},
		Car{ID: "CAR2", Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo", Amount: "1000"},
		Car{ID: "CAR3", Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max", Amount: "1000"},
		Car{ID: "CAR4", Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana", Amount: "1000"},
	}

	for _, car := range cars {
		carJSON, err := json.Marshal(car)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(car.ID, carJSON)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %v", err)
		}
	}

	return nil

}

// QueryCar returns the car stored in the world state with given id.
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, id string) (*Car, error) {
	carJSON, err := ctx.GetStub().GetState(id)
	fmt.Println("carJSON => ", carJSON)
	fmt.Println("Type => ", reflect.TypeOf(carJSON))
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if carJSON == nil {
		return nil, fmt.Errorf("The car %s does not exist", id)
	}

	var car Car
	err = json.Unmarshal(carJSON, &car)
	if err != nil {
		return nil, err
	}
	fmt.Println("car => ", car)
	fmt.Println("&car => ", &car)
	fmt.Println("car Type => ", reflect.TypeOf(car))
	return &car, nil
}

// In CouchDB,QueryCar returns the car stored in the world state with given id.
func (s *SmartContract) QueryCarCouchDB(ctx contractapi.TransactionContextInterface, query string) ([]*Car, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if resultsIterator == nil {
		return nil, fmt.Errorf("The car does not exist")
	}

	defer resultsIterator.Close()

	var cars []*Car
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var car Car
		err = json.Unmarshal(queryResponse.Value, &car)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}
	return cars, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]*Car, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all cars in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var cars []*Car
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var car Car
		err = json.Unmarshal(queryResponse.Value, &car)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}

	return cars, nil
}

// AddCar issues a new car to the world state with given details.
func (s *SmartContract) AddCar(ctx contractapi.TransactionContextInterface, id string, make string, model string, colour string, owner string) error {
	exists, err := s.CarExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("The car %s already exists", id)
	}

	car := Car{
		ID:     id,
		Make:   make,
		Model:  model,
		Colour: colour,
		Owner:  owner,
	}
	log.Println("car Add : ", car)
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}
	log.Println("carJSON : ", carJSON)
	err = ctx.GetStub().PutState(id, carJSON)
	if err != nil {
		return fmt.Errorf("Failed to put to world state. %v", err)
	}
	return nil
}

// ChangeOwner updates the owner field of car with given id in world state.
func (s *SmartContract) ChangeOwner(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	car, err := s.QueryCar(ctx, id)
	if err != nil {
		return err
	}

	car.Owner = newOwner
	log.Println("ChangeOwner carJSON : ", car)
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, carJSON)
}

// CarExists returns true when car with given ID exists in world state
func (s *SmartContract) CarExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	carJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state: %v", err)
	}

	return carJSON != nil, nil
}

// queryHistoryCars
func (s *SmartContract) QueryHistoryCars(ctx contractapi.TransactionContextInterface, id string) ([]*Car, error) {

	historyIer, error := ctx.GetStub().GetHistoryForKey(id)

	if error != nil {
		return nil, error
	}

	var cars []*Car
	for historyIer.HasNext() {
		queryResponse, err := historyIer.Next()
		var car Car
		if err != nil {
			return nil, err
		}
		if queryResponse.IsDelete {
			continue
		} else {
			err = json.Unmarshal(queryResponse.Value, &car)
			if err != nil {
				return nil, err
			}
		}

		cars = append(cars, &car)
	}

	return cars, nil
}

// deleteCar
func (s *SmartContract) DeleteCar(ctx contractapi.TransactionContextInterface, id string) error {

	_, err := s.QueryCar(ctx, id)
	if err != nil {
		return err
	}

	return ctx.GetStub().DelState(id)
}

// BuyCar decrease amount
func (s *SmartContract) BuyCar(ctx contractapi.TransactionContextInterface, id string, values string, writecolumn string, writevalue string) error {
	log.Println("values : ", values)
	log.Println("WriteColumn : ", writecolumn)
	log.Println("writevalue : ", writevalue)
	var bytes []byte
	bytes = []byte(values)

	var car Car
	var objmap map[string]interface{}
	err1 := json.Unmarshal(bytes, &objmap)
	if err1 != nil {
		fmt.Println("error")
	}

	objmap[writecolumn] = writevalue
	ID := objmap["ID"].(string)
	amount := objmap[writecolumn].(string)
	make := objmap["make"].(string)
	model := objmap["model"].(string)
	owner := objmap["owner"].(string)
	colour := objmap["colour"].(string)
	car = Car{
		ID:     ID,
		Make:   make,
		Model:  model,
		Colour: colour,
		Owner:  owner,
		Amount: amount,
	}
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	log.Println("carJSON : ", carJSON)
	return ctx.GetStub().PutState(id, carJSON)
}
