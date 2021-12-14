/*
 * The smart contract for v2g based on blockchain
 *
 */

package main

/* Imports
 * utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the states of the v2gcontract
const DEPLOY = 1     // the contract has been deployed
const SETTLED = 2       // the contract has been settled


// Define the V2Gcontract structure, with 6 properties.  Structure tags are used by encoding/json library
// tx of the V2Gcontract as the key, as the following are the value
type V2Gcontract struct {
	Tx   string `json:"Tx"`
	Usr    string `json:"Usr"`
	Caas string `json:"Caas"`
	StartTime string `json:"StartTime"`
	EndTime string `json:"EndTime"`
	Signals string `json:"Signals"`
	VcState string `json:"VcState"`
	Pay string `json:"Pay"`
}


/*
 * The Init method is called when the Smart Contract "V2Gcontract" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "V2Gcontract"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "deploy" {
		return s.deploy(APIstub, args)
	} else if function == "control" {
		return s.control(APIstub, args)
	} else if function == "settlement" {
		return s.settlement(APIstub, args)
	} else if function == "queryTx" {
		return s.queryTx(APIstub, args)
	}else if function == "queryUsr" {
		return s.queryUsr(APIstub, args)
	}else if function == "initLedger" {
		return s.initLedger(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/**
  根据 Tx 键查询合约内容
*/
func (s *SmartContract) queryTx(APIstub shim.ChaincodeStubInterface, arg []string) sc.Response {

	args := strings.Split(arg[0], ",") // split the args

	// args[0]: Tx
	// return v2gcontract string of the corresponding Tx

	if len(args) != 1 {
		return shim.Error("Incorrect format of arguments. Expecting 1.")
	}

	// get the v2gcontract from the world state
	vcAsBytes, _ := APIstub.GetState(args[0])

	return shim.Success(vcAsBytes) //vcAsBytes include the Tx
}

/**
  根据 Usr 键查询该用户拥有的充电合约
*/
func (s *SmartContract) queryUsr(APIstub shim.ChaincodeStubInterface, arg []string) sc.Response {

	args := strings.Split(arg[0], ",") // split the args

	// args[0]: Usr
	// return Contract string owned by the Usr
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"Usr\":\"%s\"}}", args[0]) // 查询构造

	queryResults, err := getQueryResultForQueryString(APIstub, queryString) // 富查询

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults) //vcAsBytes

}


/**
  账本初始化，不需要操作。仅仅为了完整性而保留
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	return shim.Success(nil)
}

/**
  部署V2G合约
*/
func (s *SmartContract) deploy(APIstub shim.ChaincodeStubInterface, arg []string) sc.Response {

	args := strings.Split(arg[0], ",") // split the args

	// args: Tx, Usr, Caas, startTime, endTime
	// create a v2gcontract and return the args


	// check the length of args
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// if Tx already used
	// fetch the vc according to the Tx
	vcAsBytes1, _ := APIstub.GetState(args[0])
	v2gcontract1 := V2Gcontract{}
	_ = json.Unmarshal(vcAsBytes1, &v2gcontract1)

	// check if the vs is in effect (deployed)
	if v2gcontract1.VcState == strconv.Itoa(DEPLOY) {
		return shim.Error("Tx already used")
	}

	// create a new v2gcontract
	var v2gcontract = V2Gcontract{Tx: args[0], Usr: args[1], Caas: args[2], StartTime: args[3], EndTime: args[4], Signals: ""}

	v2gcontract.VcState = strconv.Itoa(DEPLOY)// 设置状态

	// 生成充电合约,并更新状态
	vcAsBytes, _ := json.Marshal(v2gcontract)
	APIstub.PutState(args[0], vcAsBytes)

	return shim.Success(vcAsBytes)
}


/**
  控制函数：上传控制信号
*/
func (s *SmartContract) control(APIstub shim.ChaincodeStubInterface, arg []string) sc.Response {

	args := strings.Split(arg[0], ",") // split the args
	// args: Tx, Signals
	// Append control Signals to the Signal of the vc according to Tx

	if len(args) < 2 || len(args)%3 != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 4,7,...")
	}

	// fetch the vc according to the Tx
	vcAsBytes, _ := APIstub.GetState(args[0])
	v2gcontract := V2Gcontract{}
	_ = json.Unmarshal(vcAsBytes, &v2gcontract)

	// check if the vs is in effect (deployed)
	if v2gcontract.VcState != strconv.Itoa(DEPLOY) {
		return shim.Error("The v2gcontract is not in effect!")
	}

	// add signal to the Signals
	// TODO: the signal should be further checked
	for idx := 1; idx < len(args); idx ++ {
		if len(v2gcontract.Signals) != 0 {
			v2gcontract.Signals += "," // seperate with comma
		}
		v2gcontract.Signals += args[idx] // append
	}

	// marshal the vc and update the state according to Tx
	vcAsBytes, _ = json.Marshal(v2gcontract)
	APIstub.PutState(args[0], vcAsBytes)

	return shim.Success(vcAsBytes)
}

/**
  结算
*/
func (s *SmartContract) settlement(APIstub shim.ChaincodeStubInterface, arg []string) sc.Response {

	args := strings.Split(arg[0], ",") // split the args
	// args: Tx
	// settle the vc according to Tx

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// fetch the vc according to the Tx
	vcAsBytes, _ := APIstub.GetState(args[0])
	v2gcontract := V2Gcontract{}
	_ = json.Unmarshal(vcAsBytes, &v2gcontract)

	// check the state of the vc
	if v2gcontract.VcState != strconv.Itoa(DEPLOY) {
		return shim.Error("The v2gcontract has already been settled!")
	}

	// get the Signals of the vc
	Signals := strings.Split(v2gcontract.Signals, ",")
	// calculate the price
	Pay := float32(0.0)
	// TODO 加上上传完整性检验，一旦失败，需要惩罚
	for idx := 0; 3 * idx < len(Signals); idx ++ {
		p, _ := strconv.ParseFloat(Signals[3 * idx + 1], 32)
		soc, _ := strconv.ParseFloat(Signals[3 * idx + 2], 32)

		Pay += deg(float32(p), float32(soc))
	}

	// update the state of the vc and on-block it
	v2gcontract.Pay = fmt.Sprintf("%f", Pay) 
	v2gcontract.VcState = strconv.Itoa(SETTLED)

	vcAsBytes, _ = json.Marshal(v2gcontract)
	APIstub.PutState(args[0], vcAsBytes)

	return shim.Success(vcAsBytes)
}

// The degradation funtion
func deg(p float32, soc float32) float32 {
	// args: power p, soc
	// return Pay in the time slot
	const coef = 4.0
	const soc_opt = 0.5
	const base = 1.0
	return p * (coef * (soc - soc_opt) * (soc - soc_opt) + base)
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

// Two functions to support rich query of listed v2gcontracts
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		//buffer.WriteString("\"Deployed V2Gcontracts\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(queryResponse.Key)
		//buffer.WriteString("\"")

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("\n")

	}

	return &buffer, nil
}
