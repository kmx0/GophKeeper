package cli

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret/types"
)

func PrintSecret(secret *models.Secret, writer io.Writer, saveFile string) {
	if secret.Type != types.File {
		fmt.Fprintf(writer, "%s:%s", secret.Key, secret.Value)

		return
	}
	err := ioutil.WriteFile(fmt.Sprintf("%s_%s", saveFile, secret.Key), []byte(secret.Value), 0777)
	if err != nil {
		fmt.Fprintf(writer, "err: %s", err)

	}
	fmt.Fprintf(writer, "Secret in file %s_%s\n", saveFile, secret.Key)
}
