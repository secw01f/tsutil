package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

var (
	app           = kingpin.New("tsutil", "A tool for installing Tailscale and adding Nodes, Subnet Routers, and App Connectors using the Tailscale OAuth client")
	install       = app.Command("install", "Install Tailscale on the device")
	tscs          = app.Flag("secret", "Tailscale client secret key").String()
	smarn         = app.Flag("arn", "ARN for AWS SecretsManager secret for Tailscale client id and secret").String()
	subnet        = app.Command("subnetrouter", "Configures endpoint to be a Tailscale Subnet Router (Auth Required)")
	subnetroutes  = subnet.Flag("routes", "Route to advertise for Tailscale Subnet Router").Required().String()
	subnettags    = subnet.Flag("tags", "Comma seperated tags required by the OAuth key for Subnet Router config)").Default("tag:subnetrouter").String()
	connector     = app.Command("app-connector", "Configures endpoint as a Tailscale App Connector (Auth Required)")
	connectortags = connector.Flag("tags", "Comma seperated tags required by the OAuth key for App Connector config)").Default("tag:connector").String()
	node          = app.Command("node", "Configures endpoint as a standard Tailscale Node (Auth Required)")
	nodetags      = node.Flag("tags", "Comma seperated tags required by the OAuth key for Node config)").Required().String()
	status        = app.Command("status", "Check the status of Tailscaled")
	deviceid      = app.Command("id", "Get the device ID of this Tailscale endpoint")
	delete        = app.Command("delete", "Delete a Tailscale endpoint (Auth Required)")
	deleteself    = delete.Flag("self", "Delete the this Tailscale endpoint (Auth Required)").Bool()
	deleteid      = delete.Flag("id", "Device ID to delete (Auth Required)").String()
	device        = app.Command("device", "Details about a device (Auth Required)")
	devicejson    = device.Flag("json", "Output device details in JSON (Auth Required)").Bool()
	detailself    = device.Flag("self", "Details about this device (Auth Required)").Bool()
	detailsid     = device.Flag("id", "Device ID to gather details (Auth Required)").String()
	devices       = app.Command("devices", "List all devices on the tailnet (Auth Required)")
	devicesjson   = devices.Flag("json", "Output list of all devices on the tialnet in JSON (Auth Required)").Bool()
	logout        = app.Command("logout", "Logout of Tailscale")
)

var authkey string

type Device struct {
	NodeID        string
	Username      string
	Hostname      string
	OS            string
	ClientVersion string
	Update        bool
	Tags          []string
}

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
			if *devicejson {
				devicedetails := Device{NodeID: details.NodeID, Username: details.User, Hostname: details.Hostname, OS: details.OS, ClientVersion: details.ClientVersion, Update: details.UpdateAvailable, Tags: details.Tags}
				devicedata, err := json.MarshalIndent(devicedetails, "", "    ")
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(devicedata))
			} else {
				fmt.Println("DeviceID:", details.NodeID)
				fmt.Println("Username:", details.User)
				fmt.Println("Hostname:", details.Hostname)
				fmt.Println("OS:", details.OS)
				fmt.Println("Client Version:", details.ClientVersion)
				fmt.Println("Update:", details.UpdateAvailable)
				fmt.Println("Tags:", details.Tags)
			}
		} else {
			tscid := strings.Split(*tscs, "-")[2]
			authkey = GetTsApi(tscid, *tscs)
			details := TsDevice(authkey, *detailsid)
			if *devicejson {
				devicedetails := Device{NodeID: details.NodeID, Username: details.User, Hostname: details.Hostname, OS: details.OS, ClientVersion: details.ClientVersion, Update: details.UpdateAvailable, Tags: details.Tags}
				devicedata, err := json.MarshalIndent(devicedetails, "", "    ")
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(devicedata))
			} else {
				fmt.Println("DeviceID:", details.NodeID)
				fmt.Println("Username:", details.User)
				fmt.Println("Hostname:", details.Hostname)
				fmt.Println("OS:", details.OS)
				fmt.Println("Client Version:", details.ClientVersion)
				fmt.Println("Update:", details.UpdateAvailable)
				fmt.Println("Tags:", details.Tags)
			}
		}
	case devices.FullCommand():
		type Devices struct {
			Devices []Device `json:"devices"`
		}
		tscid := strings.Split(*tscs, "-")[2]
		authkey = GetTsApi(tscid, *tscs)
		devices := TsDevices(authkey)
		if *devicesjson {
			jsondevices := Devices{Devices: []Device{}}
			for _, device := range devices {
				devicedetails := Device{NodeID: device.NodeID, Username: device.User, Hostname: device.Hostname, OS: device.OS, ClientVersion: device.ClientVersion, Update: device.UpdateAvailable, Tags: device.Tags}
				jsondevices.Devices = append(jsondevices.Devices, devicedetails)
			}
			devicesdata, err := json.MarshalIndent(jsondevices, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(devicesdata))
		} else {
			for _, device := range devices {
				fmt.Println("\nDeviceID:", device.NodeID)
				fmt.Println("Username:", device.User)
				fmt.Println("Hostname:", device.Hostname)
				fmt.Println("OS:", device.OS)
				fmt.Printf("Client Version: %s\n", device.ClientVersion)
				fmt.Println("Update:", device.UpdateAvailable)
				fmt.Printf("Tags: %s\n", device.Tags)
			}
		}
	case logout.FullCommand():
		Logout()
	}
}
