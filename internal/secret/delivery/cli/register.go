package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var secretKey string
var secretValue string

// createCmd represents the credit command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new secret",
	Long: `
  This command creates a credit transaction for a particular user.
  Usage: gophkeeper create --key=<key> --value=<value|path_to_secret>.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) < 1 {
		// 	log.Fatal("Username not specified")
		// }
		// if len(args) < 2 {
		// 	log.Fatal("Password not specified")
		// }
		// fmt.Println(cmd.Flags().Args())
		fmt.Println(secretKey, secretValue)
		fmt.Println("Logged succesfully")
	},
}

// getCmd represents the credit command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secret by key",
	Long: `
  This command creates a credit transaction for a particular user.
  Usage: gophkeeper create --key=<key> --value=<value|path_to_secret>.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) < 1 {
		// 	log.Fatal("Username not specified")
		// }
		// if len(args) < 2 {
		// 	log.Fatal("Password not specified")
		// }
		// fmt.Println(cmd.Flags().Args())
		fmt.Println(secretKey)
		fmt.Println("Logged succesfully")
	},
}

// listCmd represents the credit command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets",
	Long: `
  This command creates a credit transaction for a particular user.
  Usage: gophkeeper create --key=<key> --value=<value|path_to_secret>.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) < 1 {
		// 	log.Fatal("Username not specified")
		// }
		// if len(args) < 2 {
		// 	log.Fatal("Password not specified")
		// }
		// fmt.Println(cmd.Flags().Args())
		// fmt.Println(secretKey)
		fmt.Println("Logged succesfully")
	},
}

// listCmd represents the credit command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete secret by key",
	Long: `
  This command creates a credit transaction for a particular user.
  Usage: gophkeeper create --key=<key> --value=<value|path_to_secret>.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) < 1 {
		// 	log.Fatal("Username not specified")
		// }
		// if len(args) < 2 {
		// 	log.Fatal("Password not specified")
		// }
		// fmt.Println(cmd.Flags().Args())
		fmt.Println(secretKey)
		fmt.Println("Logged succesfully")
	},
}

// authEndpoints.POST("/create",  h.Create)
// authEndpoints.POST("/get",     h.Get)
// authEndpoints.POST("/list", h.List)
// authEndpoints.POST("/delete",  h.Delete)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&secretKey, "key", "k", "", "Key for new secret")
	createCmd.Flags().StringVarP(&secretValue, "value", "v", "", "Value of secret")
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
