package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

func (s *Server) ListExecutionFunctions(_ context.Context, _ *execution.ListExecutionFunctionsRequest) (*execution.ListExecutionFunctionsResponse, error) {
	return &execution.ListExecutionFunctionsResponse{
		Functions: s.ListActionFunctions(),
	}, nil
}

func (s *Server) ListExecutionMethods(_ context.Context, _ *execution.ListExecutionMethodsRequest) (*execution.ListExecutionMethodsResponse, error) {
	return &execution.ListExecutionMethodsResponse{
		Methods: s.ListGRPCMethods(),
	}, nil
}

func (s *Server) ListExecutionServices(_ context.Context, _ *execution.ListExecutionServicesRequest) (*execution.ListExecutionServicesResponse, error) {
	return &execution.ListExecutionServicesResponse{
		Services: s.ListGRPCServices(),
	}, nil
}

func (s *Server) SetExecution(ctx context.Context, req *execution.SetExecutionRequest) (*execution.SetExecutionResponse, error) {
	set := &command.SetExecution{
		Targets:  req.GetTargets(),
		Includes: req.GetIncludes(),
	}

	var err error
	var details *domain.ObjectDetails
	switch t := req.GetCondition().GetConditionType().(type) {
	case *execution.SetConditions_Request:
		cond := &command.ExecutionAPICondition{
			Method:  t.Request.GetMethod(),
			Service: t.Request.GetService(),
			All:     t.Request.GetAll(),
		}
		details, err = s.command.SetExecutionRequest(ctx, cond, set, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *execution.SetConditions_Response:
		cond := &command.ExecutionAPICondition{
			Method:  t.Response.GetMethod(),
			Service: t.Response.GetService(),
			All:     t.Response.GetAll(),
		}
		details, err = s.command.SetExecutionResponse(ctx, cond, set, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *execution.SetConditions_Event:
		cond := &command.ExecutionEventCondition{
			Event: t.Event.GetEvent(),
			Group: t.Event.GetGroup(),
			All:   t.Event.GetAll(),
		}
		details, err = s.command.SetExecutionEvent(ctx, cond, set, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *execution.SetConditions_Function:
		details, err = s.command.SetExecutionFunction(ctx, command.ExecutionFunctionCondition(t.Function), set, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	}
	return &execution.SetExecutionResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeleteExecution(ctx context.Context, req *execution.DeleteExecutionRequest) (*execution.DeleteExecutionResponse, error) {
	var err error
	var details *domain.ObjectDetails
	switch t := req.GetCondition().GetConditionType().(type) {
	case *execution.SetConditions_Request:
		cond := &command.ExecutionAPICondition{
			Method:  t.Request.GetMethod(),
			Service: t.Request.GetService(),
			All:     t.Request.GetAll(),
		}
		details, err = s.command.DeleteExecutionRequest(ctx, cond, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *execution.SetConditions_Response:
		cond := &command.ExecutionAPICondition{
			Method:  t.Response.GetMethod(),
			Service: t.Response.GetService(),
			All:     t.Response.GetAll(),
		}
		details, err = s.command.DeleteExecutionResponse(ctx, cond, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *execution.SetConditions_Event:
		cond := &command.ExecutionEventCondition{
			Event: t.Event.GetEvent(),
			Group: t.Event.GetGroup(),
			All:   t.Event.GetAll(),
		}
		details, err = s.command.DeleteExecutionEvent(ctx, cond, authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	case *execution.SetConditions_Function:
		details, err = s.command.DeleteExecutionFunction(ctx, command.ExecutionFunctionCondition(t.Function), authz.GetInstance(ctx).InstanceID())
		if err != nil {
			return nil, err
		}
	}
	return &execution.DeleteExecutionResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}
