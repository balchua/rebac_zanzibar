package main

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func dialOpts(token string) []grpc.DialOption {
	opts := []grpc.DialOption{}

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpcutil.WithInsecureBearerToken(token))

	return opts
}

func main() {
	client, err := authzed.NewClient(
		"localhost:50051", dialOpts("supersecretthingy")...)

	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	ctx := context.Background()

	checkPortfolioPermission(ctx, client, "topdawg", "shell", "create")
	checkPortfolioPermission(ctx, client, "topdawg", "sgx", "create")
	checkPortfolioPermission(ctx, client, "madame_oracle", "sgx", "create")
	checkPortfolioPermission(ctx, client, "minime", "sgx", "create")
	checkPortfolioPermission(ctx, client, "minime", "sgx", "read")
	checkPortfolioPermission(ctx, client, "minime", "shell", "read")
	checkDocumentPermission(ctx, client, "minime", "findoc", "read")
	checkDocumentPermission(ctx, client, "minime", "findoc", "update")

	// resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_NO_PERMISSION

	readRelationship(ctx, client, "industry")
	start := time.Now()

	for i := 0; i < 2; i++ {
		checkPortfolioPermission(ctx, client, "topdawg", "shell", "create")
		// checkPortfolioPermission(ctx, client, "topdawg", "sgx", "create")
		// checkPortfolioPermission(ctx, client, "madame_oracle", "sgx", "create")
		// checkPortfolioPermission(ctx, client, "minime", "sgx", "create")
		// checkPortfolioPermission(ctx, client, "minime", "sgx", "read")
		// checkPortfolioPermission(ctx, client, "minime", "shell", "read")
		// checkDocumentPermission(ctx, client, "minime", "findoc", "read")
		// checkDocumentPermission(ctx, client, "minime", "findoc", "update")
	}
	end := time.Since(start)
	log.Printf("Elapsed time : %d", end.Milliseconds())
}

func checkPortfolioPermission(ctx context.Context, client *authzed.Client, userId string, portfolioId string, action string) {
	permission := checkUserPermission(ctx, client, "portfolio", userId, portfolioId, action)
	log.Printf("%s permission is %t for user %s on portfolio %s", action, permission, userId, portfolioId)
}

func checkDocumentPermission(ctx context.Context, client *authzed.Client, userId string, documentId string, action string) {
	permission := checkUserPermission(ctx, client, "document", userId, documentId, action)
	log.Printf("%s permission is %t for user %s on document %s", action, permission, userId, documentId)
}

func checkUserPermission(ctx context.Context, client *authzed.Client, resourceType string, userId string, portfolioId string, action string) bool {
	user := &pb.SubjectReference{Object: &pb.ObjectReference{
		ObjectType: "user",
		ObjectId:   userId,
	}}

	portfolio := &pb.ObjectReference{
		ObjectType: resourceType,
		ObjectId:   portfolioId,
	}

	resp, err := client.CheckPermission(ctx, &pb.CheckPermissionRequest{
		Resource:   portfolio,
		Permission: action,
		Subject:    user,
	})
	if err != nil {
		log.Fatalf("failed to check permission: %s", err)
	}

	return resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION
}

func readRelationship(ctx context.Context, client *authzed.Client, object string) {
	readFilter := &v1.RelationshipFilter{ResourceType: object}
	readRelationshipReq := &v1.ReadRelationshipsRequest{
		RelationshipFilter: readFilter,
	}

	resp, err := client.ReadRelationships(context.Background(), readRelationshipReq)
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
