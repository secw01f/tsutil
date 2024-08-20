package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

var (
	app           = kingpin.New("tsutil", "A tool for installing Tailscale and adding Nodes, Subnet Routers, and App Connectors using the Tailscale OAuth client")
	install       = app.Command("install", "Install Tailscale on the device")
	tscs          = app.Flag("secret", "Tailscale client secret key").String()
	smarn         = app.Flag("arn", "ARN for AWS SecretsManager secret for Tailscale client id and secret").String()
	subnet        = app.Command("subnetrouter", "Configures endpoint to be a Tailscale Subnet Router")
	subnetroutes  = subnet.Flag("routes", "Route to advertise for Tailscale Subnet Router").Required().String()
	subnettags    = subnet.Flag("tags", "Comma seperated tags required by the OAuth key for Subnet Router config)").Default("tag:subnetrouter").String()
	connector     = app.Command("app-connector", "Configures endpoint as a Tailscale App Connector")
	connectortags = connector.Flag("tags", "Comma seperated tags required by the OAuth key for App Connector config)").Default("tag:connector").String()
	node          = app.Command("node", "Configures endpoint as a standard Tailscale Node")
	nodetags      = node.Flag("tags", "Comma seperated tags required by the OAuth key for Node config)").Required().String()
	status        = app.Command("status", "Check the status of Tailscaled")
	deviceid      = app.Command("id", "Get the device ID of this Tailscale endpoint")
	delete        = app.Command("delete", "Delete a Tailscale endpoint")
	deleteself    = delete.Flag("self", "Delete the this Tailscale endpoint").Bool()
	deleteid      = delete.Flag("id", "Device ID to delete").String()
	device        = app.Command("device", "Details about a device")
	detailself    = device.Flag("self", "Details about this device").Bool()
	detailsid     = device.Flag("id", "Device ID").String()
	logout        = app.Command("logout", "Logout of Tailscale")
)

var authkey string

func main() {

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case install.FullCommand():
		Install()
	case subnet.FullCommand():
		if *subnetroutes == "" {
			fmt.Println("You must provide a route to advertise for the Subnet Router")
			os.Exit(1)
		}

		if *tscs != "" && *smarn == "" {
			tscid := strings.Split(*tscs, "-")[2]
			tags := fmt.Sprintf("tag:subnetrouter,%s", *subnettags)
			authkey = GetTsAuth(tscid, *tscs, tags)
		} else if *smarn != "" && *tscs == "" {
			key := GetSMSecret(*smarn)
			id := strings.Split(key, "-")[2]
			tags := fmt.Sprintf("tag:subnetrouter,%s", *subnettags)
			authkey = GetTsAuth(id, key, tags)
		} else {
			os.Exit(1)
		}
		SubnetRouter(authkey, *subnetroutes, *subnettags)
	case connector.FullCommand():
		if *tscs != "" && *smarn == "" {
			tscid := strings.Split(*tscs, "-")[2]
			tags := "tag:connector"
			authkey = GetTsAuth(tscid, *tscs, tags)
		} else if *smarn != "" && *tscs == "" {
			key := GetSMSecret(*smarn)
			id := strings.Split(key, "-")[2]
			tags := "tag:connector"
			authkey = GetTsAuth(id, key, tags)
		} else {
			os.Exit(1)
		}
		AppConnector(authkey, *connectortags)
	case node.FullCommand():
		if *tscs != "" && *smarn == "" {
			tscid := strings.Split(*tscs, "-")[2]
			tags := *nodetags
			authkey = GetTsAuth(tscid, *tscs, tags)
		} else if *smarn != "" && *tscs == "" {
			key := GetSMSecret(*smarn)
			id := strings.Split(key, "-")[2]
			tags := *nodetags
			authkey = GetTsAuth(id, key, tags)
		} else {
			os.Exit(1)
		}
		Node(authkey, *nodetags)
	case status.FullCommand():
		fmt.Println("Tailscaled is:", TsStatus())
	case deviceid.FullCommand():
		fmt.Println("Device ID:", TsDeviceId())
	case delete.FullCommand():
		if *deleteself {
			tscid := strings.Split(*tscs, "-")[2]
			authkey = GetTsApi(tscid, *tscs)
			TsDelete(authkey, TsDeviceId())
		} else {
			tscid := strings.Split(*tscs, "-")[2]
			authkey = GetTsApi(tscid, *tscs)
			TsDelete(authkey, *deleteid)
		}
	case device.FullCommand():
		if *detailself {
			tscid := strings.Split(*tscs, "-")[2]
			authkey = GetTsApi(tscid, *tscs)
			details := TsDevice(authkey, TsDeviceId())
			fmt.Println(*details)
		} else {
			tscid := strings.Split(*tscs, "-")[2]
			authkey = GetTsApi(tscid, *tscs)
			details := TsDevice(authkey, *detailsid)
			fmt.Println(*details)
		}
	case logout.FullCommand():
		Logout()
	}
}
