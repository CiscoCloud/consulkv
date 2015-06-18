package command

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type WriteCommand struct {
	UI		cli.Ui
	Consul		*ConsulFlags
}

func (c *WriteCommand) Help() string {
	helpText := `
Usage: consulkv write [options] path value

  Write a value to a given path.

Options:

  --consul=127.0.0.1:8500	HTTP address of the Consul Agent
  --ssl				Use HTTPS while talking to Consul.
				(default: false)
  --ssl-verify			Verify certificates when connecting via SSL.
				(default: true)
  --ssl-cert			Path to an SSL certificate to use to authenticate
				to the consul server.
				(default: not set)
  --ssl-ca-cert			Path to an SSL client certificate to use to authenticate
				to the consul server.
				(default: not set)
  --token			The Consul API token.
				(default: not set)
  --modifyindex=<ModifyIndex>	Perform a Check-and-Set write.
				(default: not set)
  --flags=<number>		Integer value between 0 and 2^64 - 1
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *WriteCommand) Run(args []string) int {
	var modifyIndex string
	var flags string

	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.StringVar(&modifyIndex, "cas", "", "")
	cmdFlags.StringVar(&flags, "flags", "", "")
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	extra := cmdFlags.Args()
	if len(extra) < 2 {
		c.UI.Error("Key path and value must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}

	path := extra[0]
	value := strings.Join(extra[1:], " ")

	kv := new(consulapi.KVPair)

	kv.Key = path
	if strings.HasPrefix(value, "@") {
		v, err := ioutil.ReadFile(value[1:])
		if err != nil {
			c.UI.Error(fmt.Sprintf("ReadFile error: %v", err))
			return 1
		}
		kv.Value = v
	} else {
		kv.Value = []byte(value)
	}

	// &flags=
	//
	if flags != "" {
		f, err := strconv.ParseUint(flags, 0, 64)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing flags: %v", flags))
			c.UI.Error("")
			return 1
		}
		kv.Flags = f
	}

	consul, err := NewConsulClient(c.Consul, &c.UI)
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.KV()

	if modifyIndex == "" {
		_, err := client.Put(kv, nil)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	} else {
		// Check-and-Set
		i, err := strconv.ParseUint(modifyIndex, 0, 64)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing modifyIndex: %v", modifyIndex))
			c.UI.Error("")
			return 1
		}
		kv.ModifyIndex = i

		success, _, err := client.CAS(kv, nil)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		if !success {
			return 1
		}
	}


	return 0
}

func (c *WriteCommand) Synopsis() string {
	return "Write a value"
}
