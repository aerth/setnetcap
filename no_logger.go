//+build nolog

package main

func init(){
  tattle = func(s string){}
}
