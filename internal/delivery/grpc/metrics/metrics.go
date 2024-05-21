package metrics

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dkmelnik/metrics/internal/apperrors"
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/models"
	grpcmetrics "github.com/dkmelnik/metrics/proto/metrics"
)

type Handler struct {
	grpcmetrics.MetricsServer
	pgDB    *sqlx.DB
	service *metrics.Service
}

func NewHandler(pgDB *sqlx.DB, s *metrics.Service) *Handler {
	return &Handler{
		pgDB:    pgDB,
		service: s,
	}
}

func (h *Handler) Create(ctx context.Context, req *grpcmetrics.CreateRequest) (*grpcmetrics.CreateResponse, error) {
	if req.Delta == 0 && req.Value == 0 {
		return nil, status.Error(codes.InvalidArgument, "must not be zero")
	}
	model, err := models.NewMetric(req.Id, req.Mtype)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.Delta != 0 {
		model.SetDelta(req.Delta)
	}
	if req.Value != 0 {
		model.SetValue(req.Value)
	}

	if err = h.service.CreateOrUpdate(model); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrTypeNotCorrect):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, apperrors.ErrParse):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, codes.Internal.String())
		}
	}

	return &grpcmetrics.CreateResponse{
		Id:    model.ID,
		Mtype: model.MType,
		Delta: model.Delta.Int64,
		Value: model.Value.Float64,
	}, nil
}
func (h *Handler) CreateMany(ctx context.Context, req *grpcmetrics.CreateManyRequest) (*grpcmetrics.CreateManyResponse, error) {
	if len(req.Metrics) == 0 {
		return nil, status.Error(codes.InvalidArgument, "metrics list empty")
	}
	var mds = make([]models.Metric, 0, len(req.Metrics))
	for _, v := range req.Metrics {
		if v.Delta == 0 && v.Value == 0 {
			return nil, status.Error(codes.InvalidArgument, "must not be zero")
		}

		model, err := models.NewMetric(v.Id, v.Mtype)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if v.Delta != 0 {
			model.SetDelta(v.Delta)
		}
		if v.Value != 0 {
			model.SetValue(v.Value)
		}

		mds = append(mds, model)
	}
	if err := h.service.CreateOrUpdateMany(mds); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrTypeNotCorrect):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, apperrors.ErrParse):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, codes.Internal.String())
		}
	}

	var out = make([]*grpcmetrics.CreateResponse, len(mds))
	for idx, v := range mds {
		out[idx] = &grpcmetrics.CreateResponse{
			Id:    v.ID,
			Mtype: v.MType,
			Delta: v.Delta.Int64,
			Value: v.Value.Float64,
		}
	}
	return &grpcmetrics.CreateManyResponse{Metrics: out}, nil
}
func (h *Handler) Get(ctx context.Context, req *grpcmetrics.GetRequest) (*grpcmetrics.GetResponse, error) {
	value, err := h.service.GetMetric(req.Mtype, req.Id)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			return nil, status.Error(codes.NotFound, codes.NotFound.String())
		default:
			return nil, status.Error(codes.Internal, codes.Internal.String())
		}
	}
	out := &grpcmetrics.GetResponse{
		Id:    value.ID,
		Mtype: value.MType,
	}
	if value.Delta != nil {
		out.Delta = *value.Delta
	}
	if value.Value != nil {
		out.Value = *value.Value
	}

	return out, nil
}
func (h *Handler) List(ctx context.Context, req *grpcmetrics.ListRequest) (*grpcmetrics.ListResponse, error) {
	list, err := h.service.GetAllInHTML()
	if err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}
	return &grpcmetrics.ListResponse{Html: list}, nil
}
func (h *Handler) CheckPGSQL(ctx context.Context, req *grpcmetrics.CheckPGSQLRequest) (*grpcmetrics.CheckPGSQLResponse, error) {
	if h.pgDB == nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	if err := h.pgDB.Ping(); err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &grpcmetrics.CheckPGSQLResponse{}, nil
}
