package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret/types"
)

func PrintSecret(secret *models.Secret) {
	if secret.Type != types.File {
		fmt.Printf("%+v", secret)
		return
	}
	err := ioutil.WriteFile(fmt.Sprintf("/tmp/%s_%s", secret.UserID, secret.Key), []byte(secret.Value), 0777)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	fmt.Printf("Secret in file /tmp/%s_%s\n", secret.UserID, secret.Key)
}
