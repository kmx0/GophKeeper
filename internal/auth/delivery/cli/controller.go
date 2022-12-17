package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/kmx0/GophKeeper/internal/auth"
)

type Controller struct {
	writer  io.Writer
	useCase auth.UseCase
}

func NewController(writer io.Writer, useCase auth.UseCase) *Controller {
	return &Controller{
		writer:  writer,
		useCase: useCase,
	}
}

func (c *Controller) SignUp(ctx context.Context, login, password string) {

	if len(login) == 0 {
		fmt.Fprintf(c.writer, "login is empty!")
		return
	}
	if len(password) == 0 {
		fmt.Fprintf(c.writer, "password is empty!")
		return
	}

	if err := c.useCase.SignUp(ctx, login, password); err != nil {
		fmt.Fprintf(c.writer, "internal server error: %s", err.Error())
		return
	}
	fmt.Fprintf(c.writer, "Successfully registered")
}

func (c *Controller) SignIn(ctx context.Context, login, password string) {

	if len(login) == 0 {
		fmt.Fprintf(c.writer, "login is empty!")
		return
	}
	if len(password) == 0 {
		fmt.Fprintf(c.writer, "password is empty!")
		return
	}

	_, err := c.useCase.SignIn(ctx, login, password)
	if err != nil {
		if errors.Unwrap(err) == auth.ErrUserNotFound {
			fmt.Fprintf(c.writer, "%s", err.Error())
			return
		}
		if errors.Unwrap(err) == auth.ErrIncorrectPassword {
			fmt.Fprintf(c.writer, "%s", err.Error())
			return
		}

		fmt.Fprintf(c.writer, "internal server error: %s", err.Error())
		return
	}

	fmt.Fprintf(c.writer, "Successfully logged in")
}
