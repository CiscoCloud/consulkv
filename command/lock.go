package command

import (
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type LockCommand struct {
	UI	cli.Ui
	Consul	*ConsulFlags
}

func (c *LockCommand) Help() string {
	helpText := `
Usage: consulkv lock [options] path

  Acquire a lock on a given path

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
  --behavior=release		Lock behavior. One of 'release' or 'delete'
				(default: release)
  --ttl=15s			Lock time to live
				(default: 15s)
  --lock-delay=5s		Lock delay
				(default: 5s)
`

	return strings.TrimSpace(helpText)
}

func (c *LockCommand) Run(args []string) int {
	var behavior string
	var ttl string
	var lockDelay time.Duration

	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.StringVar(&behavior, "behavior", "release", "")
	cmdFlags.StringVar(&ttl, "ttl", "15s", "")
	cmdFlags.DurationVar(&lockDelay, "lock-delay", 5 * time.Second, "")

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

	client, err := NewConsulClient(c.Consul, &c.UI)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	sessionClient := client.Session()
	writeOpts := new(consulapi.WriteOptions)

	// Create the Consul session
	sessionId, _, err := sessionClient.CreateNoChecks(&consulapi.SessionEntry{
						Name:		"Session for consulkv",
						LockDelay:	lockDelay,
						Behavior:	behavior,
						TTL:		ttl,
						}, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Set the session to renew periodically
	sessionRenew := make(chan struct{})
	go sessionClient.RenewPeriodic(ttl, sessionId, nil, sessionRenew)
	defer func() {
		close(sessionRenew)
		sessionRenew = nil
	}()

	// Create the Lock Structure
	lockOpts := consulapi.LockOptions{
		Key:		path,
		Session:	sessionId,
		}
	l, err := client.LockOpts(&lockOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	_, err = l.Lock(nil)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(sessionId)

	return 0
}

func (c *LockCommand) Synopsis() string {
	return "Lock a node"
}
