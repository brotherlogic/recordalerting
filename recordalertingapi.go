package main

import (
	"golang.org/x/net/context"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

//ClientUpdate process new record
func (s *Server) ClientUpdate(ctx context.Context, req *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	r, err := s.rc.getRecord(ctx, req.GetInstanceId())
	if err != nil {
		return nil, err
	}

	config, err := s.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	err = s.assessRecord(ctx, config, r)

	errt := s.saveConfig(ctx, config)
	if errt != nil {
		return nil, errt
	}

	return &rcpb.ClientUpdateResponse{}, err
}
