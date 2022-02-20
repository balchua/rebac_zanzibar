package deal

import (
	"context"
	"errors"
	"time"

	"github.com/balchua/demo-spicedb/pkg/permission"
	"go.uber.org/zap"
)

type DealServiceImpl struct {
	permissionService permission.DealPermissionService
}

func NewDealService(permissionService permission.DealPermissionService) (*DealServiceImpl, error) {

	return &DealServiceImpl{
		permissionService: permissionService,
	}, nil
}
func (d *DealServiceImpl) GetAll(ctx context.Context) ([]DealInfo, error) {
	dealInfos := make([]DealInfo, 1)
	return dealInfos, nil
}

func (d *DealServiceImpl) Get(ctx context.Context, dealId string) (DealInfo, error) {

	var dealInfo DealInfo

	userId := ctx.Value("user")
	coreView := d.permissionService.CanView(ctx, dealId, userId.(string))
	zap.S().Infof("Can user %s view deal %s ==> %t", userId.(string), dealId, coreView)
	if coreView == true {
		dealInfo = DealInfo{
			DealId:     dealId,
			ClientName: "ClientA",
			ClientId:   "12345",
			ContractId: "ABCDEFG",
			Date:       time.Now(),
		}
		supplResp := d.permissionService.CanViewSupplementaryInfo(ctx, dealId, userId.(string))
		zap.S().Infof("Can user %s view supplementary info %s ==> %t", userId.(string), dealId, supplResp)
		if supplResp == true {

			supplInfo := DealSupplementaryInfo{
				SupplInfo1: "1111",
				SupplInfo2: "2222",
				SupplInfo3: "3333",
			}
			dealInfo.Supplementary = supplInfo
		}
		servicingResp := d.permissionService.CanViewServicingInfo(ctx, dealId, userId.(string))
		zap.S().Infof("Can user %s view servicing info %s ==> %t", userId.(string), dealId, servicingResp)
		if servicingResp == true {
			servInfo := DealServicingInfo{
				Servicing1: "Srv1",
				Servicing2: "Srv2",
				Servicing3: "Srv3",
			}
			dealInfo.Servicing = servInfo
		}
	}

	return dealInfo, nil
}

func (d *DealServiceImpl) CreateDeal(ctx context.Context, dealInfo DealInfo) (DealInfo, error) {
	userId := ctx.Value("user")

	resp := d.permissionService.CanCreateDeal(ctx, dealInfo.Origin, userId.(string))

	if resp == true {

		err := d.permissionService.WriteDealRelationship(ctx, dealInfo.DealId, dealInfo.Origin, "")
		if err != nil {
			return DealInfo{}, err
		}

		return DealInfo{
			DealId:     "99",
			ClientName: "ClientA",
			ClientId:   "12345",
			ContractId: "ABCDEFG",
			Date:       time.Now(),
			Origin:     dealInfo.Origin,
		}, nil
	}
	return DealInfo{}, errors.New("no permission to create deal")
}

func (d *DealServiceImpl) Save(context.Context, DealInfo) (DealInfo, error) {
	return DealInfo{}, nil
}

func (d *DealServiceImpl) Delete(ctx context.Context, dealId string) (int64, error) {
	return 0, nil
}
