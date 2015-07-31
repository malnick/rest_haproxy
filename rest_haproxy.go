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

type Backend struct {
	Name      string
	Endpoints []string
}

func parsefile(filename string) (backends, error) {
	// Backend Array Defined
	var backends []*Backend
	endpoints := make(map[string][]*Backend)

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
		for _, b := range backends {
			if match_bkend.MatchString(line) {
				log.Println("MATCHED BACKEND: ", line)
				larry := strings.Fields(line)
				b.Name = larry[1]
				for _, e := range b.Endpoints {
					if match_srv.MatchString(line) {
						log.Println("MATCHED SERVER:\n", line)
						larry := strings.Fields(line)
						log.Println("LENGTH: ", len(larry))
						dest := strings.Split(larry[2], ":")
						ip := dest[0]
						port := dest[1]
						if len(larry) == 6 {
							mgmt := larry[5]
						} else {
							mgmt := port
						}
						endpoint := ip + ":" + mgmt
						log.Println("IP: ", ip)
						log.Println("MGMT: ", mgmt)
						log.Println("PORT: ", port)
						log.Println("ENDPOINT ", endpoint)
						endpoints[e] = append(endpoints[e], b)
					}
				}
			}
		}
	}

	return &Services{
			Services: backends,
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
