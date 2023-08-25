package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/elmasy-com/elnet/dns"
	"github.com/elmasy-com/elnet/dns/hetzner"
	"golang.org/x/sys/unix"
)

func mustGetCommandArg() string {

	if len(os.Args) < 2 {

		fmt.Fprintf(os.Stderr, "Command is missing!\n")
		fmt.Printf("Use \"%s help\" for help!\n", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] != "set" && os.Args[1] != "unset" && os.Args[1] != "init" && os.Args[1] != "test" && os.Args[1] != "help" {
		fmt.Fprintf(os.Stderr, "Invalid command: %s!\n", os.Args[1])
		fmt.Printf("Use \"%s help\" for help!\n", os.Args[0])
		os.Exit(1)
	}

	return os.Args[1]
}

func mustGetDomainArg() string {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "domain is missing!\n")
		fmt.Printf("Use \"%s help\" for help!\n", os.Args[0])
		os.Exit(1)
	}

	return dns.GetDomain(os.Args[2])
}

func mustGetValidationNameArg() string {

	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "validation_name is missing!\n")
		fmt.Printf("Use \"%s help\" for help!\n", os.Args[0])
		os.Exit(1)
	}

	return dns.GetSub(os.Args[3])
}

func mustGetValidationContextArg() string {

	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "validation_context is missing!\n")
		fmt.Printf("Use \"%s help\" for help!\n", os.Args[0])
		os.Exit(1)
	}

	return os.Args[4]
}

func getToken() (string, error) {

	path := os.ExpandEnv("$HOME/.tahtoken")

	out, err := os.ReadFile(path)

	return strings.Trim(string(out), "\n"), err
}

func testTokenFile() error {

	path := os.ExpandEnv("$HOME/.tahtoken")

	statT := new(unix.Stat_t)

	err := unix.Stat(path, statT)
	if err != nil {
		return err
	}

	if statT.Uid != uint32(os.Geteuid()) {
		fmt.Printf("Different user for config file: %d/%d\n", statT.Uid, os.Geteuid())
	}

	if statT.Gid != uint32(os.Getegid()) {
		fmt.Printf("Different group for config file: %d/%d\n", statT.Uid, os.Geteuid())
	}

	if statT.Mode != 0o100600 {
		fmt.Printf("Dangerous file mode: %o, use 0600!\n", statT.Mode)
	}

	return nil
}

func Set() {

	domain := mustGetDomainArg()
	validationName := mustGetValidationNameArg()
	validationContext := mustGetValidationContextArg()

	token, err := getToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get token: %s\n", err)
		os.Exit(1)
	}

	hc := hetzner.NewClient(token)

	zone, err := hc.GetZoneByName(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get zone ID for %s: %s\n", domain, err)
		os.Exit(1)
	}

	_, err = hc.CreateRecord(validationName, 3600, "TXT", validationContext, zone.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create test record: %s\n", err)
		os.Exit(1)
	}
}

func Unset() {

	domain := mustGetDomainArg()
	validationName := mustGetValidationNameArg()
	validationContext := mustGetValidationContextArg()

	token, err := getToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get token: %s\n", err)
		os.Exit(1)
	}

	hc := hetzner.NewClient(token)

	zone, err := hc.GetZoneByName(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get zone ID for %s: %s\n", domain, err)
		os.Exit(1)
	}

	records, err := hc.GetAllRecordsByZone(zone.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create test record: %s\n", err)
		os.Exit(1)
	}

	recordId := ""

	for i := range records {
		if records[i].Name == validationName && records[i].Value == validationContext && records[i].Type == "TXT" {
			recordId = records[i].ID
			break
		}
	}

	if recordId == "" {
		fmt.Fprintf(os.Stderr, "Failed to get record ID: not found\n")
		os.Exit(1)
	}

	err = hc.DeleteRecord(recordId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete record: %s\n", err)
		os.Exit(1)
	}
}

func Init() {

	path := os.ExpandEnv("$HOME/.tahtoken")

	fmt.Printf("Creating %s...\n", path)

	_, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create %s: %s\n", path, err)
		os.Exit(1)
	}

	fmt.Printf("Change mode of %s to 0600...\n", path)
	err = os.Chmod(path, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to chmod %s: %s\n", path, err)
		os.Exit(1)
	}
}

func Test() {

	domain := mustGetDomainArg()

	err := testTokenFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to test token: %s\n", err)
		os.Exit(1)
	}

	token, err := getToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get token: %s\n", err)
		os.Exit(1)
	}

	hc := hetzner.NewClient(token)

	hetznerZone, err := hc.GetZoneByName(dns.GetDomain(domain))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get zone ID for %s: %s\n", dns.GetDomain(domain), err)
		os.Exit(1)
	}

	r, err := hetzner.NewClient(token).CreateRecord("hcdcadclk", 3600, "TXT", "TEST-TRUENAS-ACME-HETZNER", hetznerZone.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create test record: %s\n", err)
		os.Exit(1)
	}

	err = hetzner.NewClient(token).DeleteRecord(r.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete test record: %s\n", err)
		os.Exit(1)
	}
}

func Help() {

	fmt.Printf("Usage: %s <command> <domain> <validation_name> <validation_context>\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("Command:\n")
	fmt.Printf("\tset\n")
	fmt.Printf("\tunset\n")
	fmt.Printf("\tinit - Initialize script\n")
	fmt.Printf("\ttest - Test script configuration\n")
	fmt.Printf("\thelp - Print help\n")

}

func main() {

	switch mustGetCommandArg() {
	case "set":
		Set()
	case "unset":
		Unset()
	case "init":
		Init()
	case "test":
		Test()
	case "help":
		Help()
	}

}
