package cell

import (
	"fmt"
)

var evalService = startevalService()

func newEvalError(e string) EvalError {
	r := EvalError{
		Err: e,
	}
	return r
}

func newEvalResult(c Cell, e error) EvalResult {
	r := EvalResult{
		Cell: c,
		Err:  e,
	}
	return r
}

func newEvalPositiveResult(c Cell) EvalResult {
	return newEvalResult(c, nil)
}

func newEvalErrorResult(e error) EvalResult {
	return newEvalResult(nil, e)
}

type evalRequest struct {
	Cell      Cell
	Env       *environmentEntry
	ReplyChan chan EvalResult
}

func newEvalRequest(c Cell, env *environmentEntry, replChan chan EvalResult) evalRequest {
	r := evalRequest{
		Cell:      c,
		Env:       env,
		ReplyChan: replChan,
	}
	return r
}

func startevalService() chan evalRequest {
	service := make(chan evalRequest)
	go server(service)
	return service
}

func server(service <-chan evalRequest) {
	for {
		req := <-service
		go serve(req)
	}
}

func serve(req evalRequest) {
	result := eval(req.Cell, req.Env)
	req.ReplyChan <- result
}

func eval(toEval Cell, env *environmentEntry) EvalResult {
	if toEval == nil {
		return newEvalPositiveResult(nil)
	}
	switch c := toEval.(type) {
	case *IntCell:
		return newEvalPositiveResult(c)
	case *StringCell:
		return newEvalPositiveResult(c)
	case *SymbolCell:
		return assoc(c, env)
	case *ConsCell:
		switch car := c.Car.(type) {
		case *BuiltinMacroCell:
			return car.Macro(c.Cdr, env)
		default:
			argsResult := c.Evlis(c.Cdr, env)
			if argsResult.Err != nil {
				return newEvalErrorResult(argsResult.Err)
			} else {
				return apply(car, argsResult.Cell, env)
			}
		}
	// builtin symbols autoquote: allows higer order functions
	case *BuiltinMacroCell:
		return newEvalPositiveResult(c)
	case *BuiltinLambdaCell:
		return newEvalPositiveResult(c)
	default:
		return newEvalErrorResult(newEvalError("[eval] Unknown cell type: " + fmt.Sprintf("%v", toEval)))
	}
}

func evlisParallel(args Cell, env *environmentEntry) EvalResult {
	n := getNumberOfArgs(args)

	if n == 0 {
		return newEvalPositiveResult(nil)
	}

	var replyChans []chan EvalResult
	act := args
	for act != nil && cdr(act) != nil {
		newChan := make(chan EvalResult)
		replyChans = append(replyChans, newChan)
		go serve(newEvalRequest(car(act), env, newChan))
		act = cdr(act)
	}

	lastArgResult := eval(car(act), env)
	if lastArgResult.Err != nil {
		return lastArgResult
	}

	var top Cell
	var actCons Cell
	for i := 0; i < n; i++ {
		if i == n-1 {
			appendCellToArgs(&top, &actCons, &(lastArgResult.Cell))
		} else {
			evaluedArg := <-replyChans[i]
			if evaluedArg.Err != nil {
				return newEvalErrorResult(evaluedArg.Err)
			}
			appendCellToArgs(&top, &actCons, &(evaluedArg.Cell))
		}
	}

	return newEvalPositiveResult(top)
}

func getNumberOfArgs(c Cell) int {
	count := 0
	act := c
	actNotNil := (act != nil)
	for actNotNil {
		count++
		if cdr(act) == nil {
			actNotNil = false
		}
		act = cdr(act)
	}
	return count
}

func evlisSequential(args Cell, env *environmentEntry) EvalResult {
	actArg := args
	var top Cell
	var actCons Cell
	var evaluedArg EvalResult

	for actArg != nil {
		evaluedArg = eval(actArg.(*ConsCell).Car, env)
		if evaluedArg.Err != nil {
			return evaluedArg
		}
		appendCellToArgs(&top, &actCons, &(evaluedArg.Cell))
		actArg = (actArg.(*ConsCell)).Cdr
	}
	return newEvalResult(top, nil)
}

func extractCars(args Cell) []Cell {
	act := args
	var argsArray []Cell
	if args == nil {
		return argsArray
	}
	var actCons *ConsCell
	for act != nil {
		actCons = act.(*ConsCell)
		argsArray = append(argsArray, actCons.Car)
		act = actCons.Cdr
	}
	return argsArray
}

// appends to append after actCell, maybe initializing top. Has side effects
func appendCellToArgs(top, actCell, toAppend *Cell) {
	if *top == nil {
		*top = makeCons((*toAppend), nil)
		*actCell = *top
	} else {
		tmp := makeCons((*toAppend), nil)
		actConsCasted := (*actCell).(*ConsCell)
		actConsCasted.Cdr = tmp
		*actCell = actConsCasted.Cdr
	}
}

func apply(function Cell, args Cell, env *environmentEntry) EvalResult {
	switch functionCasted := function.(type) {
	case *BuiltinLambdaCell:
		return functionCasted.Lambda(args, env)
	case *ConsCell:
		if lisp.isLambdaSymbol(functionCasted.Car) {
			formalParameters := cadr(function)
			lambdaBody := caddr(function)
			if isClosure(formalParameters, args) {
				return newEvalPositiveResult(buildClosure(lambdaBody, formalParameters, args))
			}
			newEnv, err := pairlis(formalParameters, args, env)
			if err != nil {
				return newEvalErrorResult(err)
			}
			return eval(lambdaBody, newEnv)
		} else {
			// partial apply
			partiallyAppliedFunction := eval(function, env)
			if partiallyAppliedFunction.Err != nil {
				return partiallyAppliedFunction
			}
			return apply(partiallyAppliedFunction.Cell, args, env)
		}
	case *SymbolCell:
		evaluedFunction := eval(function, env)
		if evaluedFunction.Err != nil {
			return newEvalErrorResult(evaluedFunction.Err)
		}
		return apply(evaluedFunction.Cell, args, env)
	default:
		return newEvalErrorResult(newEvalError("[apply] trying to apply non-builtin, non-lambda, non-symbol"))
	}
}

func assoc(symbol *SymbolCell, env *environmentEntry) EvalResult {
	if res, isInglobalEnv := globalEnv[symbol.Sym]; isInglobalEnv {
		return newEvalPositiveResult(res)
	}
	if env == nil {
		return newEvalErrorResult(newEvalError("[assoc] symbol " + symbol.Sym + " not in env"))
	}
	act := env
	for act != nil {
		if (act.Pair.Symbol.Sym) == symbol.Sym {
			return newEvalPositiveResult(act.Pair.Value)
		}
		act = act.Next
	}
	return newEvalErrorResult(newEvalError("[assoc] symbol " + symbol.Sym + " not in env"))
}

func pairlis(formalParameters, actualParameters Cell, oldEnv *environmentEntry) (*environmentEntry, error) {
	actFormal := formalParameters
	actActual := actualParameters
	newEntry := oldEnv
	for actFormal != nil {
		if actActual == nil {
			return nil, newEvalError("[parilis] not enough actual parameters")
		}
		newEntry = newenvironmentEntry((car(actFormal)).(*SymbolCell), car(actActual), newEntry)
		actFormal = (actFormal.(*ConsCell)).Cdr
		actActual = (actActual.(*ConsCell)).Cdr
	}
	return newEntry, nil
}
