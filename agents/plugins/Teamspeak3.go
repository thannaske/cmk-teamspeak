package main

import (
	"fmt"
	"github.com/multiplay/go-ts3"
	"gopkg.in/ini.v1"
	"os"
)

type ServerQueryConfig struct {
	ServerAddress string
	Username      string
	Password      string
}

// Set to TRUE during testing and make sure that you define the dev_* variables
// The if-structure gets automagically removed by the golang compiler during optimization
// because DEVELOPER_MODE is a const that can't change during runtime
const DEVELOPER_MODE bool = false

// Only relevant when testing
var (
	dev_server   = ""
	dev_username = ""
	dev_password = ""
)

func ReadConfig() *ServerQueryConfig {
	// Gets removed during build when DEVELOPER_MODE is set to false
	// Do not worry about loss of performance - Go compiler is awesome!
	if DEVELOPER_MODE {
		return &ServerQueryConfig{ServerAddress: dev_server, Username: dev_username, Password: dev_password}
	}

	// Check_MK provides an env variable to the configuration directory
	configBaseDir := os.Getenv("MK_CONFDIR")

	// Load the user configuration and process the required sections and keys
	// Exit application with non-zero code when configuration was not read successful
	// Also display a Check_MK Agent-compatible error code to process on monitoring server
	if cfg, err := ini.Load(configBaseDir + "/teamspeak3.cfg"); err != nil {
		fmt.Println("<<<Teamspeak3>>>")
		fmt.Println("ConfigError: Yes, 1")
		os.Exit(1)
	} else {
		if sect, err := cfg.GetSection("serverquery"); err != nil {
			fmt.Println("<<<Teamspeak3>>>")
			fmt.Println("ConfigError: Yes, 2")
			os.Exit(1)
		} else {
			var conf_address, conf_user, conf_password string

			if val, err := sect.GetKey("address"); err != nil {
				fmt.Println("<<<Teamspeak3>>>")
				fmt.Println("ConfigError: Yes, 3")
				os.Exit(1)
			} else {
				conf_address = val.String()
			}

			if val, err := sect.GetKey("user"); err != nil {
				fmt.Println("<<<Teamspeak3>>>")
				fmt.Println("ConfigError: Yes, 4")
				os.Exit(1)
			} else {
				conf_user = val.String()
			}

			if val, err := sect.GetKey("password"); err != nil {
				fmt.Println("<<<Teamspeak3>>>")
				fmt.Println("ConfigError: Yes, 5")
				os.Exit(1)
			} else {
				conf_password = val.String()
			}

			// Build ServerQueryConfig struct and return pointer
			return &ServerQueryConfig{ServerAddress: conf_address, Username: conf_user, Password: conf_password}
		}
	}

	// Return empty ServerQueryConfig struct
	return &ServerQueryConfig{}
}

func main() {
	var queryConfig *ServerQueryConfig = ReadConfig()

	// Print Check_MK section header
	fmt.Println("<<<Teamspeak3>>>")
	fmt.Println("ConfigError: No")

	// Establish connection to Teamspeak3 server query
	c, err := ts3.NewClient(queryConfig.ServerAddress)

	// Determine if we can actually reach the server's query console
	// In case of error we want to exit with zero-code because the check application itself ran correctly
	if err != nil {
		fmt.Println("QueryPortReachable: No")
		os.Exit(0)
	} else {
		fmt.Println("QueryPortReachable: Yes")
	}

	// Make sure the query connection will be closed when application terminates
	defer c.Close()

	// Try to authenticate with Teamspeak3 server query
	if err := c.Login(queryConfig.Username, queryConfig.Password); err != nil {
		fmt.Println("AuthSuccess: No")
		os.Exit(0)
	} else {
		fmt.Println("AuthSuccess: Yes")
	}

	// Try to get server's current version
	if v, err := c.Version(); err != nil {
		fmt.Println("Version: None")
		fmt.Println("Platform: None")
		fmt.Println("Build: None")
		os.Exit(0)
	} else {
		fmt.Println("Version:", v.Version)
		fmt.Println("Platform:", v.Platform)
		fmt.Println("Build:", v.Build)
	}

	// Iterate through list of virtual servers
	if l, err := c.Server.List(); err != nil {
		os.Exit(0)
	} else {
		for _, server := range l {
			c.Use(server.ID)

			var serverAutoStart string = "no"
			var trafficIngressBytesTotal uint64 = 0
			var trafficEgressBytesTotal uint64 = 0

			// Convert boolean value to string
			if server.AutoStart == true {
				serverAutoStart = "yes"
			}

			// When the server is stopped this method exits with an error
			if connInfo, err := c.Server.ServerConnectionInfo(); err == nil {
				trafficIngressBytesTotal = connInfo.BytesReceivedTotal
				trafficEgressBytesTotal = connInfo.BytesSentTotal
			}

			// Scheme: "VirtualServer: ($PORT $STATUS $ONLINE_CLIENTS $MAX_CLIENTS $CURRENT_CHANNELS $AUTO_START $BANDWIDTH_INGRESS_TOTAL $BANDWIDTH_EGRESS_TOTAL )"
			fmt.Printf("VirtualServer: (%d %s %d %d %d %s %d %d)\n", server.Port, server.Status, server.ClientsOnline, server.MaxClients, server.ChannelsOnline, serverAutoStart, trafficIngressBytesTotal, trafficEgressBytesTotal)
		}
	}

	// Ran successful
	os.Exit(0)
}
