package deal

import "context"

type DealService interface {
	CreateDeal(ctx context.Context, dealInfo DealInfo) (DealInfo, error)
	GetAll(ctx context.Context) ([]DealInfo, error)
	Get(ctx context.Context, dealId string) (DealInfo, error)
	Save(context.Context, DealInfo) (DealInfo, error)
	Delete(ctx context.Context, dealId string) (int64, error)
}
