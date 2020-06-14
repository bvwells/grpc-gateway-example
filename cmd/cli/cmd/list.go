package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/bvwells/grpc-gateway-example/proto/beers"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const address = "localhost:50000"

var page int32 = 1

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "list",
	Short:   "list lists beers",
	Long:    `list lists beers`,
	Example: `cli list beers`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 || args[0] != "beers" {
			fmt.Println("invalid command line arguments")
			os.Exit(1)
		}

		// Set up a connection to the server.
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			fmt.Println("unable to connect to beer grpc server")
			os.Exit(1)
		}
		defer conn.Close()

		c := beers.NewBeerServiceClient(conn)

		ctx := context.Background()
		resp, err := c.ListBeers(ctx, &beers.ListBeersRequest{
			Page: page,
		})
		if err != nil {
			return
		}

		for _, beer := range resp.Beers {
			fmt.Println(beer)
		}
	},
}

func init() {
	getCmd.Flags().Int32VarP(&page, "page", "p", 1, "page number")

	rootCmd.AddCommand(getCmd)
}
