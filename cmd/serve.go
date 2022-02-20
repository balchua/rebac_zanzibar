/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/balchua/demo-spicedb/pkg/deal"
	"github.com/balchua/demo-spicedb/pkg/permission"
	"github.com/labstack/echo"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the application",
	Long:  `Starts serving the application`,
	Run:   setupRoutes,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func setupRoutes(cmd *cobra.Command, args []string) {
	//ctx := context.Background()
	e := echo.New()
	permissionService, err := permission.NewDealPermissionService("supersecretthingy", "localhost:50051")

	if err != nil {
		zap.S().Fatalf("unable to connect to spicedb %v", err)
	}

	dealService, err := deal.NewDealService(permissionService)

	if err != nil {
		zap.S().Fatalf("unable to create deal service %v", err)
	}

	dealRoutes := deal.NewDealRoutes(dealService)

	e.GET("/api/v1/deals/all", dealRoutes.GetAllDeals)
	e.GET("/api/v1/deal/:id", dealRoutes.GetDeal)
	e.POST("/api/v1/deal", dealRoutes.CreateDeal)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		zap.S().Error(err.Error())

		// Call the default handler to return the HTTP response
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.Logger.Fatal(e.Start(":8181"))
}
