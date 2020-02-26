package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/olekukonko/tablewriter"
	"strconv"
	"strings"
)

func main() {

	IP := flag.String("ip", "", "IP or hostname of the switch")
	Community := flag.String("community", "public", "SNMP Community of the switch")
	NoAlias := flag.Bool("noalias", false, "Don't output Alias Values")
	flag.Parse()
	if *IP == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}
	a := HPSwitch_new(*IP, *Community)
	table := tablewriter.NewWriter(os.Stdout)

	var data [][]string

	fmt.Printf("Switch %s at IP %s in Location %s\n", a.Name, a.IP, a.Location)
	for _, port := range a.ports {
		row := []string{strconv.Itoa(port.ID) + " (" + port.Name + ")"}
		if ! *NoAlias {
			row = append(row, port.Description)
		}


		untagged := ""
		for _, vlan := range port.UntaggedVlan {
			untagged += vlan.vlan_name + ", "
		}


		for _, vlan := range port.TaggedVlans {
			untagged += vlan.vlan_name + "*, "
		}
		row = append(row, strings.TrimRight(untagged, ", "))

		data = append(data, row)

	}
	header := []string{"Port ID", "Port Alias", "VLANS"}
	if *NoAlias {
		header = []string{"Port ID", "VLANS"}
	}
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render()
}
