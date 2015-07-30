package main

import (
	"bufio"
	//"encoding/json"
	//"io/ioutil"
	"log"
	//"net/http"
	"os"
	"regexp"
	"strings"
)

type Available struct {
	Name string
	Ip   string
	Port string
}

type Services struct {
	Servers []Available
}

//func (services *Services) AddServer(available Available) []Available {
//	services.Servers = append(services.Servers, available)
//	return services.Servers
//}

func parsefile(filename string) (*Services, error) {
	// Temp array
	var temp_avail []Available
	var a Available
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

			larry := strings.Split(line, " ")

			a.Name = larry[1]

			dest := strings.Split(larry[2], ":")

			a.Ip, a.Port = dest[0], dest[1]

			temp_avail = append(temp_avail, a)
		}
	}

	return &Services{
			Servers: temp_avail,
		},
		nil
}

func main() {
	avail, _ := parsefile("/etc/haproxy/haproxy.cfg")
	log.Println(avail)
}
