package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/kmx0/GophKeeper/internal/auth/delivery/cli"
	"github.com/kmx0/GophKeeper/internal/secret"
	"github.com/kmx0/GophKeeper/internal/secret/types"
)

type Controller struct {
	useCase    secret.UseCase
	authStatus cli.AuthStatus
	writer     io.Writer
}

func NewController(useCase secret.UseCase, authStatus *cli.AuthStatus, writer io.Writer) *Controller {
	return &Controller{
		useCase:    useCase,
		authStatus: *authStatus,
		writer:     writer,
	}
}

func (c *Controller) Create(ctx context.Context, key, value, secretType string) {

	// check user bearer

	// user := c.MustGet(auth.CtxUserKey).(*models.User)
	// get uset from bearer
	user, err := c.authStatus.CheckAuthStatus(ctx)
	if err != nil {
		fmt.Fprintf(c.writer, "err: %s", err)
		return
	}
	var byteValue []byte
	if secretType == types.File {
		byteValue, err = ioutil.ReadFile(value)
		if err != nil {
			fmt.Fprintf(c.writer, "err: %s", err)
			return
		}
	} else {
		byteValue = []byte(value)
	}
	value = B64Encode(byteValue)
	if err := c.useCase.CreateSecret(ctx, user, key, value, secretType); err != nil {
		fmt.Fprintf(c.writer, "err: %s", err)
		return
	}
	fmt.Fprintf(c.writer, "Secret succesfully created")

}

func (c *Controller) Get(ctx context.Context, key, saveFile string) {

	user, err := c.authStatus.CheckAuthStatus(ctx)
	if err != nil {
		fmt.Fprintf(c.writer, "%s", err)
		return
	}
	sc, err := c.useCase.GetSecret(ctx, user, key)
	if err != nil {
		fmt.Fprintf(c.writer, "%s", err)
		return
	}
	valueByte, err := B64Decode(sc.Value)
	if err != nil {
		fmt.Fprintf(c.writer, "%s", err)
		return
	}
	sc.Value = string(valueByte)
	PrintSecret(sc, c.writer, saveFile)
}
func (c *Controller) List(ctx context.Context, saveFile string) {

	user, err := c.authStatus.CheckAuthStatus(ctx)
	if err != nil {
		fmt.Fprintf(c.writer, "err: %s", err)
		return
	}
	scs, err := c.useCase.GetSecrets(ctx, user)
	if err != nil {
		fmt.Fprintf(c.writer, "err: %s", err)
		return
	}
	for _, sc := range scs {

		valueByte, err := B64Decode(sc.Value)
		if err != nil {
			fmt.Fprintf(c.writer, "err: %s", err)
			return
		}
		sc.Value = string(valueByte)
		PrintSecret(sc, c.writer, saveFile)
	}

}

func (c *Controller) Delete(ctx context.Context, key string) {

	user, err := c.authStatus.CheckAuthStatus(ctx)
	if err != nil {
		fmt.Fprintf(c.writer, "err: %s", err)
		return
	}
	err = c.useCase.DeleteSecret(ctx, user, key)
	if err != nil {
		fmt.Fprintf(c.writer, "err: %s", err)
		return
	}
	fmt.Fprintf(c.writer, "Secret %s successfully deleted", key)
}
