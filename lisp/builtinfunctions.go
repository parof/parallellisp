package lisp

import (
	"fmt"
	"io/ioutil"
	"time"
)

func condMacro(args Cell, env *environmentEntry) EvalResult {
	actBranch := args
	var condAndBody Cell
	var cond Cell
	var body Cell
	var condResult EvalResult
	for actBranch != nil {
		condAndBody = car(actBranch)
		cond = car(condAndBody)
		body = cadr(condAndBody)
		condResult = eval(cond, env)
		if condResult.Err != nil {
			return condResult
		} else if condResult.Cell != nil {
			return eval(body, env)
		}
		actBranch = cdr(actBranch)
	}
	return newEvalErrorResult(newEvalError("[cond] none condition was verified"))
}

func quoteMacro(args Cell, env *environmentEntry) EvalResult {
	switch cons := args.(type) {
	case *consCell:
		return newEvalPositiveResult(cons.Car)
	default:
		return newEvalErrorResult(newEvalError("[quote] Can't quote" + fmt.Sprint(cons)))
	}
}

func timeMacro(args Cell, env *environmentEntry) EvalResult {
	if args == nil {
		return newEvalErrorResult(newEvalError("[time] too few arguments"))
	}
	now := time.Now()
	start := now.UnixNano()

	result := eval(car(args), env)
	if result.Err != nil {
		return result
	}

	now = time.Now()
	afterEvalTime := now.UnixNano()
	elapsedMillis := (afterEvalTime - start) / 1000000
	fmt.Println("time:", elapsedMillis, "ms")

	return result
}

func lambdaMacro(args Cell, env *environmentEntry) EvalResult {
	// lambda autoquote
	return newEvalPositiveResult(makeCons(makeSymbol("lambda"), args))
}

func defunMacro(args Cell, env *environmentEntry) EvalResult {
	argsSlice := extractCars(args)
	if len(argsSlice) != 3 {
		return newEvalErrorResult(newEvalError("[defun] wrong number of arguments"))
	}
	name := argsSlice[0]
	formalParameters := argsSlice[1]
	lambdaBody := argsSlice[2]
	bodyCons := makeCons(lambdaBody, nil)
	argsAndBodyCons := makeCons(formalParameters, bodyCons)
	ret := makeCons(makeSymbol("lambda"), argsAndBodyCons)
	switch nameSymbolCell := name.(type) {
	case *symbolCell:
		globalEnv[nameSymbolCell.Sym] = ret
	default:
		return newEvalErrorResult(newEvalError("[defun] the name of the lambda must be a symbol"))
	}
	return newEvalPositiveResult(ret)
}

func setqMacro(args Cell, env *environmentEntry) EvalResult {
	argsSlice := extractCars(args)
	if len(argsSlice) != 2 {
		return newEvalErrorResult(newEvalError("[setq] wrong number of arguments"))
	}
	name := argsSlice[0]
	value := argsSlice[1]

	switch name.(type) {
	case (*symbolCell):
		evaluedVal := eval(value, env)
		if evaluedVal.Err != nil {
			return evaluedVal
		}
		newArgs := makeCons(name, makeCons(evaluedVal.Cell, nil))
		return setLambda(newArgs, env)
	default:
		return newEvalErrorResult(newEvalError("[setq] first argument must be a symbol"))
	}
}

func letMacro(args Cell, env *environmentEntry) EvalResult {
	pairs := car(args)
	newEnv := env
	for pairs != nil {
		evaluedValue := eval(cadar(pairs), env)
		if evaluedValue.Err != nil {
			return evaluedValue
		}
		newEnv = newEnvironmentEntry(caar(pairs).(*symbolCell), evaluedValue.Cell, newEnv)
		pairs = cdr(pairs)
	}
	return eval(cadr(args), newEnv)
}

func dotimesMacro(args Cell, env *environmentEntry) EvalResult {
	firstArg := car(args)
	body := cadr(args)
	varName := car(firstArg)
	varValue := cadr(firstArg)
	for i := 0; i < (varValue.(*intCell)).Val; i++ {
		newEnv := newEnvironmentEntry(varName.(*symbolCell), makeInt(i), env)
		eval(body, newEnv)
	}
	return newEvalPositiveResult(nil)
}

func carLambda(args Cell, env *environmentEntry) EvalResult {
	switch topCons := args.(type) {
	case *consCell:
		switch cons := topCons.Car.(type) {
		case *consCell:
			return newEvalResult(cons.Car, nil)
		default:
			return newEvalResult(nil, newEvalError("[car] car applied to atom"))
		}
	default:
		return newEvalResult(nil, newEvalError("[car] not enough arguments"))
	}
}

