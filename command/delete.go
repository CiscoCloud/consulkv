package command

import (
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type DeleteCommand struct {
	UI	cli.Ui
	Consul	*ConsulFlags
}

func (c *DeleteCommand) Help() string {
	helpText := `
Usage: consulkv delete [options] path

  Delete a given path

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
  --modifyindex=<ModifyIndex>	Perform a Check-and-Set delete
				(default: not set)
  --recurse			Perform a recursive delete
				(default: false)
`

	return strings.TrimSpace(helpText)
}

func (c *DeleteCommand) Run (args[]string) int {
	var modifyIndex string
	var doRecurse bool

	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.StringVar(&modifyIndex, "modifyindex", "", "")
	cmdFlags.BoolVar(&doRecurse, "recurse", false, "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	extra := cmdFlags.Args()
	if len(extra) < 1 {
		c.UI.Error("Key path must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}

	path := extra[0]

	writeOpts := new(consulapi.WriteOptions)
	consul, err := NewConsulClient(c.Consul, &c.UI)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.KV()

	switch {
	case doRecurse:
		_, err := client.DeleteTree(path, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	case modifyIndex != "":
		m, err := strconv.ParseUint(modifyIndex, 0, 64)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		kv := consulapi.KVPair{
			Key:		path,
			ModifyIndex:	m,
		}

		success, _, err := client.DeleteCAS(&kv, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		if !success {
			return 1
		}
	default:
		_, err := client.Delete(path, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	}

	return 0
}

func (c *DeleteCommand) Synopsis() string {
	return "Delete a path"
}
