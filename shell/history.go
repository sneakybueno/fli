package shell

import (
	"container/ring"
)

// CmdHistory keeps track of cmds using a ring buffer
type CmdHistory struct {
	max    int
	cmds   *ring.Ring
	offset int
}

// InitCmdHistory inits the ring buffer and sets the history's max capacity
func InitCmdHistory(max int) *CmdHistory {
	return &CmdHistory{
		max:  max,
		cmds: ring.New(0),
	}
}

// All returns a list of all previously added commands
func (c *CmdHistory) All() []string {
	cmds := []string{}
	c.reset()
	for i := 0; i < c.cmds.Len(); i++ {
		cmds = append(cmds, c.cmds.Value.(string))
		c.cmds = c.cmds.Next()
	}
	return cmds
}

// Next iterates "clockwise" around the ring from the current position
func (c *CmdHistory) Next() string {
	if c.cmds.Len() == 0 {
		return ""
	}
	c.offset = (c.offset + 1) % c.cmds.Len()
	c.cmds = c.cmds.Next()
	return c.cmds.Value.(string)
}

// Prev iterates "counter-clockwise" around the ring from the current position
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

// Add inserts a new cmd into the ring. If the ring is not empty, the new cmd
// is the predecessor of the last cmd that was most recently inserted. If the
// ring's capacity is reached, the oldest cmd is bumped.
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
