package cell

import (
	"fmt"
)

var EvalService = startEvalService()

type EvalError struct {
	Err string
}

func newEvalError(e string) EvalError {
	r := EvalError{
		Err: e,
	}
	return r
}

type EvalResult struct {
	Cell Cell
	Err  error
}

func newEvalResult(c Cell, e error) EvalResult {
	r := EvalResult{
		Cell: c,
		Err:  e,
	}
	return r
}

type EvalRequest struct {
	Cell      Cell
	Env       *EnvironmentEntry
	ReplyChan chan EvalResult
}

func NewEvalRequest(c Cell, env *EnvironmentEntry, replChan chan EvalResult) EvalRequest {
	r := EvalRequest{
		Cell:      c,
		Env:       env,
		ReplyChan: replChan,
	}
	return r
}

func (e EvalError) Error() string {
	return e.Err
}

func startEvalService() chan EvalRequest {
	service := make(chan EvalRequest)
	go server(service)
	return service
}

func server(service <-chan EvalRequest) {
	for {
		req := <-service
		go serve(req)
	}
}

func serve(req EvalRequest) {
	replyChan := req.ReplyChan
	replyChan <- eval(req)
}

func eval(req EvalRequest) EvalResult {
	toEval := req.Cell
	env := req.Env
	if toEval == nil {
		return newEvalResult(nil, nil)
	}
	switch c := toEval.(type) {
	case *IntCell:
		return newEvalResult(c, nil)
	case *StringCell:
		return newEvalResult(c, nil)
	case *SymbolCell:
		return newEvalResult(c, nil)
	case *ConsCell:
		switch car := c.Car.(type) {
		case *BuiltinMacroCell:
			return car.Macro(c.Cdr, env)
		default:
			argsResult := evlis(c.Cdr)
			if argsResult.Err != nil {
				return newEvalResult(nil, argsResult.Err)
			} else {
				return apply(car, argsResult.Cell, env)
			}
		}
	default:
		error := newEvalError("[eval] Unknown cell type: " + fmt.Sprintf("%v", toEval))
		return newEvalResult(nil, error)
	}
}

func evlis(args Cell) EvalResult {
	unvaluedArgs := extractArgs(args)

	if len(*unvaluedArgs) == 0 {
		return newEvalResult(nil, nil)
	}

	var replyChans []chan EvalResult
	n := len(*unvaluedArgs)
	for i := 0; i < n-1; i++ {
		newChan := make(chan EvalResult)
		replyChans = append(replyChans, newChan)
		go serve(NewEvalRequest((*unvaluedArgs)[i], EmptyEnv(), newChan)) // TODO: empty env!!
	}

	lastArgResult := eval(NewEvalRequest((*unvaluedArgs)[n-1], EmptyEnv(), nil))
	if lastArgResult.Err != nil {
		return lastArgResult
	}

	var top Cell
	var actCons Cell
	for i := 0; i < n; i++ {
		if i == n-1 {
			if top == nil {
				top = MakeCons(lastArgResult.Cell, nil)
				actCons = top
			} else {
				tmp := MakeCons(lastArgResult.Cell, nil)
				switch actConsCasted := actCons.(type) {
				case *ConsCell:
					actConsCasted.Cdr = tmp
					actCons = actConsCasted.Cdr
				}
			}
		} else {
			evaluedArg := <-replyChans[i]
			if evaluedArg.Err != nil {
				return newEvalResult(nil, evaluedArg.Err)
			}
			if top == nil {
				top = MakeCons(evaluedArg.Cell, nil)
				actCons = top
			} else {
				tmp := MakeCons(evaluedArg.Cell, nil)
				switch actConsCasted := actCons.(type) {
				case *ConsCell:
					actConsCasted.Cdr = tmp
					actCons = actConsCasted.Cdr
				}
			}
		}
	}

	return newEvalResult(top, nil)
}

func extractArgs(args Cell) *[]Cell {

	act := args
	var argsArray = new([]Cell)
	if args == nil {
		return argsArray
	}
	for act != nil {
		switch actCons := act.(type) {
		case *ConsCell:
			*argsArray = append(*argsArray, actCons.Car)
			act = actCons.Cdr
		default:
			panic("wrong argument format")
		}
	}
	return argsArray
}

func apply(function Cell, args Cell, env *EnvironmentEntry) EvalResult {
	switch functionCasted := function.(type) {
	case *BuiltinLambdaCell:
		return functionCasted.Lambda(args, env)
	default:
		return newEvalResult(nil, newEvalError("[apply] partial implementation"))
	}
}
