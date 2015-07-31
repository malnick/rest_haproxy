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
	Ip       string
	Port     string
	MgmtPort string
}

type Available struct {
	Name string
	Info []Processes
}

type Services struct {
	Services []Available
}

func parsefile(filename string) (*Services, error) {
	var a_arry []Available
	var p_arry []Processes
	a := new(Available)
	proc := new(Processes)

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

	//	Loop:
	for scanner.Scan() {
		line := scanner.Text()
		if match_bkend.MatchString(line) {
			log.Println("MATCHED BACKEND: ", line)
			larry := strings.Fields(line)
			// Define a new backend
			backend := larry[1]
			a.Name = backend
			continue
		}
		if match_srv.MatchString(line) {
			log.Println("MATCHED SERVER:\n", line)
			larry := strings.Fields(line)
			log.Println("LENGTH: ", len(larry))
			dest := strings.Split(larry[2], ":")
			proc.Ip = larry[2]
			port := dest[1]
			log.Println("IP: ", dest[0])
			log.Println("PORT: ", dest[1])
			if len(larry) == 6 {
				proc.MgmtPort = larry[5]
			} else {
				proc.MgmtPort = proc.Port
			}
			log.Println("MGMT PORT: ", proc.MgmtPort)

			p_arry = append(p_arry, proc)
			a.Info = append(a.Info, p_arry)
			continue
		}
	}

	return &Services{
			Services: a_arry,
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
