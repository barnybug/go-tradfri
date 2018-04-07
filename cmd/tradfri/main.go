package main

import (
	"errors"
	"fmt"
	"os"

	tradfri "github.com/barnybug/go-tradfri"
	"github.com/barnybug/go-tradfri/log"
	"github.com/urfave/cli"
)

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Error: '%+v'\n", err)
		os.Exit(1)
	}
}

func main() {
	commonFlags := []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable debug logging",
		},
		cli.StringFlag{
			Name:  "gateway",
			Usage: "hostname or IP (required)",
		},
		cli.StringFlag{
			Name:  "key",
			Usage: "gateway key (required)",
		},
	}
	app := cli.NewApp()
	app.Name = "tradfri"
	app.Usage = "Command line tool for the Ikea Tradfri gateway"
	app.Version = "0.0.1"
	app.Flags = commonFlags
	app.Commands = []cli.Command{
		{
			Name:   "devices",
			Usage:  "scan for devices",
			Action: devicesCommand,
		},
		{
			Name:   "groups",
			Usage:  "scan for groups",
			Action: groupsCommand,
		},
		{
			Name:   "set",
			Usage:  "switch/dim/color a device or group",
			Action: setCommand,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "id",
					Usage: "device or group id",
				},
				cli.BoolFlag{
					Name:  "off",
					Usage: "switch off",
				},
				cli.IntFlag{
					Name:  "level",
					Usage: "dim level (0-100)",
				},
				cli.IntFlag{
					Name:  "temp",
					Usage: "colour temperature (2200-4000K)",
				},
				cli.StringFlag{
					Name:  "color",
					Usage: "hex string (6 chars)",
				},
				cli.IntFlag{
					Name:  "duration",
					Usage: "transition duration (ms)",
				},
			},
		},
		{
			Name:   "info",
			Usage:  "get gateway info",
			Action: infoCommand,
		},
		{
			Name:   "reboot",
			Usage:  "reboot gateway",
			Action: rebootCommand,
		},
		{
			Name:   "watch",
			Usage:  "watch for events",
			Action: watchCommand,
		},
		{
			Name:   "factory_reset",
			Usage:  "factory reset the gateway",
			Action: factoryResetCommand,
		},
	}
	err := app.Run(os.Args)
	checkErr(err)
	log.Println("Done")
}

func connect(c *cli.Context) (*tradfri.Client, error) {
	log.Debug = c.GlobalBool("debug")
	gateway := c.GlobalString("gateway")
	if gateway == "" {
		return nil, errors.New("--gateway required")
	}
	key := c.GlobalString("key")
	client := tradfri.NewClient(gateway)
	err := client.LoadPSK()
	if err != nil {
		if key == "" {
			return nil, errors.New("--key required")
		} else {
			client.Key = key
		}
	}
	err = client.Connect()
	if err != nil {
		client.SavePSK()
	}
	return client, err
}

func devicesCommand(c *cli.Context) error {
	client, err := connect(c)
	checkErr(err)

	devices, err := client.ListDevices()
	checkErr(err)

	fmt.Printf("Found %d devices\n\n", len(devices))
	for _, device := range devices {
		fmt.Printf("%s\n", device)
	}
	return nil
}

func setCommand(c *cli.Context) error {
	power := 1
	if c.BoolT("off") {
		power = 0
	}
	change := tradfri.LightControl{}
	change.Power = &power
	if c.IsSet("level") {
		dim := tradfri.PercentageToDim(c.Int("level"))
		change.Dim = &dim
	}
	if c.IsSet("color") {
		color := c.String("color")
		change.Color = &color
	}
	if c.IsSet("temp") {
		mired := tradfri.KelvinToMired(c.Int("temp"))
		change.Mireds = &mired
	}
	if c.IsSet("duration") {
		d := tradfri.MsToDuration(c.Int("duration"))
		change.Duration = &d
	}

	if !c.IsSet("id") {
		return errors.New("required arguments: --id")
	}
	client, err := connect(c)
	checkErr(err)
	id := c.Int("id")
	if id&(1<<17) == 0 {
		err = client.SetDevice(id, change)
	} else {
		err = client.SetGroup(id, change)
	}
	checkErr(err)
	return nil
}

func groupsCommand(c *cli.Context) error {
	client, err := connect(c)
	checkErr(err)

	groups, err := client.ListGroups()
	checkErr(err)

	for _, group := range groups {
		fmt.Printf("%s\n", group)
	}
	return nil
}

func watchCommand(c *cli.Context) error {
	client, err := connect(c)
	checkErr(err)

	fmt.Println("Listing devices...")
	deviceIds, err := client.ListDeviceIds()
	checkErr(err)

	fmt.Printf("Watching %d devices...\n", len(deviceIds))
	for _, id := range deviceIds {
		client.Observe(id)
	}

	for msg := range client.Events() {
		fmt.Printf("%s\n", msg)
	}
	return nil
}

func infoCommand(c *cli.Context) error {
	client, err := connect(c)
	checkErr(err)

	info, err := client.GetGatewayInfo()
	checkErr(err)

	fmt.Printf("%s\n", info)
	return nil
}

func rebootCommand(c *cli.Context) error {
	client, err := connect(c)
	checkErr(err)

	err = client.Reboot()
	checkErr(err)
	return nil
}

func factoryResetCommand(c *cli.Context) error {
	client, err := connect(c)
	checkErr(err)

	err = client.FactoryReset()
	checkErr(err)
	return nil
}
