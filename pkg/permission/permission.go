package permission

import (
	"context"
	"errors"
	"io"
	"log"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"

	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DealPermissionServiceImpl struct {
	client *authzed.Client
}

func NewDealPermissionService(token string, host string) (*DealPermissionServiceImpl, error) {
	opts := []grpc.DialOption{}

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpcutil.WithInsecureBearerToken(token))
	client, err := authzed.NewClient(
		host, opts...)
	if err != nil {
		zap.S().Errorf("unable to initialize client: %v", err)
		return nil, err
	}

	return &DealPermissionServiceImpl{
		client: client,
	}, nil

}

func (d *DealPermissionServiceImpl) CanView(ctx context.Context, dealId string, userId string) bool {
	resp, err := d.checkPermission(ctx, "deal", dealId, userId, "read_core_section")

	if err != nil {
		return false
	}

	return resp
}

func (d *DealPermissionServiceImpl) CanViewSupplementaryInfo(ctx context.Context, dealId string, userId string) bool {
	resp, err := d.checkPermission(ctx, "deal", dealId, userId, "read_supplementary_section")

	if err != nil {
		return false
	}

	return resp
}

func (d *DealPermissionServiceImpl) CanViewServicingInfo(ctx context.Context, dealId string, userId string) bool {
	resp, err := d.checkPermission(ctx, "deal", dealId, userId, "read_servicing_section")

	if err != nil {
		return false
	}

	return resp
}

func (d *DealPermissionServiceImpl) CanCreateDeal(ctx context.Context, groupId string, userId string) bool {
	resp, err := d.checkPermission(ctx, "group", groupId, userId, "can_create_deal")

	if err != nil {
		return false
	}

	return resp
}

func (d *DealPermissionServiceImpl) CanUpdateCoreInformation(ctx context.Context, dealId string, userId string) bool {
	resp, err := d.checkPermission(ctx, "deal", dealId, userId, "update_core_section")

	if err != nil {
		return false
	}

	return resp
}

func (d *DealPermissionServiceImpl) checkPermission(ctx context.Context, resourceType string, dealId string, userId string, permissionAction string) (bool, error) {
	user := &pb.SubjectReference{Object: &pb.ObjectReference{
		ObjectType: "user",
		ObjectId:   userId,
	}}

	deal := &pb.ObjectReference{
		ObjectType: resourceType,
		ObjectId:   dealId,
	}

	resp, err := d.client.CheckPermission(ctx, &pb.CheckPermissionRequest{
		Resource:   deal,
		Permission: permissionAction,
		Subject:    user,
	})

	if err != nil {
		zap.S().Errorf("unable to get permission %v", err)
		return false, err
	}
	return resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION, nil
}

func (d *DealPermissionServiceImpl) WriteDealRelationship(ctx context.Context, dealId string, subjectId string, subjectRel string) error {

	request := &pb.WriteRelationshipsRequest{
		Updates: []*pb.RelationshipUpdate{
			{
				Operation: pb.RelationshipUpdate_OPERATION_CREATE,
				Relationship: &pb.Relationship{
					Resource: &pb.ObjectReference{
						ObjectType: "deal",
						ObjectId:   dealId,
					},
					Relation: "group",
					Subject: &pb.SubjectReference{
						Object: &pb.ObjectReference{
							ObjectType: "group",
							ObjectId:   subjectId,
						},
						OptionalRelation: subjectRel,
					},
				},
			},
		},
		OptionalPreconditions: nil,
	}

	response, err := d.client.WriteRelationships(ctx, request)
	if err != nil {
		return err
	}

	zap.S().Infof("Token : %s", response.WrittenAt.Token)
	return nil
}

func (d *DealPermissionServiceImpl) readRelationship(ctx context.Context, object string) {
	readFilter := &v1.RelationshipFilter{ResourceType: object}
	readRelationshipReq := &v1.ReadRelationshipsRequest{
		RelationshipFilter: readFilter,
	}

	resp, err := d.client.ReadRelationships(context.Background(), readRelationshipReq)
	if err == nil {
		for {
			msg, recErr := resp.Recv()
			if recErr != nil && !errors.Is(recErr, io.EOF) {
				log.Printf("%v", recErr)
			} else if recErr == nil {
				log.Printf("relationship %s", msg.String())
			} else {
				log.Print("no more relationships")
				return
			}
		}

	} else {
		log.Printf("%v", err)
	}

}
