package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

//ClientUpdate process new record
func (s *Server) ClientUpdate(ctx context.Context, req *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	config, err := s.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	r, err := s.rc.getRecord(ctx, req.GetInstanceId())
	if err != nil {
		if status.Convert(err).Code() == codes.OutOfRange {
			for _, problem := range config.GetProblems() {
				if problem.GetInstanceId() == req.GetInstanceId() {
					err := s.DeleteIssue(ctx, problem.GetIssueNumber())
					if err != nil {
						return nil, err
					}
				}
			}
			return &rcpb.ClientUpdateResponse{}, nil
		}
		return nil, err
	}

	// Don't validate records that need a gram pull
	if r.GetMetadata().GetNeedsGramUpdate() {
		return &rcpb.ClientUpdateResponse{}, nil
	}

	err = s.assessRecord(ctx, config, r)

	errt := s.saveConfig(ctx, config)
	if errt != nil {
		return nil, errt
	}

	return &rcpb.ClientUpdateResponse{}, err
}
