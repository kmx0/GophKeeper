package cli

import (
	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/spf13/cobra"
)

var login string
var password string


func RegisterAuthCmdEndpoints(rootCmd *cobra.Command, uc auth.UseCase) {
	c := NewController(uc)
	signInCmd := &cobra.Command{
		Use:   "sign-in",
		Short: "Sign In to gophkeeper",
		Long: `
	  This command creates a credit transaction for a particular user.
	  Usage: gophkeeper sign-up --login=<login> --password=<password>.`,
		Run: func(cmd *cobra.Command, args []string) {
			c.SignIn(rootCmd.Context(), login, password)
		},
	}

	// creditCmd represents the credit command
	signUpCmd := &cobra.Command{
		Use:   "sign-up",
		Short: "Sign Up to gophkeeper",
		Long: `
  This command creates a credit transaction for a particular user.
  Usage: gophkeeper sign-up --login=<login> --password=<password>.`,
		Run: func(cmd *cobra.Command, args []string) {
			c.SignUp(cmd.Context(), login, password)
		},
	}

	rootCmd.AddCommand(signUpCmd)
	signUpCmd.Flags().StringVarP(&login, "login", "l", "", "Login for user")
	signUpCmd.Flags().StringVarP(&password, "password", "p", "", "Password for user")
	signUpCmd.MarkFlagRequired("login")
	signUpCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(signInCmd)
	signInCmd.Flags().StringVarP(&login, "login", "l", "", "Login for user")
	signInCmd.Flags().StringVarP(&password, "password", "p", "", "Password for user")
	signInCmd.MarkFlagRequired("login")
	signInCmd.MarkFlagRequired("password")
}
