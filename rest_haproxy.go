package main

import (
	"bufio"
	//"encoding/json"
	//"io/ioutil"
	"log"
	//"net/http"
	"os"
	"regexp"
)

type Available struct {
	Name string
	Ip   string
	Port string
}

type Services struct {
	Servers []string //[]Availablei
}

//func (services *Services) AddServer(available Available) []Available {
//	services.Servers = append(services.Servers, available)
//	return services.Servers
//}

func parsefile(filename string) (*Services, error) {
	// Temp array
	var temp []string
	// Define our regex to parse
	regex, err := regexp.Compile(`^\s*server`)
	if err != nil {
		return nil, err // there was a problem with the regular expression.
	}

	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		if regex.MatchString(line) {
			log.Println("Matched: %s\n", line)
			temp = append(temp, line)
		}
	}

	return &Services{
			Servers: temp,
		},
		nil
}

func main() {
	parsefile("/etc/haproxy/haproxy.cfg")
	//log.Println(avail)
}
