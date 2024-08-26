package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"tailscale.com/client/tailscale"
)

func SubnetRouter(authkey string, routes string, tags string) {
	aerr := exec.Command("bash", "-c", "echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf").Run()
	if aerr != nil {
		log.Fatal(aerr)
	}

	berr := exec.Command("bash", "-c", "echo 'net.ipv6.conf.all.forwarding = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf").Run()
	if berr != nil {
		log.Fatal(berr)
	}

	cerr := exec.Command("bash", "-c", "sudo sysctl -p /etc/sysctl.d/99-tailscale.conf && sudo systemctl enable --now tailscaled").Run()
	if cerr != nil {
		log.Fatal(cerr)
	}

	command := fmt.Sprintf("sudo tailscale up --authkey %s --advertise-routes=%s --advertise-tags=%s", authkey, routes, tags)

	cfg := exec.Command("bash", "-c", command)
	err := cfg.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func AppConnector(authkey string, tags string) {
	aerr := exec.Command("bash", "-c", "echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf").Run()
	if aerr != nil {
		log.Fatal(aerr)
	}

	berr := exec.Command("bash", "-c", "echo 'net.ipv6.conf.all.forwarding = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf").Run()
	if berr != nil {
		log.Fatal(berr)
	}

	cerr := exec.Command("bash", "-c", "sudo sysctl -p /etc/sysctl.d/99-tailscale.conf && sudo systemctl enable --now tailscaled").Run()
	if cerr != nil {
		log.Fatal(cerr)
	}

	command := fmt.Sprintf("sudo tailscale up --authkey %s --advertise-connector --advertise-tags=%s", authkey, tags)

	cfg := exec.Command("bash", "-c", command)
	err := cfg.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Node(authkey string, tags string) {
	command := fmt.Sprintf("sudo tailscale up --authkey %s --advertise-tags=%s", authkey, tags)

	cfg := exec.Command("bash", "-c", command)
	err := cfg.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Install() {
	command := "curl -fsSL https://tailscale.com/install.sh | sh"

	cfg := exec.Command("bash", "-c", command)
	err := cfg.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func TsStatus() string {
	lclient := &tailscale.LocalClient{
		Dial:          nil,
		Socket:        "",
		UseSocketOnly: false,
	}

	status, terr := lclient.Status(context.Background())
	if terr != nil {
		log.Fatal(terr)
	}

	return status.BackendState
}

func TsDeviceId() string {
	lclient := &tailscale.LocalClient{
		Dial:          nil,
		Socket:        "",
		UseSocketOnly: false,
	}

	status, terr := lclient.Status(context.Background())
	if terr != nil {
		log.Fatal(terr)
	}

	return string(status.Self.ID)
}

func TsDelete(authkey string, deviceid string) {
	tailscale.I_Acknowledge_This_API_Is_Unstable = true

	client := tailscale.NewClient("-", tailscale.APIKey(authkey))

	client.DeleteDevice(context.Background(), deviceid)
}

func TsDevice(authkey string, deviceid string) *tailscale.Device {
	tailscale.I_Acknowledge_This_API_Is_Unstable = true

	client := tailscale.NewClient("-", tailscale.APIKey(authkey))

	details, err := client.Device(context.Background(), deviceid, nil)
	if err != nil {
		log.Fatal(err)

	}

	return details
}

func TsDevices(authkey string) []*tailscale.Device {
	tailscale.I_Acknowledge_This_API_Is_Unstable = true

	client := tailscale.NewClient("-", tailscale.APIKey(authkey))

	details, err := client.Devices(context.Background(), nil)
	if err != nil {
		log.Fatal(err)

	}

	return details
}

func Logout() {
	lclient := &tailscale.LocalClient{
		Dial:          nil,
		Socket:        "",
		UseSocketOnly: false,
	}

	err := lclient.Logout(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
