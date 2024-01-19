package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/juliansunn/simple_bank/db/sqlc"
	"github.com/juliansunn/simple_bank/pb"
	"github.com/juliansunn/simple_bank/util"
	"github.com/juliansunn/simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO: add authorization layer
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, InvalidArgumentError(violations)
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != "",
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != "",
		},
	}

	if req.Password != "" {
		hashedPassword, err := util.HashedPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}
	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil

}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, FieldViolation("username", err))
	}
	if req.GetPassword() != "" {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, FieldViolation("password", err))
		}
	}
	if req.GetFullName() != "" {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, FieldViolation("full_name", err))
		}
	}
	if req.GetEmail() != "" {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, FieldViolation("email", err))
		}
	}
	return violations
}
