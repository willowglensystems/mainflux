// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"

	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/things"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ mainflux.ThingsServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	canAccessByKey kitgrpc.Handler
	canAccessByID  kitgrpc.Handler
	isChannelOwner kitgrpc.Handler
	identify       kitgrpc.Handler
}

// NewServer returns new ThingsServiceServer instance.
func NewServer(tracer opentracing.Tracer, svc things.Service) mainflux.ThingsServiceServer {
	return &grpcServer{
		canAccessByKey: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "can_access")(canAccessEndpoint(svc)),
			decodeCanAccessByKeyRequest,
			encodeIdentityResponse,
		),
		canAccessByID: kitgrpc.NewServer(
			canAccessByIDEndpoint(svc),
			decodeCanAccessByIDRequest,
			encodeEmptyResponse,
		),
		isChannelOwner: kitgrpc.NewServer(
			isChannelOwnerEndpoint(svc),
			decodeIsChannelOwnerRequest,
			encodeEmptyResponse,
		),
		identify: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "identify")(identifyEndpoint(svc)),
			decodeIdentifyRequest,
			encodeIdentityResponse,
		),
	}
}

func (gs *grpcServer) CanAccessByKey(ctx context.Context, req *mainflux.AccessByKeyReq) (*mainflux.ThingID, error) {
	_, res, err := gs.canAccessByKey.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*mainflux.ThingID), nil
}

func (gs *grpcServer) CanAccessByID(ctx context.Context, req *mainflux.AccessByIDReq) (*empty.Empty, error) {
	_, res, err := gs.canAccessByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*empty.Empty), nil
}

func (gs *grpcServer) IsChannelOwner(ctx context.Context, req *mainflux.ChannelOwnerReq) (*empty.Empty, error) {
	_, res, err := gs.isChannelOwner.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*empty.Empty), nil
}

func (gs *grpcServer) Identify(ctx context.Context, req *mainflux.Token) (*mainflux.ThingID, error) {
	_, res, err := gs.identify.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*mainflux.ThingID), nil
}

func decodeCanAccessByKeyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*mainflux.AccessByKeyReq)
	return AccessByKeyReq{thingKey: req.GetToken(), chanID: req.GetChanID()}, nil
}

func decodeCanAccessByIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*mainflux.AccessByIDReq)
	return accessByIDReq{thingID: req.GetThingID(), chanID: req.GetChanID()}, nil
}

func decodeIsChannelOwnerRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*mainflux.ChannelOwnerReq)
	return channelOwnerReq{owner: req.GetOwner(), chanID: req.GetChanID()}, nil
}

func decodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*mainflux.Token)
	return identifyReq{key: req.GetValue()}, nil
}

func encodeIdentityResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(identityRes)
	return &mainflux.ThingID{Value: res.id}, nil
}

func encodeEmptyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(emptyRes)
	return &empty.Empty{}, encodeError(res.err)
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case things.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid can access request")
	case things.ErrUnauthorizedAccess:
		return status.Error(codes.PermissionDenied, "missing or invalid credentials provided")
	case things.ErrEntityConnected:
		return status.Error(codes.PermissionDenied, "entities are not connected")
	case things.ErrNotFound:
		return status.Error(codes.NotFound, "entity does not exist")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