func cdrLambda(args Cell, env *environmentEntry) EvalResult {
	switch topCons := args.(type) {
	case *consCell:
		switch cons := topCons.Car.(type) {
		case *consCell:
			return newEvalResult(cons.Cdr, nil)
		default:
			return newEvalResult(nil, newEvalError("[cdr] cdr applied to atom"))
		}
	default:
		return newEvalResult(nil, newEvalError("[cdr] not enough arguments"))
	}
}

func consLambda(args Cell, env *environmentEntry) EvalResult {
	switch firstCons := args.(type) {
	case *consCell:
		switch cons := firstCons.Cdr.(type) {
		case *consCell:
			return newEvalPositiveResult(makeCons(firstCons.Car, cons.Car))
		default:
			return newEvalErrorResult(newEvalError("[cons] not enough arguments"))
		}
	default:
		return newEvalErrorResult(newEvalError("[cons] not enough arguments"))
	}
}

func eqLambda(args Cell, env *environmentEntry) EvalResult {
	switch firstArg := args.(type) {
	case *consCell:
		switch secondArg := firstArg.Cdr.(type) {
		case *consCell:
			if eq(firstArg.Car, secondArg.Car) {
				return newEvalPositiveResult(lisp.getTrueSymbol())
			}
			return newEvalPositiveResult(nil)
		default:
			return newEvalErrorResult(newEvalError("[eq] not enough arguments"))
		}
	default:
		return newEvalErrorResult(newEvalError("[eq] not enough arguments"))
	}
}

func atomLambda(args Cell, env *environmentEntry) EvalResult {
	switch firstCons := args.(type) {
	case *consCell:
		switch firstCons.Car.(type) {
		case *consCell:
			return newEvalPositiveResult(nil)
		default:
			return newEvalPositiveResult(lisp.getTrueSymbol())
		}
	default:
		return newEvalErrorResult(newEvalError("[atom] not enough arguments"))
	}
}

func plusLambda(args Cell, env *environmentEntry) EvalResult {
	tot := 0
	act := args
	for act != nil {
		tot += (car(act).(*intCell)).Val
		act = cdr(act)
	}
	return newEvalPositiveResult(makeInt(tot))
}

func multLambda(args Cell, env *environmentEntry) EvalResult {
	tot := 1
	act := args
	for act != nil {
		tot *= (car(act).(*intCell)).Val
		act = cdr(act)
	}
	return newEvalPositiveResult(makeInt(tot))
}

func minusLambda(args Cell, env *environmentEntry) EvalResult {
	if args == nil {
		return newEvalErrorResult(newEvalError("[-] too few arguments"))
	}
	tot := (car(args).(*intCell)).Val
	act := cdr(args)
	for act != nil {
		tot -= (car(act).(*intCell)).Val
		act = cdr(act)
	}
	return newEvalPositiveResult(makeInt(tot))
}

func orLambda(args Cell, env *environmentEntry) EvalResult {
	act := args
	for act != nil {
		if car(act) != nil {
			return newEvalPositiveResult(car(act))
		}
		act = cdr(act)
	}
	return newEvalPositiveResult(nil)
}

func andLambda(args Cell, env *environmentEntry) EvalResult {
	act := args
	var last Cell
	for act != nil {
		last = car(act)
		if last == nil {
			return newEvalPositiveResult(nil)
		}
		act = cdr(act)
	}
	return newEvalPositiveResult(last)
}

func notLambda(args Cell, env *environmentEntry) EvalResult {
	toNegate := car(args)
	if toNegate == nil {
		return newEvalPositiveResult(lisp.getTrueSymbol())
	}
	return newEvalPositiveResult(nil)
}

func greaterLambda(args Cell, env *environmentEntry) EvalResult {
	return listRelationalComparison(args, env, func(left, right int) bool { return left > right })
}

func greaterEqLambda(args Cell, env *environmentEntry) EvalResult {
	return listRelationalComparison(args, env, func(left, right int) bool { return left >= right })
}

func lessLambda(args Cell, env *environmentEntry) EvalResult {
	return listRelationalComparison(args, env, func(left, right int) bool { return left < right })
}

func lessEqLambda(args Cell, env *environmentEntry) EvalResult {
	return listRelationalComparison(args, env, func(left, right int) bool { return left <= right })
}

func listRelationalComparison(args Cell, env *environmentEntry, operator func(int, int) bool) EvalResult {
	act := cdr(args)
	last := car(args)
	for act != nil {
		if !(operator((last.(*intCell)).Val, (car(act).(*intCell)).Val)) {
			return newEvalPositiveResult(nil)
		}
		last = car(act)
		act = cdr(act)
	}
	return newEvalPositiveResult(lisp.getTrueSymbol())
}

