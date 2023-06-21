package main

import (
	"fmt"
	"os/exec"
)

type Ip struct {
    addr string
}

func main() {

    ip := string(out)
	fmt.Println(ip)
}

func dig() ([]byte, error) {
	args := []string{"+short", "myip.opendns.com", "@resolver1.opendns.com"}
	cmd := exec.Command("dig", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error getting DNS: ", err)
        return []byte{}, nil
	}

    return out, nil
}

func (ip *Ip) Get() string {
    return ip.addr
}

func (ip *Ip) Set(data []byte) {
    ip.addr = string(data)
}
