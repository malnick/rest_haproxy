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

type Services struct {
	Available []map[string]map[string]string
}

func parsefile(filename string) (*Services, error) {
	backend := make(map[string]map[string]string)
	backends := []map[string]map[string]string{}

	// Regex for Server
	match_srv, err := regexp.Compile(`^\s*server`)
	if err != nil {
		return nil, err // there was a problem with the regular expression.
	}

	// Regex for Bmeackend
	match_bkend, err := regexp.Compile(`^\s*backend`)
	if err != nil {
		return nil, err
	}

	// Handle the file
	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	// For each line in the file...
	for scanner.Scan() {
		line := scanner.Text()
		if match_bkend.MatchString(line) {
			log.Println("MATCHED BACKEND: ", line)
			larry := strings.Fields(line)
			name := larry[1]
			backend[name] = make(map[string]string)
			log.Println(backends)
			for scanner.Scan() {
				line := scanner.Text()
				if match_bkend.MatchString(line) {
					break
				} else {
					if match_srv.MatchString(line) {
						log.Println("MATCHED SERVER:\n", line)
						larry := strings.Fields(line)
						log.Println("LENGTH: ", len(larry))
						dest := strings.Split(larry[2], ":")
						ip := dest[0]
						port := dest[1]
						mgmt := ""
						if len(larry) == 6 {
							mgmt := larry[5]
							log.Println("MGMT: ", mgmt)
						} else {
							mgmt := port
							log.Println("MGMT Set to SVC PORT: ", mgmt)
						}
						//backend[name].Ip = ip
						log.Println("IP: ", ip)
						log.Println("PORT: ", port)
						backend[name]["ip"] = ip
						backend[name]["port"] = port
						backend[name]["mgmt"] = mgmt
						backends = append(backends, backend)
					}
				}
			}
		}
	}

	return &Services{
			Available: backends,
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
