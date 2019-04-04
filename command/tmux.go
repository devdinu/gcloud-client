package command

import (
	"fmt"
	"strings"
)

type tmux struct {
	hosts []string
	cmd   string
	TmuxConfig
}

func (cmd *tmux) Name() string { return "tmuxinator" }

func (cmd *tmux) Args() []string {
	return strings.Fields(cmd.startCommand())
}

func (cmd *tmux) startCommand() string {
	return fmt.Sprintf(`start %s cmd="%s" hosts=%s session_name=%s %s`, cmd.Project, cmd.cmd, strings.Join(cmd.hosts, ","), cmd.TmuxConfig.SessionName(), strings.Join(cmd.TmuxConfig.Flags(), " "))
}

func (cmd *tmux) String() string {
	return fmt.Sprintf(`%s %s`, cmd.Name(), cmd.startCommand())
}

type TmuxConfig struct {
	Session string
	user    string
	cmd     string
	Project string
	keyVals map[string]string
}

func (c *TmuxConfig) AddArg(key, val string) {
	if c.keyVals == nil {
		c.keyVals = make(map[string]string, 1)
	}
	c.keyVals[key] = val
}

func (c *TmuxConfig) SessionName() string {
	if c.Session != "" {
		return c.Session
	}
	return "default-sesion"
}

func (c *TmuxConfig) Flags() []string {
	var keyargs []string
	for k, v := range c.keyVals {
		keyargs = append(keyargs, fmt.Sprintf("%s=%s", k, v))
	}
	return keyargs
}
