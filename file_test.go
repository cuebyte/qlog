package qlog

import (
	"fmt"
	"log"
	"testing"
	//"time"
)

func TestFile(t *testing.T) {
	l, err := InitLogger(TRAC)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100000; i++ {
		l.Debug(fmt.Sprintf("test %d", i))
	}
	l.Close()
}
