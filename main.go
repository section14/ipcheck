package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Site struct {
	Url  string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Creds struct {
	Sites     []Site `json:"sites"`
	CurrentIp string `json:"current_ip"`
}

func main() {
	creds, err := parseConfig("creds.json")
	if err != nil {
		fmt.Println("error opening creds: ", err)
	}

	for {
		//get LAN ip
		ip, err := dig()
		if err != nil {
			fmt.Println("error getting ip: ", err)
		}

		if ip != creds.CurrentIp {
			//do stuff
			err := updateIp(ip, creds.Sites)
			if err != nil {
				fmt.Println("error updating DNS: ", err)
			}
		}

		time.Sleep(2 * time.Hour)
	}

	//string for curl
	//URL="https://${USER}:${PASS}@domains.google.com/nic/update?hostname=${HOST}&myip=${IP}"
}

func updateIp(ip string, sites []Site) error {
	//dnsUrl := "https://%s:%s@domains.google.com/nic/update?hostname=%s&myip=%s"
    url := "http://ghost.sectionfourteen.com"

	for _, site := range sites {
		//url := fmt.Sprintf(dnsUrl, site.User, site.Pass, site.Url, ip)

		res, err := http.Get(url)
		if err != nil {
			fmt.Println("error updating site DNS: ", site.Url)
			//return errors.New("one or more DNS updates failed")
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		fmt.Println(string(body))
	}

	return nil
}

//parse JSON config file
func parseConfig(path string) (Creds, error) {
	var creds Creds

	//open local file
	file, err := os.Open(path)
	if err != nil {
		return creds, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&creds)
	if err != nil {
		return creds, err
	}

	return creds, nil
}

//run dig command
func dig() (string, error) {
	args := []string{"+short", "myip.opendns.com", "@resolver1.opendns.com"}
	cmd := exec.Command("dig", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error getting DNS: ", err)
		return "", nil
	}

	return strings.TrimSuffix(string(out), "\n"), nil
}
