package shell

import (
	"container/ring"
)

type CmdHistory struct {
	max    int
	cmds   *ring.Ring
	offset int
}

func InitCmdHistory(max int) *CmdHistory {
	return &CmdHistory{
		max:  max,
		cmds: ring.New(0),
	}
}

func (c *CmdHistory) All() []string {
	cmds := []string{}
	c.reset()
	for i := 0; i < c.cmds.Len(); i++ {
		cmds = append(cmds, c.cmds.Value.(string))
		c.cmds = c.cmds.Next()
	}
	return cmds
}

func (c *CmdHistory) Next() string {
	if c.cmds.Len() == 0 {
		return ""
	}
	c.offset = (c.offset + 1) % c.cmds.Len()
	c.cmds = c.cmds.Next()
	return c.cmds.Value.(string)
}

func (c *CmdHistory) Prev() string {
	if c.cmds.Len() == 0 {
		return ""
	}
	c.offset -= 1
	if c.offset < 0 {
		c.offset = c.cmds.Len()
	}
	c.cmds = c.cmds.Prev()
	return c.cmds.Value.(string)
}

func (c *CmdHistory) Add(cmd string) {
	if c.cmds.Len() == 0 {
		c.cmds = &ring.Ring{Value: cmd}
		return
	}

	c.reset()
	if c.cmds.Len() >= c.max {
		// remove the oldest cmd in the buffer
		c.cmds = c.cmds.Move(c.cmds.Len() - 1)
		c.cmds.Unlink(1)
		c.cmds = c.cmds.Next()
	}

	c.cmds = c.cmds.Move(-1)
	c.cmds = c.cmds.Link(&ring.Ring{Value: cmd})
}

func (c *CmdHistory) reset() {
	if c.cmds.Len() == 0 {
		return
	}
	c.cmds = c.cmds.Move(-c.offset)
}
