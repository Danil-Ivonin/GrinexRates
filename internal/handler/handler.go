package handler

import (
	"context"
	"strings"

	"github.com/Danil-Ivonin/GrinexRates/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	pb "github.com/Danil-Ivonin/GrinexRates/gen/rates"
)

// RatesHandler implements pb.RatesServiceServer
type RatesHandler struct {
	pb.UnimplementedRatesServiceServer
	svc ports.RatesService
}

// NewRatesHandler New creates a RatesHandler
func NewRatesHandler(svc ports.RatesService) *RatesHandler {
	return &RatesHandler{svc: svc}
}

// GetRates implements the GetRates RPC
func (h *RatesHandler) GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	snap, err := h.svc.GetRates(ctx, req.GetN(), req.GetM())
	if err != nil {
		if strings.Contains(err.Error(), "calculator:") {
			return nil, status.Errorf(codes.InvalidArgument, "invalid parameters: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "get rates failed: %v", err)
	}

	return &pb.GetRatesResponse{
		AskPrice:  snap.AskPrice.String(),
		BidPrice:  snap.BidPrice.String(),
		TopN:      snap.TopN.String(),
		AvgNm:     snap.AvgNM.String(),
		FetchedAt: snap.FetchedAt.Unix(),
	}, nil
}

// RegisterHealth registers the standard grpc_health_v1 health server and sets the serving status to SERVING
func RegisterHealth(s *grpc.Server) {
	hs := health.NewServer()
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, hs)
}
