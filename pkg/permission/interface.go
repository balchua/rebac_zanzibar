package permission

import "context"

type DealPermissionService interface {
	CanCreateDeal(ctx context.Context, groupId string, userId string) bool
	CanView(ctx context.Context, dealId string, userId string) bool
	CanViewSupplementaryInfo(ctx context.Context, dealId string, userId string) bool
	CanViewServicingInfo(ctx context.Context, dealId string, userId string) bool
	CanUpdateCoreInformation(ctx context.Context, dealId string, userId string) bool
}
