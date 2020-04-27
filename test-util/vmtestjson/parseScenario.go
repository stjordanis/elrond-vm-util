package vmtestjson

import (
	"errors"
	"fmt"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

// ParseScenarioFile converts a scenario json string to scenario object representation
func (p *Parser) ParseScenarioFile(jsonString []byte) (*Scenario, error) {
	jobj, err := oj.ParseOrderedJSON(jsonString)
	if err != nil {
		return nil, err
	}

	topMap, isMap := jobj.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled test top level object is not a map")
	}

	scenario := &Scenario{
		CheckGas: true,
	}
	for _, kvp := range topMap.OrderedKV {
		switch kvp.Key {
		case "name":
			scenario.Name, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("bad scenario name: %w", err)
			}
		case "comment":
			scenario.Comment, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("bad scenario comment: %w", err)
			}
		case "checkGas":
			checkGasOJ, isBool := kvp.Value.(*oj.OJsonBool)
			if !isBool {
				return nil, errors.New("scenario checkGas flag is not boolean")
			}
			scenario.CheckGas = bool(*checkGasOJ)
		case "steps":
			scenario.Steps, err = p.processScenarioStepList(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error processing steps: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown step field: %s", kvp.Key)
		}
	}
	return scenario, nil
}

func (p *Parser) processScenarioStepList(obj interface{}) ([]Step, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("steps not a JSON list")
	}
	var stepList []Step
	for _, elemRaw := range listRaw.AsList() {
		step, err := p.processScenarioStep(elemRaw)
		if err != nil {
			return nil, err
		}
		stepList = append(stepList, step)
	}
	return stepList, nil
}

func (p *Parser) processScenarioStep(stepObj oj.OJsonObject) (Step, error) {
	stepMap, isStepMap := stepObj.(*oj.OJsonMap)
	if !isStepMap {
		return nil, errors.New("unmarshalled step object is not a map")
	}

	var err error
	step := ""
	for _, kvp := range stepMap.OrderedKV {
		if kvp.Key == "step" {
			step, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("step type not a string: %w", err)
			}
		}
	}

	switch step {
	case "":
		return nil, errors.New("no step type field provided")
	case stepNameSetState:
		step := &SetStateStep{}
		for _, kvp := range stepMap.OrderedKV {
			switch kvp.Key {
			case "step":
			case "comment":
				step.Comment, err = p.parseString(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("bad set state step comment: %w", err)
				}
			case "accounts":
				step.Accounts, err = p.processAccountMap(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("cannot parse set state step: %w", err)
				}
			case "newAddresses":
				step.NewAddressMocks, err = p.processNewAddressMocks(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("error parsing new addresses: %w", err)
				}
			case "blockhashes":
				step.BlockHashes, err = p.parseByteArrayList(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("error parsing block hashes: %w", err)
				}
			default:
				return nil, fmt.Errorf("invalid set state field: %s", kvp.Key)
			}
		}
		return step, nil
	case stepNameCheckState:
		step := &CheckStateStep{}
		for _, kvp := range stepMap.OrderedKV {
			switch kvp.Key {
			case "step":
			case "comment":
				step.Comment, err = p.parseString(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("bad check state step comment: %w", err)
				}
			case "accounts":
				step.CheckAccounts, err = p.processCheckAccountMap(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("cannot parse check state step: %w", err)
				}
			default:
				return nil, fmt.Errorf("invalid check state field: %s", kvp.Key)
			}
		}
		return step, nil
	case stepNameScCall:
		return p.parseTxStep(ScCall, stepMap)
	case stepNameScDeploy:
		return p.parseTxStep(ScDeploy, stepMap)
	case stepNameTransfer:
		return p.parseTxStep(Transfer, stepMap)
	default:
		return nil, fmt.Errorf("unknown step type: %s", step)
	}
}

func (p *Parser) parseTxStep(txType TransactionType, stepMap *oj.OJsonMap) (*TxStep, error) {
	step := &TxStep{}
	var err error
	for _, kvp := range stepMap.OrderedKV {
		switch kvp.Key {
		case "step":
		case "txId":
			step.TxIdent, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("bad tx step id: %w", err)
			}
		case "comment":
			step.Comment, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("bad tx step comment: %w", err)
			}
		case "tx":
			step.Tx, err = p.processTx(txType, kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("cannot parse tx step transaction: %w", err)
			}
		case "expect":
			step.ExpectedResult, err = p.processTxExpectedResult(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("cannot parse tx expected result: %w", err)
			}
		default:
			return nil, fmt.Errorf("invalid tx step field: %s", kvp.Key)
		}
	}
	return step, nil
}
