package tradfri

import (
	"fmt"
	"time"
)

// types to unmarshal json data from tradfri
// generated with help from https://mholt.github.io/json-to-go/
// struct names derived from
// - https://github.com/IPSO-Alliance/pub/blob/master/reg/README.md
// - https://github.com/hardillb/TRADFRI2MQTT/blob/master/src/main/java/uk/me/hardill/TRADFRI2MQTT/TradfriConstants.java
// - http://www.openmobilealliance.org/wp/OMNA/LwM2M/LwM2MRegistry.html#resources

type LightControl struct {
	Color    *string `json:"5706,omitempty"`
	ColorHue *int    `json:"5707,omitempty"`
	ColorSat *int    `json:"5708,omitempty"`
	ColorX   *int    `json:"5709,omitempty"`
	ColorY   *int    `json:"5710,omitempty"`
	Power    *int    `json:"5850,omitempty"`
	Dim      *int    `json:"5851,omitempty"`
	Mireds   *int    `json:"5711,omitempty"`
	Duration *int    `json:"5712,omitempty"`
}

type DeviceDescription struct {
	Device struct {
		Manufacturer          string `json:"0"`
		ModelNumber           string `json:"1"`
		Serial                string `json:"2"`
		FirmwareVersion       string `json:"3"`
		AvailablePowerSources int    `json:"6"`
		BatteryLevel          int    `json:"9"`
	} `json:"3"`
	LightControl      []LightControl `json:"3311"`
	ApplicationType   int            `json:"5750"`
	DeviceName        string         `json:"9001"`
	CreatedAt         int            `json:"9002"`
	DeviceID          int            `json:"9003"`
	ReachabilityState int            `json:"9019"`
	LastSeen          int            `json:"9020"`
	OTAUpdateState    int            `json:"9054"`
}

var PowerSources = map[int]string{
	1: "Internal Battery",
	2: "External Battery",
	3: "Battery",
	4: "Power over Ethernet",
	5: "USB",
	6: "AC (Mains) power",
	7: "Solar",
}

func (d *DeviceDescription) AvailablePowerSource() string {
	if s, ok := PowerSources[d.Device.AvailablePowerSources]; ok {
		return s
	} else {
		return "Unknown"
	}
}

func (d *DeviceDescription) String() string {
	s := fmt.Sprintf("ID: %d Name: %q\nType: %d Model: %q\n", d.DeviceID, d.DeviceName, d.ApplicationType, d.Device.ModelNumber)
	s += fmt.Sprintf("Firmware: %s Manufacturer: %q\n", d.Device.FirmwareVersion, d.Device.Manufacturer)
	s += fmt.Sprintf("Power: %s", d.AvailablePowerSource())
	if d.ApplicationType == Remote || d.ApplicationType == Remote2 {
		s += fmt.Sprintf(" Level: %v%%", d.Device.BatteryLevel)
	}
	s += "\n"
	lastSeen := time.Unix(int64(d.LastSeen), 0)
	s += fmt.Sprintf("Last seen: %s\n", lastSeen.Format(time.RFC1123))
	if d.ApplicationType == Lamp {
		for count, entry := range d.LightControl {
			power := "off"
			if *entry.Power != 0 {
				power = "on"
			}
			pc := DimToPercentage(*entry.Dim)
			s += fmt.Sprintf("Light Control Set %d, Power: %s, Dim: %d%%\n",
				count, power, pc)
			s += "Color: "
			if entry.Mireds != nil {
				s += fmt.Sprintf("%dK ", MiredToKelvin(*entry.Mireds))
			}
			if entry.Color != nil {
				s += fmt.Sprintf("#%s ", *entry.Color)
			}
			if entry.ColorX != nil {
				s += fmt.Sprintf("X:%d/Y:%d ", *entry.ColorX, *entry.ColorY)
			}
			if entry.ColorHue != nil {
				s += fmt.Sprintf("Hue: %d Sat: %d ", *entry.ColorHue, *entry.ColorSat)
			}
			s += "\n"
		}
	}
	return s
}

func (d *DeviceDescription) SupportsMired() bool {
	return d.LightControl[0].Mireds != nil
}

func (d *DeviceDescription) SupportsColorXY() bool {
	return d.LightControl[0].ColorX != nil
}

func (d *DeviceDescription) SupportsHueSat() bool {
	return d.LightControl[0].ColorHue != nil
}

type DeviceSet struct {
	LightControl []LightControl `json:"3311"`
}

type GroupDescription struct {
	Power         int    `json:"5850"`
	Dim           int    `json:"5851"`
	GroupName     string `json:"9001"`
	CreatedAt     int    `json:"9002"`
	GroupID       int    `json:"9003"`
	AccessoryLink struct {
		LinkedItems struct {
			DeviceIDs []int `json:"9003"`
		} `json:"15002"`
	} `json:"9018"`
	Num9039 int `json:"9039"`
}

func (g *GroupDescription) String() string {
	createdAt := time.Unix(int64(g.CreatedAt), 0)
	s := fmt.Sprintf("ID: %d Name: %q Created: %s\n", g.GroupID, g.GroupName, createdAt.Format(time.RFC1123))
	s += fmt.Sprintf("Power: %d Dim: %d\n", g.Power, g.Dim)
	s += fmt.Sprintf("Linked devices: %v\n", g.AccessoryLink.LinkedItems.DeviceIDs)
	return s
}

type PSKRequest struct {
	Ident string `json:"9090"`
}

type PSKResponse struct {
	PSK string `json:"9091"`
}

type GatewayInfo struct {
	ID               string `json:"9081"`
	NTPServer        string `json:"9023"`
	FirmwareVersion  string `json:"9029"`
	CurrentTimestamp int    `json:"9059"`
	CurrentTime      string `json:"9060"`
}

func (g *GatewayInfo) String() string {
	return fmt.Sprintf("ID: %s\nNTPServer: %s\nFirmware Version: %s\nCurrent Time: %s", g.ID, g.NTPServer, g.FirmwareVersion, g.CurrentTime)
}
