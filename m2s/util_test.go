package m2s

import (
	"fmt"
	"testing"
)

type Animal struct {
	Name  string
	Order int64
	Yes   bool
}

func TestM2S(t *testing.T) {
	var targ Animal
	m := map[string]string{"Name": "Bill", "Order": "2", "Yes": "True"}
	Map2Struct(m, &targ)
	fmt.Printf("%v", targ)
}
