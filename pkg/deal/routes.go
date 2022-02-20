package deal

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type DealRoutes struct {
	dealService DealService
}

func NewDealRoutes(dealService DealService) *DealRoutes {
	return &DealRoutes{
		dealService: dealService,
	}
}

func (d *DealRoutes) GetAllDeals(c echo.Context) error {
	ctx := c.Request().Context()
	zap.S().Info("retrieving all characters...")
	characters, err := d.dealService.GetAll(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, characters)
}

func (d *DealRoutes) GetDeal(c echo.Context) error {
	ctx := c.Request().Context()
	dealId := c.Param("id")
	userId := c.Request().Header.Get("user")
	zap.S().Infof("retrieving a deal %s", dealId)
	dealCtx := context.WithValue(ctx, "user", userId)
	deal, err := d.dealService.Get(dealCtx, dealId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, deal)
}

func (d *DealRoutes) CreateDeal(c echo.Context) (err error) {
	ctx := c.Request().Context()
	userId := c.Request().Header.Get("user")
	di := new(DealInfo)
	if err = c.Bind(di); err != nil {
		return
	}
	zap.S().Info("creating a deal")
	dealCtx := context.WithValue(ctx, "user", userId)
	deal, err := d.dealService.CreateDeal(dealCtx, *di)

	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err)
	}

	return c.JSON(http.StatusOK, deal)
}
