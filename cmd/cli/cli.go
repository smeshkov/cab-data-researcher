package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/smeshkov/cab-data-researcher/app"
)

var (
	baseURL = "http://localhost:8080/api/v1"

	client = &http.Client{
		Timeout: time.Second * 120,
	}

	// trips command
	tripsCmd        = flag.NewFlagSet("trips", flag.ExitOnError)
	tripsMedallions = tripsCmd.String("medallions", "D7D598CD99978BD012A87A76A7C891B7", "Comma separated list of cab medallions.")
	tripsPickupDate = tripsCmd.String("pickupDate", "2013-12-01", "Pickup date.")
	tripsMoCache    = tripsCmd.Bool("noCache", false, "Do not use cache.")

	// cache command
	cacheCmd   = flag.NewFlagSet("cache", flag.ExitOnError)
	cacheClear = cacheCmd.Bool("clr", true, "Clears cache.")
)

func main() {
	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	for strings.HasPrefix(cmd, "-") {
		cmd = strings.TrimPrefix(cmd, "-")
	}
	if cmd == "help" || cmd == "h" {
		printUsage()
		os.Exit(1)
	}

	switch cmd {
	case tripsCmd.Name():
		getTripsCount(tripsCmd)
	case cacheCmd.Name():
		clearCache(cacheCmd)
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		os.Exit(2)
	}
}
func getTripsCount(cmd *flag.FlagSet) {
	parseCommand(cmd)
	req := &app.TripCountReq{
		Medallions: strings.Split(*tripsMedallions, ","),
		PickupDate: *tripsPickupDate,
		NoCache:    *tripsMoCache,
	}
	bs, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	r, _ := http.NewRequest(http.MethodPost, baseURL+"/trip/count", bytes.NewReader(bs))
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("bad response status: %s\n", resp.Status)
		return
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
}

func clearCache(cmd *flag.FlagSet) {
	parseCommand(cmd)

	if !*cacheClear {
		fmt.Println("noop")
	}

	r, _ := http.NewRequest(http.MethodPost, baseURL+"/cache/clear", nil)
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		fmt.Printf("bad response status: %s\n", resp.Status)
		return
	}
	fmt.Println("cleared cache")
}

func printDefaults(cmd *flag.FlagSet) {
	println(cmd.Name())
	cmd.PrintDefaults()
}

func printUsage() {
	println("Usage:")
	println()
	printDefaults(tripsCmd)
	println()
	printDefaults(cacheCmd)
	println()
}

func parseCommand(cmd *flag.FlagSet) {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("error in parsing arguments: %v \n", err)
		printDefaults(cmd)
		os.Exit(3)
	}
}
