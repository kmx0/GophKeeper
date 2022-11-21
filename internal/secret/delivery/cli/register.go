package cli

import (
	"github.com/kmx0/GophKeeper/internal/auth/delivery/cli"
	"github.com/kmx0/GophKeeper/internal/secret"
	"github.com/kmx0/GophKeeper/internal/secret/types"
	"github.com/spf13/cobra"
)

var secretKey string
var secretValue string
var secretIsFile bool

// authEndpoints.POST("/create",  h.Create)
// authEndpoints.POST("/get",     h.Get)
// authEndpoints.POST("/list", h.List)
// authEndpoints.POST("/delete",  h.Delete)

func Init(rootCmd *cobra.Command) {

}

func RegisterSecretCmdEndpoints(rootCmd *cobra.Command, uc secret.UseCase, authMiddleware *cli.AuthMiddleware) {
	c := NewController(uc, *authMiddleware)
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create new secret",
		Long: `
	  This command creates a credit transaction for a particular user.
	  Usage: gophkeeper create --key=<key> --value=<value|path_to_secret> --file=<true|false> .`,
		Run: func(cmd *cobra.Command, args []string) {
			if secretIsFile {
				c.Create(cmd.Context(), secretKey, secretValue, types.File)
				return
			}
			c.Create(cmd.Context(), secretKey, secretValue, types.String)

		},
	}

	// getCmd represents the credit command
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get secret by key",
		Long: `
  This command get a credit transaction for a particular user.
  Usage: gophkeeper get --key=<key>.`,
		Run: func(cmd *cobra.Command, args []string) {
			c.Get(cmd.Context(), secretKey)
		},
	}

	// listCmd represents the credit command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all secrets",
		Long: `
  This command creates a credit transaction for a particular user.
  Usage: gophkeeper create --key=<key> --value=<value|path_to_secret>.`,
		Run: func(cmd *cobra.Command, args []string) {
			c.List(cmd.Context())

		},
	}

	// listCmd represents the credit command
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete secret by key",
		Long: `
  This command creates a credit transaction for a particular user.
  Usage: gophkeeper create --key=<key> --value=<value|path_to_secret>.`,
		Run: func(cmd *cobra.Command, args []string) {
			// c.Delete()
			c.Delete(cmd.Context(), secretKey)
		},
	}

	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&secretKey, "key", "k", "", "Key for new secret")
	createCmd.Flags().StringVarP(&secretValue, "value", "v", "", "Value of secret")
	createCmd.Flags().BoolVarP(&secretIsFile, "file", "f", false, "secret is file")
	createCmd.MarkFlagRequired("key")
	createCmd.MarkFlagRequired("value")

	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&secretKey, "key", "k", "", "Key for secret")
	getCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(listCmd)

	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&secretKey, "key", "k", "", "Key for secret")
	deleteCmd.MarkFlagRequired("key")
}
