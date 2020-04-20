package vmtestjson

import (
	"fmt"
	"math/big"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

// TestToJSONString converts a test object to its JSON representation.
func TestToJSONString(testTopLevel []*Test) string {
	jobj := TestToOrderedJSON(testTopLevel)
	return oj.JSONString(jobj)
}

// TestToOrderedJSON converts a test object to an ordered JSON object.
func TestToOrderedJSON(testTopLevel []*Test) oj.OJsonObject {
	result := oj.NewMap()
	for _, test := range testTopLevel {
		result.Put(test.TestName, testToOJ(test))
	}

	return result
}

func testToOJ(test *Test) oj.OJsonObject {
	testOJ := oj.NewMap()

	if !test.CheckGas {
		ojFalse := oj.OJsonBool(false)
		testOJ.Put("checkGas", &ojFalse)
	}

	testOJ.Put("pre", accountsToOJ(test.Pre))

	var blockList []oj.OJsonObject
	for _, block := range test.Blocks {
		blockList = append(blockList, blockToOJ(block))
	}
	blocksOJ := oj.OJsonList(blockList)
	testOJ.Put("blocks", &blocksOJ)

	testOJ.Put("network", stringToOJ(test.Network))

	var blockhashesList []oj.OJsonObject
	for _, blh := range test.BlockHashes {
		blockhashesList = append(blockhashesList, byteArrayToOJ(blh))
	}
	blockHashesOJ := oj.OJsonList(blockhashesList)
	testOJ.Put("blockhashes", &blockHashesOJ)

	testOJ.Put("postState", accountsToOJ(test.PostState))
	return testOJ
}

func accountsToOJ(accounts []*Account) oj.OJsonObject {
	acctsOJ := oj.NewMap()
	for _, account := range accounts {
		acctOJ := oj.NewMap()
		acctOJ.Put("nonce", intToOJ(account.Nonce))
		acctOJ.Put("balance", intToOJ(account.Balance))
		storageOJ := oj.NewMap()
		for _, st := range account.Storage {
			storageOJ.Put(byteArrayToString(st.Key), byteArrayToOJ(st.Value))
		}
		acctOJ.Put("storage", storageOJ)
		acctOJ.Put("code", stringToOJ(account.OriginalCode))

		acctsOJ.Put(byteArrayToString(account.Address), acctOJ)
	}

	return acctsOJ
}

func blockToOJ(block *Block) oj.OJsonObject {
	blockOJ := oj.NewMap()

	var resultList []oj.OJsonObject
	for _, blr := range block.Results {
		resultList = append(resultList, resultToOJ(blr))
	}
	resultsOJ := oj.OJsonList(resultList)
	blockOJ.Put("results", &resultsOJ)

	var txList []oj.OJsonObject
	for _, tx := range block.Transactions {
		txList = append(txList, transactionToOJ(tx))
	}
	txsOJ := oj.OJsonList(txList)
	blockOJ.Put("transactions", &txsOJ)

	blockHeaderOJ := oj.NewMap()
	blockHeaderOJ.Put("gasLimit", intToOJ(block.BlockHeader.GasLimit))
	blockHeaderOJ.Put("number", intToOJ(block.BlockHeader.Number))
	blockHeaderOJ.Put("difficulty", intToOJ(block.BlockHeader.Difficulty))
	blockHeaderOJ.Put("timestamp", uint64ToOJ(block.BlockHeader.Timestamp))
	blockHeaderOJ.Put("coinbase", intToOJ(block.BlockHeader.Beneficiary))
	blockOJ.Put("blockHeader", blockHeaderOJ)

	return blockOJ
}

func transactionToOJ(tx *Transaction) oj.OJsonObject {
	transactionOJ := oj.NewMap()
	transactionOJ.Put("nonce", uint64ToOJ(tx.Nonce))
	transactionOJ.Put("function", stringToOJ(tx.Function))
	transactionOJ.Put("gasLimit", uint64ToOJ(tx.GasLimit))
	transactionOJ.Put("value", intToOJ(tx.Value))
	transactionOJ.Put("to", byteArrayToOJ(tx.To))

	var argList []oj.OJsonObject
	for _, arg := range tx.Arguments {
		argList = append(argList, byteArrayToOJ(arg))
	}
	argOJ := oj.OJsonList(argList)
	transactionOJ.Put("arguments", &argOJ)

	transactionOJ.Put("contractCode", byteArrayToOJ(tx.Code))
	transactionOJ.Put("gasPrice", uint64ToOJ(tx.GasPrice))
	transactionOJ.Put("from", byteArrayToOJ(tx.From))

	return transactionOJ
}

func resultToOJ(res *TransactionResult) oj.OJsonObject {
	resultOJ := oj.NewMap()

	var outList []oj.OJsonObject
	for _, out := range res.Out {
		outList = append(outList, byteArrayToOJ(out))
	}
	outOJ := oj.OJsonList(outList)
	resultOJ.Put("out", &outOJ)

	resultOJ.Put("status", intToOJ(res.Status))
	resultOJ.Put("message", stringToOJ(res.Message))
	resultOJ.Put("gas", uint64ToOJ(res.Gas))
	if res.IgnoreLogs {
		resultOJ.Put("logs", stringToOJ("*"))
	} else {
		if len(res.LogHash) > 0 {
			resultOJ.Put("logs", stringToOJ(res.LogHash))
		} else {
			resultOJ.Put("logs", logsToOJ(res.Logs))
		}
	}
	resultOJ.Put("refund", intToOJ(res.Refund))

	return resultOJ
}

// LogToString returns a json representation of a log entry, we use it for debugging
func LogToString(logEntry *LogEntry) string {
	logOJ := logToOJ(logEntry)
	return oj.JSONString(logOJ)
}

func logToOJ(logEntry *LogEntry) oj.OJsonObject {
	logOJ := oj.NewMap()
	logOJ.Put("address", byteArrayToOJ(logEntry.Address))
	logOJ.Put("identifier", byteArrayToOJ(logEntry.Identifier))

	var topicsList []oj.OJsonObject
	for _, topic := range logEntry.Topics {
		topicsList = append(topicsList, byteArrayToOJ(topic))
	}
	topicsOJ := oj.OJsonList(topicsList)
	logOJ.Put("topics", &topicsOJ)

	logOJ.Put("data", byteArrayToOJ(logEntry.Data))

	return logOJ
}

func logsToOJ(logEntries []*LogEntry) oj.OJsonObject {
	var logList []oj.OJsonObject
	for _, logEntry := range logEntries {
		logOJ := logToOJ(logEntry)
		logList = append(logList, logOJ)
	}
	logOJList := oj.OJsonList(logList)
	return &logOJList
}

func intToString(i *big.Int) string {
	if i == nil {
		return ""
	}
	if i.Sign() == 0 {
		return "0x00"
	}

	isNegative := i.Sign() == -1
	str := i.Text(16)
	if isNegative {
		str = str[1:] // drop the minus in front
	}
	if len(str)%2 != 0 {
		str = "0" + str
	}
	str = "0x" + str
	if isNegative {
		str = "-" + str
	}
	return str
}

func intToOJ(i JSONBigInt) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func byteArrayToString(byteArray JSONBytes) string {
	return byteArray.Original
}

func byteArrayToOJ(byteArray JSONBytes) oj.OJsonObject {
	return &oj.OJsonString{Value: byteArrayToString(byteArray)}
}

func uint64ToOJ(i uint64) oj.OJsonObject {
	return stringToOJ(fmt.Sprintf("%d", i))
}

func stringToOJ(str string) oj.OJsonObject {
	return &oj.OJsonString{Value: str}
}