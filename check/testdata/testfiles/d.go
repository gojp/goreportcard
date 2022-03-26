package testfiles

import (
	"fmt"
	"time"
)

func foo() {
	a := time.Now()
	b := a.Sub(time.Now())
	fmt.Println(a, b)
}
