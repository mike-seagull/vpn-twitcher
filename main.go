package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/glendc/go-external-ip"
	"github.com/spf13/viper"
	"net"
	"os"
	"errors"
	"plugin"
	"strings"
)	

func validate_ip(actual_public_ip string) (string, error) {
	if govalidator.IsDNSName(actual_public_ip) {
		// get ips from domain name
		ips, err := net.LookupIP(actual_public_ip)
		if err != nil {
			return "", err
		}
		for _, ip := range ips {
			if govalidator.IsIPv4(ip.String()) {
				actual_public_ip = ip.String()
			}
		}
		if !govalidator.IsIPv4(actual_public_ip) {
			fmt.Println()
			return "", errors.New("Failed getting IPv4 from " + actual_public_ip)
		}		
	} else if govalidator.IsIPv4(actual_public_ip) {
		return actual_public_ip, nil
	} else if govalidator.IsIPv6(actual_public_ip) {
		return "", errors.New("IPv6 is currently not supported")
	} else {
		return "", errors.New("Invalid input: " + actual_public_ip)
	}
	return actual_public_ip, nil
}
func notify(notification_config map[string]string, message string) {
	for plugin_name, _ := range notification_config {
		if plugin_name != "up_message" && plugin_name != "down_message" {
			plug, err := plugin.Open(plugin_name+".so")
		        if err != nil {
                fmt.Println(err)
                os.Exit(1)
	        }
			symbol_notification_plugin, err := plug.Lookup("NotificationPlugin")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			var np NotificationPlugin
			np, ok := symbol_notification_plugin.(NotificationPlugin)
			if !ok {
				fmt.Println("unexpected type from module symbol")
				os.Exit(1)
			}
			np.Init(viper.GetStringMapString("notification."+plugin_name))
			np.SendMessage(message)	
		}
	}
}
func main() {
	var config_file string
	if len(os.Args) < 2 {
		config_file = "config.yaml"
	} else {
		config_file = os.Args[1]
	}
	fmt.Println(config_file)
	viper.SetConfigName(strings.TrimSuffix(config_file, ".yaml"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./test/")
	config_err := viper.ReadInConfig()
	if config_err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", config_err))
	}
	actual_public_ip := viper.GetString("ip_address")
	if actual_public_ip == "" {
		panic(fmt.Errorf("ip_address is required"))
	}
	ipv4, err := validate_ip(actual_public_ip)
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	// get public IP
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		panic(fmt.Errorf("failed getting public ip"))
	}
	if ip.String() == ipv4 {
		fmt.Println("VPN is down")
		if viper.IsSet("notification") { // load notification plugins
			notify(viper.GetStringMapString("notification"), viper.GetString("notification.down_message"))
		}
	} else {
		fmt.Println("VPN is active")
		if viper.IsSet("notification") {
			notify(viper.GetStringMapString("notification"), viper.GetString("notification.up_message"))
		}
	}
}
