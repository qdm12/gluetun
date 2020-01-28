package command

func (c *command) Unbound() error {
	_, err := c.command.Run("unbound")
	return err
}
