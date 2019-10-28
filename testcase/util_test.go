package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	fmt.Println(time.Now().UTC().String())
	before := time.Now()
	time.Sleep(3 * time.Second)
	fmt.Println(time.Now().Sub(before).Seconds())
}
