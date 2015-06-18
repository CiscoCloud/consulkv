package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

type ReadCommand struct {
	UI	cli.Ui
	Consul	*ConsulFlags
}

type KVOutput struct {
	Key		string		`json:",omitempty"`
	CreateIndex	uint64		`json:",omitempty"`
	ModifyIndex	uint64		`json:",omitempty"`
	LockIndex	uint64		`json:",omitempty"`
	Flags		uint64		`json:",omitempty"`
	Value		string		`json:",omitempty"`
	Session		string		`json:",omitempty"`
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
  --format=raw			Output format. Supported options: raw, json, prettyjson
				(default: raw)
`

	return strings.TrimSpace(helpText)
}

func (c *ReadCommand) Run(args []string) int {
	var fieldsRaw string
	var outputFormat string

	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.StringVar(&fieldsRaw, "fields", "value", "")
	cmdFlags.StringVar(&outputFormat, "format", "raw", "")
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

	// Copy the fields that are to be output
	//
	var output KVOutput
	for _,field := range strings.Split(fieldsRaw, ",") {
		f := strings.ToLower(field)

		switch {
		case f == "key":
			output.Key = kv.Key
		case f == "createindex":
			output.CreateIndex = kv.CreateIndex
		case f == "modifyindex":
			output.ModifyIndex = kv.ModifyIndex
		case f == "lockindex":
			output.LockIndex = kv.LockIndex
		case f == "flags":
			output.Flags = kv.Flags
		case f == "value":
			output.Value = string(kv.Value)
		case f == "session":
			output.Session = kv.Session
		default:
			c.UI.Warn(fmt.Sprintf("Ignoring invalid field: %s", field))
		}
	}

	o := strings.ToLower(outputFormat)
	switch {
	case o == "json":
		jsonRaw, err := json.Marshal(output)
		if err != nil {
			c.UI.Error("Error marshalling output")
			return 1
		}
		c.UI.Output(string(jsonRaw))
	case o == "prettyjson":
		jsonRaw, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			c.UI.Error("Error marshalling output")
			return 1
		}
		c.UI.Output(string(jsonRaw))
	case o == "raw":
		for _,field := range strings.Split(fieldsRaw, ",") {
			f := strings.ToLower(field)

			switch {
			case f == "key":
				c.UI.Output(output.Key)
			case f == "createindex":
				c.UI.Output(fmt.Sprintf("%d", output.CreateIndex))
			case f == "modifyindex":
				c.UI.Output(fmt.Sprintf("%d", output.ModifyIndex))
			case f == "lockindex":
				c.UI.Output(fmt.Sprintf("%d", output.LockIndex))
			case f == "flags":
				c.UI.Output(fmt.Sprintf("%d", output.Flags))
			case f == "value":
				c.UI.Output(output.Value)
			case f == "session":
				c.UI.Output(output.Session)
			}
		}
	default:
		c.UI.Error(fmt.Sprintf("Invalid output format: '%s", outputFormat))
		return 1
	}

	return 0
}

func (c *ReadCommand) Synopsis() string {
	return "Read a value"
}
