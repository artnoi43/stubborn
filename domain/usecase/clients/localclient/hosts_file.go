package localclient

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strings"
)

// localNetworkTable is this package's global variable
// used to query resolve local network queries.
// One IP addr may have different hostnames.
type localDnsTable map[string][]string

var localNetworkTable localDnsTable

func InitTable(tableFile string) {
	// Using concurrency to read and parse into
	// the same map will slow this down 2-3x
	localNetworkTable = make(localDnsTable)
	initJsonTable(tableFile)
	initFromEtcHostsFile()
}

// initJsonTable parses stubborn's JSON hosts file.
// File location is specified in the config file, or via the command-line.
func initJsonTable(tableFile string) {
	fp, err := os.Open(tableFile)
	if err != nil {
		log.Panicf("failed to open hosts file %s: %s\n", tableFile, err.Error())
	}
	defer fp.Close()
	if err := json.NewDecoder(fp).Decode(&localNetworkTable); err != nil {
		log.Panicf("failed to read hosts file %s: %s\n", tableFile, err.Error())
	}
}

// initFromEtcHostsFile parses UNIX /etc/hosts file into the map.
func initFromEtcHostsFile() {
	if runtime.GOOS == "windows" {
		return
	}
	fp, err := os.Open("/etc/hosts")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		if text[0] != '#' {
			fields := strings.Fields(text)
			if len(fields) >= 2 {
				ip, hostnames := fields[0], fields[1:]
				// skip IPv6
				if strings.Contains(ip, ":::") {
					continue
				}
				for _, hostname := range hostnames {
					localNetworkTable[ip] = append(localNetworkTable[ip], hostname)
				}
			}
		}
	}
}
