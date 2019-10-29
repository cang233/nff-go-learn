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

func TestUnixTime(t *testing.T){
	fmt.Println(time.Now().Unix())
	bs := IntToBytes(time.Now().Unix())
	for _,b:=range bs{
		fmt.Println(b)
	}
}
