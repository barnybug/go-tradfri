package tradfri

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/barnybug/go-tradfri/log"
	"github.com/zubairhamed/canopus"
)

type Client struct {
	Gateway string
	Key     string
	Ident   string
	PSK     string

	connection canopus.Connection
}

func SetDebug(debug bool) {
	log.Debug = debug
}

func NewClient(gateway string) *Client {
	return &Client{
		Gateway: gateway,
	}
}

func (c *Client) Connect() error {
	if c.PSK == "" {
		err := c.generatePSK()
		if err != nil {
			return err
		}
	}

	address := fmt.Sprintf("%s:%d", c.Gateway, tradfriPort)
	log.Printf("Connecting to gateway: %s", address)
	var err error
	c.connection, err = canopus.DialDTLS(address, c.Ident, c.PSK)
	if err != nil {
		return err
	}
	return nil
}

func pskPath() string {
	u, _ := user.Current()
	return path.Join(u.HomeDir, ".tradfri-psk")
}

func (c *Client) LoadPSK() error {
	file, err := os.Open(pskPath())
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = fmt.Fscanf(file, "%s\n%s", &c.Ident, &c.PSK)
	if err == nil {
		log.Println("Loaded PSK")
	} else {
		log.Printf("Couldn't load PSK: %s", err)
	}
	return err
}

func (c *Client) SavePSK() {
	s := fmt.Sprintf("%s\n%s", c.Ident, c.PSK)
	err := ioutil.WriteFile(pskPath(), []byte(s), 0700)
	if err != nil {
		log.Printf("Error saving PSK %s: %s", pskPath(), err)
	} else {
		log.Printf("Saved PSK %s", pskPath())
	}
}

func uuid() string {
	uuid, _ := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	return strings.TrimSpace(string(uuid))
}

func (c *Client) generatePSK() error {
	if c.Ident == "" {
		c.Ident = uuid()
		log.Printf("Generated ident: X%sX", c.Ident)
	} else {
		log.Printf("Using ident: %s", c.Ident)
	}
	log.Println("Requesting PSK...")
	address := fmt.Sprintf("%s:%d", c.Gateway, tradfriPort)
	conn, err := canopus.DialDTLS(address, preauthIdentity, c.Key)
	if err != nil {
		return err
	}

	payload := PSKRequest{Ident: c.Ident}
	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post)
	req.SetRequestURI(uriIdent)
	data := canopus.NewJSONPayload(payload).GetBytes()
	req.SetPayload(data)
	resp, err := conn.Send(req)
	if err != nil {
		return err
	}
	if resp.GetMessage().GetCode() == 65 {
		var pskResp PSKResponse
		err := json.Unmarshal(resp.GetMessage().GetPayload().GetBytes(), &pskResp)
		if err != nil {
			return err
		}
		c.PSK = pskResp.PSK
		log.Printf("PSK: %s\n", c.PSK)
		return nil
	} else {
		return errors.New("Unable to get PSK")
	}
}

func (c *Client) putRequest(uri string, payload interface{}) error {
	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Put)
	data := canopus.NewJSONPayload(payload).GetBytes()
	req.SetPayload(data)
	log.Printf("PUT %s payload %s", uri, string(data))
	req.SetRequestURI(uri)
	resp, err := c.connection.Send(req)
	if err != nil {
		log.Printf("<- error: %+v", err)
		return err
	}
	rdata := resp.GetMessage().GetPayload().GetBytes()
	log.Printf("<- %s", string(rdata))
	return nil
}

func (c *Client) postRequest(uri string) error {
	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post)
	req.SetStringPayload("")
	log.Printf("POST %s", uri)
	req.SetRequestURI(uri)
	resp, err := c.connection.Send(req)
	if err != nil {
		log.Printf("<- error: %+v", err)
		return err
	}
	rdata := resp.GetMessage().GetPayload().GetBytes()
	log.Printf("<- %s", string(rdata))
	return nil
}

