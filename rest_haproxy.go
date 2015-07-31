package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Processes struct {
	Ip string
	Port string
	MgmtPort string
}

type Available struct {
	Name []Processes 
}

type Services struct {
	Services []Available
}

func parsefile(filename string) (*Services, error) {
	// Temp array
	var temp_avail []Available
	var a Available

	// Define our regex to parse
	match_bkend, err := regexp.Compile(`^\s*server`)
	match_srv, err := regexp.Compile(`^\s*backend`)

	if err != nil {
		return nil, err // there was a problem with the regular expression.
	}

	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if match_bkend.MatchString(line) {
			var n Available
			log.Println("MATCHED BACKEND: ", line)
			n.Name := strings.Fields(line)[1]
			continue
		}
		if match_srv.MatchString(line) {
			log.Println("MATCHED SERVER:\n", line)
			larry := strings.Fields(line)
			log.Println("LENGTH: ", len(larry))
			dest := strings.Split(larry[2], ":")
			a.Ip, a.Port = dest[0], dest[1]
			log.Println("IP: ", dest[0])
			log.Println("PORT: ", dest[1])
			if len(larry) == 6 {
				a.MgmtPort = larry[5]
			} else {
				a.MgmtPort = a.Port
			}
			log.Println("MGMT PORT: ", a.MgmtPort)

			temp_avail = append(temp_avail, a)
		}
	}

	return &Services{
			Services: temp_avail,
		},
		nil
}

func response(rw http.ResponseWriter, request *http.Request) {
	services, err := parsefile("/etc/haproxy/haproxy.cfg")
	if err != nil {
		log.Println("ERROR: ", err)
	}
	json, err := json.Marshal(services)
	rw.Write([]byte(json))
}

func main() {
	http.HandleFunc("/services", response)
	http.ListenAndServe(":3000", nil)
}
