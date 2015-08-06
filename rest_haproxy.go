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
	Service map[string][]string
}

type Store struct {
	Name string
}

// Store the sub backend here in the sub for loop of the file parser
func (s Store) storedBackend() string {
	return s.Name
}

func getBackend(line string) (backend string, err error) {

	// Regex for Backend
	match_bkend, err := regexp.Compile(`^\s*backend`)
	if err != nil {
		return "Failed matching backend", err
	}
	if match_bkend.MatchString(line) {
		larry := strings.Fields(line)
		name := larry[1]
		backend = name
		log.Println("MATCHED BACKEND: ", backend)
		return backend, nil
	}

	return "null", nil
}

func getIp(line string) (ip string, err error) {
	// Regex for Server
	match_srv, err := regexp.Compile(`^\s*server`)
	if err != nil {
		return "Failed.", err // there was a problem with the regular expression.
	}

	if match_srv.MatchString(line) {
		larry := strings.Fields(line)
		svc_address := larry[2]
		// Get mgmt address
		if len(larry) > 5 {
			mgmt_port := larry[len(larry)-1]
			dest := strings.Split(svc_address, ":")
			ip := dest[0]
			log.Println("IP: ", ip)
			log.Println("MGMT: ", mgmt_port)
			mgmt_temp := []string{ip, ":", mgmt_port}
			mgmt_address := strings.Join(mgmt_temp, "")
			final_temp := []string{svc_address, mgmt_address}
			final := strings.Join(final_temp, " ")
			log.Println("MATCHED SERVER: ", final)
			return final, nil
		}
		log.Println("MATCHED SERVER: ", svc_address)
		return svc_address, nil
	}
	return "null", nil
}

func parsefile(filename string) (s Services, err error) {
	// You know, a simple []map[string]map[string][]string
	s.Service = make(map[string][]string)
	// Stupid hack that could probably be solved with recursion...
	var store Store

	// Handle the file
	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	// Two nested loops to deal with scoping of next backend being the last arg in the nested loop
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("STORED BACKEND: ", store.Name)
		if len(store.Name) == 0 {
			backend_name, _ := getBackend(line)
			if backend_name != "null" {
				s.Service[backend_name] = []string{}
				for scanner.Scan() {
					line := scanner.Text()
					check, _ := getBackend(line)
					if check != "null" {
						store.Name = check
						break
					}
					server_ip, _ := getIp(line)
					if server_ip != "null" {
						s.Service[backend_name] = append(s.Service[backend_name], server_ip)
					}
				}
			}
		} else {
			backend_name := store.Name
			if backend_name != "null" {
				s.Service[backend_name] = []string{}
				for scanner.Scan() {
					line := scanner.Text()
					check, _ := getBackend(line)
					if check != "null" {
						store.Name = check
						break
					}
					server_ip, _ := getIp(line)
					if server_ip != "null" {
						s.Service[backend_name] = append(s.Service[backend_name], server_ip)
					}
				}
			}
		}
	}

	log.Println("Final Hash:\n")
	log.Println(s)
	return s, nil
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
