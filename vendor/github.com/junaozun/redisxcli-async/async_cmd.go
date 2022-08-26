package redisclix

import (
	"fmt"
)

// doCmd server命令
type doCmd struct {
	cmd      []interface{}
	ret      string
	retError error
	callback func(string, error)
}

func (c *doCmd) ExecCmd(cli IClient) error {
	ret, err := cli.Do(c.cmd...)
	c.ret = ret
	c.retError = err
	return err
}

func (c *doCmd) Key() string {
	if len(c.cmd) > 0 {
		return fmt.Sprintf("%v", c.cmd[0])
	}
	return ""
}

func (c *doCmd) Callback() {
	if c.callback != nil {
		c.callback(c.ret, c.retError)
	}
}

type hmsetCmd struct {
	key   string
	value ArrayReq
	err   error
}

func (c *hmsetCmd) ExecCmd(cli IClient) error {
	c.err = cli.HMSet(c.key, c.value)
	return c.err
}

func (c *hmsetCmd) Key() string {
	return c.key
}

type hmgetCmd struct {
	key      string
	field    []string
	ret      []string
	retError error
	callback func([]string, error)
}

func (c *hmgetCmd) ExecCmd(cli IClient) error {
	c.ret, c.retError = cli.HMGet(c.key, c.field...)
	return c.retError
}

func (c *hmgetCmd) Key() string {
	return c.key
}

func (c *hmgetCmd) Callback() {
	if c.callback != nil {
		c.callback(c.ret, c.retError)
	}
}

// delCmd delete命令
type delCmd struct {
	key      string
	err      error
	callback func(error)
}

func (c *delCmd) ExecCmd(cli IClient) error {
	c.err = cli.Del(c.key)
	return c.err
}

func (c *delCmd) Key() string {
	return c.key
}

type getCmd struct {
	key      string
	ret      string
	retErr   error
	callback func(string, error)
}

func (c *getCmd) ExecCmd(cli IClient) error {
	c.ret, c.retErr = cli.Get(c.key)
	return c.retErr
}

func (c *getCmd) Key() string {
	return c.key
}

func (c *getCmd) Callback() {
	if c.callback != nil {
		c.callback(c.ret, c.retErr)
	}
}

type setCmd struct {
	key   string
	value interface{}
	err   error
}

func (c *setCmd) ExecCmd(cli IClient) error {
	c.err = cli.Set(c.key, c.value, 0)
	return c.err
}

func (c *setCmd) Key() string {
	return c.key
}
