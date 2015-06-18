package command

import (
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

const (
	LockFlagValue = 0x2ddccbc058a50c18
)

type UnlockCommand struct {
	UI	cli.Ui
	Consul	*ConsulFlags
}

func (c *UnlockCommand) Help() string {
	helpText := `
Usage: consulkv unlock [options] path

  Release a lock on a given path

Options:

  -consul=127.0.0.1:8500	HTTP address of the Consul Agent
  -ssl				Use HTTPS while talking to Consul.
				(default: false)
  -ssl-verify			Verify certificates when connecting via SSL.
				(default: true)
  -ssl-cert			Path to an SSL certificate to use to authenticate
				to the consul server.
				(default: not set)
  -ssl-ca-cert			Path to an SSL client certificate to use to authenticate
				to the consul server.
				(default: not set)
  -token			The Consul API token.
				(default: not set)
  --session=<sessionId>		Session ID of the lock holder. Required
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *UnlockCommand) Run(args []string) int {
	var sessionId string

	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.StringVar(&sessionId, "session", "", "")
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if sessionId == "" {
		c.UI.Error("Session ID must be provided")
		c.UI.Error("")
		c.UI.Error(c.Help())
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

	consul, err := NewConsulClient(c.Consul, &c.UI)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.KV()

	kv := new(consulapi.KVPair)
	kv.Key = path
	kv.Session = sessionId
	kv.Flags = LockFlagValue

	success, _, err := client.Release(kv, nil)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if !success {
		return 1
	}

	return 0
}

func (c *UnlockCommand) Synopsis() string {
	return "Unlock a node"
}
