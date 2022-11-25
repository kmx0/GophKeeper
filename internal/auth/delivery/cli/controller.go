package cli

import (
	"context"
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

func (c *Controller) SignUp(ctx context.Context, login, password string) {
	// inp := new(signInput)
	fmt.Println(login)
	fmt.Println(password)

	if len(login) == 0 {
		fmt.Fprintf(c.writer, "login is empty!")
		// fmt.Println("login is empty!")
		return
	}
	if len(password) == 0 {
		fmt.Println("password is empty!")
		return
	}

	if err := c.useCase.SignUp(ctx, login, password); err != nil {
		fmt.Printf("internal server error: %s", err.Error())

		return
	}
	fmt.Println("Successfully registered")

}

func (c *Controller) SignIn(ctx context.Context, login, password string) {
	fmt.Println(login)
	fmt.Println(password)

	if len(login) == 0 {
		fmt.Println("login is empty!")
		return
	}
	if len(password) == 0 {
		fmt.Println("password is empty!")
		return
	}

	token, err := c.useCase.SignIn(ctx, login, password)
	if err != nil {
		if err == auth.ErrUserNotFound {
			fmt.Printf("not such login or password")
			return
		}

		fmt.Printf("internal server error: %s", err.Error())
		return
	}

	fmt.Println("successfully logged in")
	fmt.Println(token)

	// c.JSON(http.StatusOK, signInResponse{Token: token})

}