func divLambda(args Cell, env *environmentEntry) EvalResult {
	if args == nil {
		return newEvalErrorResult(newEvalError("[/] too few arguments"))
	}
	tot := (car(args).(*intCell)).Val
	act := cdr(args)
	div := 0
	for act != nil {
		div = (car(act).(*intCell)).Val
		if div == 0 {
			return newEvalErrorResult(newEvalError("[/] division by zero"))
		}
		tot /= div
		act = cdr(act)
	}
	return newEvalPositiveResult(makeInt(tot))
}

func loadLambda(args Cell, env *environmentEntry) EvalResult {
	files := extractCars(args)
	if len(files) != 1 {
		return newEvalErrorResult(newEvalError("[load] load needs exaclty one argument"))
	}
	fileName := (files[0].(*stringCell)).Str
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return newEvalErrorResult(newEvalError("[load] error opening file " + fileName))
	}
	source := string(dat)
	sexpressions, err := parseMultipleSexpressions(source)
	if err != nil {
		return newEvalErrorResult(err)
	}
	var lastEvalued EvalResult
	for _, sexpression := range sexpressions {
		lastEvalued = eval(sexpression, env)
		if lastEvalued.Err != nil {
			return lastEvalued
		}
	}
	return lastEvalued
}

func writeLambda(args Cell, env *environmentEntry) EvalResult {
	phrases := extractCars(args)
	if len(phrases) == 0 {
		fmt.Println()
		return newEvalPositiveResult(makeString(""))
	} else if len(phrases) > 1 {
		return newEvalErrorResult(newEvalError("[write] write needs at least one string argument"))
	}
	fmt.Println((phrases[0].(*stringCell)).Str)
	return newEvalPositiveResult(phrases[0])
}

func listLambda(args Cell, env *environmentEntry) EvalResult {
	var top Cell
	var actLast Cell
	var newVal Cell
	act := args
	for act != nil {
		newVal = car(act)
		appendCellToArgs(&top, &actLast, &newVal)
		act = cdr(act)
	}
	return newEvalPositiveResult(top)
}

func reverseLambda(args Cell, env *environmentEntry) EvalResult {
	var top Cell
	act := car(args)
	for act != nil {
		top = makeCons(car(act), top)
		act = cdr(act)
	}
	return newEvalPositiveResult(top)
}

func memberLambda(args Cell, env *environmentEntry) EvalResult {
	toFind := car(args)
	act := cadr(args)
	for act != nil {
		if eq(toFind, car(act)) {
			return newEvalPositiveResult(act)
		}
		act = cdr(act)
	}
	return newEvalPositiveResult(nil)
}

func nthLambda(args Cell, env *environmentEntry) EvalResult {
	n := (car(args).(*intCell)).Val
	act := cadr(args)
	for n > 0 {
		n--
		act = cdr(act)
	}
	return newEvalPositiveResult(car(act))
}

func lengthLambda(args Cell, env *environmentEntry) EvalResult {
	return newEvalPositiveResult(makeInt(listLengt(car(args))))
}

func setLambda(args Cell, env *environmentEntry) EvalResult {
	id := car(args)
	val := cadr(args)
	globalEnv[(id.(*symbolCell)).Sym] = val
	return newEvalPositiveResult(val)
}

func onePlusLambda(args Cell, env *environmentEntry) EvalResult {
	num := car(args).(*intCell)
	return newEvalPositiveResult(makeInt(num.Val + 1))
}
func oneMinusLambda(args Cell, env *environmentEntry) EvalResult {
	num := car(args).(*intCell)
	return newEvalPositiveResult(makeInt(num.Val - 1))
}

func integerpLambda(args Cell, env *environmentEntry) EvalResult {
	switch car(args).(type) {
	case *intCell:
		return newEvalPositiveResult(lisp.getTrueSymbol())
	default:
		return newEvalPositiveResult(nil)
	}
}

func symbolpLambda(args Cell, env *environmentEntry) EvalResult {
	switch car(args).(type) {
	case *builtinLambdaCell:
		return newEvalPositiveResult(lisp.getTrueSymbol())
	case *builtinMacroCell:
		return newEvalPositiveResult(lisp.getTrueSymbol())
	case *symbolCell:
		return newEvalPositiveResult(lisp.getTrueSymbol())
	default:
		return newEvalPositiveResult(nil)
	}
}

func unimplementedMacro(c Cell, env *environmentEntry) EvalResult {
	panic("unimplemented macro")
}

func unimplementedLambda(c Cell, env *environmentEntry) EvalResult {
	panic("unimplemented lambda")
}
