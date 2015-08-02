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

func (s Store) storedBackend() string {
	return s.Name
}

func getBackend(line string) (backend string, err error) {

	// Regex for Backend
	match_bkend, err := regexp.Compile(`^\s*backend`)
	if err != nil {
		return "Failed", err
	}
	log.Println("ATTEMPTED MATCH ON LINE FOR BACKEND: ", line)
	if match_bkend.MatchString(line) {
		log.Println("MATCHED BACKEND: ", line)
		larry := strings.Fields(line)
		name := larry[1]
		backend = name
		log.Println(backend)
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
		log.Println("MATCHED SERVER:\n", line)
		larry := strings.Fields(line)
		ip := larry[2]
		log.Println(ip)
		return ip, nil
	}
	return "null", nil
}

func parsefile(filename string) (s Services, err error) {

	s.Service = make(map[string][]string)

	var store Store

	// Handle the file
	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	// Create backends array

	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		log.Println("STORED: ", store.Name)
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
		}
		server_ip, _ := getIp(line)
		if server_ip != "null" {
			s.Service[store.Name] = append(s.Service[store.Name], server_ip)
		}
	}

	// Start sub process to get servers
	//		for scanner.Scan() {
	//			line := scanner.Text()
	//			// Break sub process when new backend is found
	//			check_backend, _ := getBackend(line)
	///			log.Println("CHECK ", check_backend, backend_name)
	//			if check_backend != "null" {
	//				if check_backend != backend_name {
	//					log.Println("BREAKING SUB PROCESS")
	//					break
	//				}
	//				continue
	//			}
	///
	//		}
	//		}
	//	}

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