func (c *Client) getRequest(uri string, out interface{}) error {
	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get)
	req.SetStringPayload("")
	log.Printf("GET %s", uri)
	req.SetRequestURI(uri)
	resp, err := c.connection.Send(req)
	if err != nil {
		log.Printf("<- error: %+v", err)
		return err
	}
	data := resp.GetMessage().GetPayload().GetBytes()
	log.Printf("<- %s", string(data))

	err = json.Unmarshal(data, out)
	return err
}

func (c *Client) GetGatewayInfo() (*GatewayInfo, error) {
	var gatewayInfo GatewayInfo
	err := c.getRequest(uriGatewayInfo, &gatewayInfo)
	return &gatewayInfo, err
}

func (c *Client) Reboot() error {
	return c.postRequest(uriGatewayReboot)
}

func (c *Client) FactoryReset() error {
	return c.postRequest(uriGatewayFactoryReset)
}

func (c *Client) ListDeviceIds() (deviceIds []int, err error) {
	log.Println("Looking for devices... ")
	err = c.getRequest(uriDevices, &deviceIds)
	return deviceIds, err
}

func (c *Client) ListDevices() (devices []*DeviceDescription, err error) {
	deviceIds, err := c.ListDeviceIds()
	if err != nil {
		return
	}

	log.Println("Enumerating...")
	for _, device := range deviceIds {
		var desc *DeviceDescription
		desc, err = c.GetDeviceDescription(device)
		if err != nil {
			return
		}
		log.Printf("Found device: %s\n", desc)
		devices = append(devices, desc)

		// sleep for a while to avoid flood protection
		time.Sleep(100 * time.Millisecond)
	}

	return
}

func (c *Client) GetDeviceDescription(id int) (*DeviceDescription, error) {
	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	var desc DeviceDescription
	err := c.getRequest(uri, &desc)
	return &desc, err
}

func (c *Client) SetDevice(deviceId int, change LightControl) error {
	payload := DeviceSet{
		[]LightControl{change},
	}
	uri := fmt.Sprintf("%s/%d", uriDevices, deviceId)
	return c.putRequest(uri, payload)
}

func (c *Client) ListGroups() (groups []*GroupDescription, err error) {
	log.Println("Requesting groups... ")
	var groupIds []int
	err = c.getRequest(uriGroups, &groupIds)
	if err != nil {
		return
	}

	log.Println("Enumerating...")
	for _, group := range groupIds {
		var desc *GroupDescription
		desc, err = c.GetGroupDescription(group)
		if err != nil {
			return
		}
		log.Println("Found group: %+v\n", desc)
		groups = append(groups, desc)

		// sleep for a while to avoid flood protection
		time.Sleep(100 * time.Millisecond)
	}

	return
}

func (c *Client) GetGroupDescription(id int) (*GroupDescription, error) {
	uri := fmt.Sprintf("%s/%d", uriGroups, id)
	var desc GroupDescription
	err := c.getRequest(uri, &desc)
	return &desc, err
}

func (c *Client) SetGroup(groupId int, change LightControl) error {
	payload := change
	uri := fmt.Sprintf("%s/%d", uriGroups, groupId)
	return c.putRequest(uri, payload)
}

func (c *Client) observer(in chan canopus.ObserveMessage, out chan *DeviceDescription) {
	for msg := range in {
		value := msg.GetValue()
		if value, ok := value.(canopus.MessagePayload); ok {
			dd := &DeviceDescription{}
			err := json.Unmarshal(value.GetBytes(), dd)
			if err == nil {
				out <- dd
			}
		}
	}
}

func (c *Client) Events() <-chan *DeviceDescription {
	out := make(chan *DeviceDescription, 16)
	in := make(chan canopus.ObserveMessage, 16)
	go c.connection.Observe(in)
	go c.observer(in, out)
	return out
}

func (c *Client) Observe(deviceId int) error {
	uri := fmt.Sprintf("%s/%d", uriDevices, deviceId)
	_, err := c.connection.ObserveResource(uri)
	return err
}
