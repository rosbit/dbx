package dbx

import (
	"testing"
	"fmt"
)

func TestConds(t *testing.T) {
	printAndElem(Eq("name", "rosbit"), "Eq-str")
	printAndElem(Eq("age", 10), "Eq-int")
	printAndElem(And(Eq("name", "rosbit"), Eq("age", 10), Or(Eq("name", "r"), Eq("age", 1))), "And")
	printAndElem(Op("age", ">", 10), "Op")
	printAndElem(Or(Eq("name", "rosbit"), And(Eq("name", "rosbit"), Eq("age", 10)), Op("age", ">", 10)), "Or")
	printAndElem(In("age", 1,3,10), "IN-d")
	printAndElem(In("age", []int{1,3,10}), "IN-arr")
	printAndElem(NotIn("age", 1,3,10), "NotIn-d")
	printAndElem(Not(In("age", 1,3,10)), "Not-IN-d")
	printAndElem(And(Eq("name", "rosbit"), Or(Eq("name", "john"), Eq("age", 11)), Eq("age", 1)), "And-Or")
	printAndElem(Not(And(Eq("name", "rosbit"), Or(Eq("name", "john"), Eq("age", 11)), Eq("age", 1))), "Not-And-Or")
	printAndElem(Not(Or(Eq("name", "rosbit"), And(Eq("name", "rosbit"), Eq("age", 10)), Op("age", ">", 10))), "Not-Or-And")
}

func printAndElem(e AndElem, prompt string) {
	q, v := e.mkAndElem()
	fmt.Printf("%s: %s, %#v\n", prompt, q, v)
}

func TestWhere(t *testing.T) {
	w := Where(Eq("name", "rosbit"), Or(Eq("name", "john"), Eq("age", 11)))
	fmt.Printf("%#v\n", w)
}
