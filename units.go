package main

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/coreos/go-log/log"
	"github.com/coreos/go-systemd/dbus"
	"io/ioutil"
	"strings"
	"time"
)


func main() {
	client := etcd.NewClient()
	bootID, err := getBootID()
	if err != nil {
		log.Fatalln(err)
	}

	keyPrefix := "/system/" + bootID + "/"
	statusChan, _ := dbus.SubscribeUnits(5 * time.Second)
	for statusMap := range statusChan {
		for name, unit := range statusMap {
			key := keyPrefix + name + "/status"
			if unit != nil {
				client.Set(key, unit.SubState, 0)
			} else {
				client.Delete(key)
			}
		}
	}
}

func getBootID() (string, error) {
	contents, err := ioutil.ReadFile("/proc/sys/kernel/random/boot_id")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contents)), nil
}
