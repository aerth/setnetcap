// The MIT License (MIT)
//
// Copyright (c) 2016-2017 aerth <aerth@riseup.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Source: github.com/aerth/setnetcap

// setnetcap only setcaps (bind capability) and should be setuid root
package main

import (
	"fmt"
	"bytes"
	"os"
	"os/exec"
)

const setcaploc = "/sbin/setcap"
const setcapline = "cap_net_bind_service=+ep"

func init(){
	if len(os.Args) < 2 {
		println("Need an executable as argument 1")
		os.Exit(111)
	}
}


func main() {
	// security warning: you are giving user setcap net_bind cap on their choice of binary
	target := os.Args[1]

	// try using /usr/bin/logger to write to /var/log/messages (not caring if it doesn't exist)
	tattle(target)

	// make sure we are setuid and owned by root
	err := selfcheck()
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}

	// make sure target is not a symlink, directory, etc
	err = targetcheck(target)
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}

	// passed. run setcap
	err = run(exec.Command(setcaploc, setcapline, target))
	if err != nil {
		print(err.Error())
		os.Exit(111)
	}

	// no output = no errors
	os.Exit(0)
}


func run(cmd *exec.Cmd) error {
	var buf = new(bytes.Buffer)
	cmd.Stdout = os.Stdout
	cmd.Stderr = buf
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		// first line of stderr is all we need
		lin, _ := buf.ReadString(byte('\n'))
		if lin != "" {
			return fmt.Errorf(lin)
		}
		return err
	}

	return nil
}

func selfcheck() error {

		// are we setuid?
		this, _ := os.Executable()
		stat, err := os.Stat(this)
	 	if err != nil {
			return err
	 	}
		switch mode := stat.Mode(); {
			case mode&os.ModeSetuid != 0: // good
			case mode.IsDir():
				return fmt.Errorf("setnetcap needs to be setuid, got directory")
			case mode&os.ModeSymlink != 0:
				return fmt.Errorf("setnetcap needs to be setuid, got symlink")
			case mode.IsRegular():
				return fmt.Errorf("setnetcap needs to be setuid, got regular file")
			}

		// are we owned by root?
		if os.Geteuid() != 0 {
			return fmt.Errorf("%s is setuid, but not owned by root. try again", os.Args[0])
		}

		return nil
}


func targetcheck(t string) error {
		stat, err := os.Stat(t)
	 	if err != nil {
			return err
	 	}
		switch mode := stat.Mode(); {
			case mode.IsDir():
				return fmt.Errorf("Target needs to be file, got directory")
			case mode&os.ModeSymlink != 0:
				return fmt.Errorf("Target needs to be file, got symlink")
			case mode.IsRegular():
				return nil
			default:
				return fmt.Errorf("%s", mode.String())
			}
}

func tattle(target string){
	cmd := exec.Command("logger", "setnetcap is setuid and running setcap on "+target)
	cmd.Run()
}
