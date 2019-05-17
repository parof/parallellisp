package cell

import "fmt"

func EmptyEnv() *EnvironmentEntry {
	return nil
}

func SimpleEnv() *EnvironmentEntry {
	list, _ := Parse("(1 2 3 4)")
	return NewEnvironmentEntry(MakeSymbol("l").(*SymbolCell), list, nil)
}

func NewEnvironmentEntry(sym *SymbolCell, value Cell, next *EnvironmentEntry) *EnvironmentEntry {
	newEntry := new(EnvironmentEntry)
	newEntry.Pair = new(EnvironmentPair)
	newEntry.Pair.Symbol = sym
	newEntry.Pair.Value = value
	newEntry.Next = next
	return newEntry
}

type EnvironmentEntry struct {
	Pair *EnvironmentPair
	Next *EnvironmentEntry
}

func (e *EnvironmentEntry) String() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%v -> %v\n", e.Pair.Symbol, e.Pair.Value) + fmt.Sprintf("%v", e.Next)
}

type EnvironmentPair struct {
	Symbol *SymbolCell
	Value  Cell
}

var GlobalEnv = make(map[string]Cell)

func initGlobalEnv() {
	GlobalEnv["id"], _ = Parse("(lambda (x) x)")
}
