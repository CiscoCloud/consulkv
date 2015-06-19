package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

type ReadCommand struct {
	UI	cli.Ui
	Consul	*ConsulFlags
}

func (c *ReadCommand) Help() string {
	helpText := `
Usage: consulkv read [options] path

  Read a value from a given path.

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
  --fields=value		Comma separated list of fields to return.
				(default: value)
  --format=text			Output format. Supported options: text, json, prettyjson
				(default: text)
  --delimiter=			Output field delimited.
				(default: " ")
  --header			Output a header row for text format
				(default: false)
`

	return strings.TrimSpace(helpText)
}

func (c *ReadCommand) Run(args []string) int {
	var format OutputFormat
	var fieldsRaw string

	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.StringVar(&fieldsRaw, "fields", "value", "")
	cmdFlags.StringVar(&format.Type, "format", "text", "")
	cmdFlags.StringVar(&format.Delimiter, "delimiter", " ", "")
	cmdFlags.BoolVar(&format.Header, "header", false, "")
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }

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

	consul, err := NewConsulClient(c.Consul, &c.UI)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.KV()

	kv, _, err := client.Get(path, nil)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if kv == nil {
		return 0
	}

	kvo := NewKVOutput(c.UI, fieldsRaw)

	kvo.Output(kv, format)

	return 0
}

func (c *ReadCommand) Synopsis() string {
	return "Read a value"
}
