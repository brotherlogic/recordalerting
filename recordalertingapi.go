package main

import (
	rcpb "github.com/brotherlogic/recordcollection/proto"
	"golang.org/x/net/context"
)

//ClientUpdate process new record
func (s *Server) ClientUpdate(ctx context.Context, req *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	r, err := s.rc.getRecord(ctx, req.GetInstanceId())
	if err != nil {
		return nil, err
	}

	return &rcpb.ClientUpdateResponse{}, s.assessRecord(ctx, r)
}
