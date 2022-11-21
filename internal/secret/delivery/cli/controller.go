package cli

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/kmx0/GophKeeper/internal/auth/delivery/cli"
	"github.com/kmx0/GophKeeper/internal/secret"
	"github.com/kmx0/GophKeeper/internal/secret/types"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	useCase        secret.UseCase
	authMiddleWare cli.AuthMiddleware
}

func NewController(useCase secret.UseCase, authMiddleWare cli.AuthMiddleware) *Controller {
	return &Controller{
		useCase:        useCase,
		authMiddleWare: authMiddleWare,
	}
}

// var login string
// var password string
// var secretKey string
// var secretValue string

// // creditCmd represents the credit command
// var signUpCmd = &cobra.Command{
// 	Use:   "sign-up",
// 	Short: "Sign Up to gophkeeper",
// 	Long: `
//   This command creates a credit transaction for a particular user.
//   Usage: gophkeeper sign-up --login=<login> --password=<password>.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// if len(args) < 1 {
// 		// 	log.Fatal("Username not specified")
// 		// }
// 		// if len(args) < 2 {
// 		// 	log.Fatal("Password not specified")
// 		// }
// 		// fmt.Println(cmd.Flags().Args())
// 		fmt.Println(login, password)
// 		fmt.Println("Registered succesfully")
// 	},
// }

// // creditCmd represents the credit command
// var signInCmd = &cobra.Command{
// 	Use:   "sign-in",
// 	Short: "Sign In to gophkeeper",
// 	Long: `
//   This command creates a credit transaction for a particular user.
//   Usage: gophkeeper sign-up --login=<login> --password=<password>.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// if len(args) < 1 {
// 		// 	log.Fatal("Username not specified")
// 		// }
// 		// if len(args) < 2 {
// 		// 	log.Fatal("Password not specified")
// 		// }
// 		// fmt.Println(cmd.Flags().Args())
// 		fmt.Println(login, password)
// 		fmt.Println("Logged succesfully")
// 	},
// }

// func Init(rootCmd *cobra.Command) {
// 	rootCmd.AddCommand(signUpCmd)
// 	signUpCmd.Flags().StringVarP(&login, "login", "l", "", "Login for user")
// 	signUpCmd.Flags().StringVarP(&password, "password", "p", "", "Password for user")
// 	signUpCmd.MarkFlagRequired("login")
// 	signUpCmd.MarkFlagRequired("password")

// 	rootCmd.AddCommand(signInCmd)
// 	signInCmd.Flags().StringVarP(&login, "login", "l", "", "Login for user")
// 	signInCmd.Flags().StringVarP(&password, "password", "p", "", "Password for user")
// 	signInCmd.MarkFlagRequired("login")
// 	signInCmd.MarkFlagRequired("password")
// }

func (c *Controller) Create(ctx context.Context, key, value, secretType string) {

	// check user bearer

	// user := c.MustGet(auth.CtxUserKey).(*models.User)
	// get uset from bearer
	user, err := c.authMiddleWare.Handle(ctx)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	var byteValue []byte
	logrus.Info(secretType)
	if secretType == types.File {
		byteValue, err = ioutil.ReadFile(value)
		// defer the closing of our jsonFile so that we can parse it later on
		if err != nil {
			fmt.Printf("err: %s", err)
			return
		}
	} else {
		byteValue = []byte(value)
	}
	value = B64Encode(byteValue)
	if err := c.useCase.CreateSecret(ctx, user, key, value, secretType); err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Print("Secret succesfully created")
}

func (c *Controller) Get(ctx context.Context, key string) {

	user, err := c.authMiddleWare.Handle(ctx)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	sc, err := c.useCase.GetSecret(ctx, user, key)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	valueByte, err := B64Decode(sc.Value)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	logrus.Info(sc.Key)
	logrus.Info(sc.Type)
	sc.Value = string(valueByte)
	PrintSecret(sc)
}
func (c *Controller) List(ctx context.Context) {

	user, err := c.authMiddleWare.Handle(ctx)
	if err != nil {
		logrus.Error(err)
		fmt.Printf("err: %s", err)
		return
	}
	scs, err := c.useCase.GetSecrets(ctx, user)
	if err != nil {
		logrus.Error(err)
		fmt.Printf("err: %s", err)
		return
	}
	logrus.Infof("%+v", scs[0])
	for _, sc := range scs {

		valueByte, err := B64Decode(sc.Value)
		if err != nil {
			fmt.Printf("err: %s", err)
			return
		}
		sc.Value = string(valueByte)
		logrus.Info(sc.Type)
		PrintSecret(sc)
	}

}

func (c *Controller) Delete(ctx context.Context, key string) {

	user, err := c.authMiddleWare.Handle(ctx)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	err = c.useCase.DeleteSecret(ctx, user, key)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Printf("Secret %s successfully deleted", key)
}
