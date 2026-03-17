package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/config"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/domain"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/dto"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/model"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/repository"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/apperrors"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/auth"
)

type TokenBundle struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type AuthResult struct {
	User   *model.User  `json:"user"`
	Rider  *model.Rider `json:"rider,omitempty"`
	Tokens TokenBundle  `json:"tokens"`
}

type OrderRequestView struct {
	Assignment model.DeliveryAssignment `json:"assignment"`
	Order      model.Order              `json:"order"`
	Items      []model.OrderItem        `json:"items"`
}

type EarningsSlice struct {
	Total   float64              `json:"total"`
	Count   int                  `json:"count"`
	Records []model.RiderEarning `json:"records"`
}

type Services struct {
	Auth         *AuthService
	Rider        *RiderService
	Shift        *ShiftService
	Dispatch     *DispatchService
	Finance      *FinanceService
	Feedback     *FeedbackService
	Notification *NotificationService
	Support      *SupportService
	Admin        *AdminService
}

type AuthService struct {
	repo    *repository.Repository
	cfg     config.Config
	manager *auth.Manager
}

type RiderService struct {
	repo *repository.Repository
}

type ShiftService struct {
	repo *repository.Repository
	cfg  config.Config
}

type DispatchService struct {
	repo *repository.Repository
	cfg  config.Config
}

type FinanceService struct {
	repo *repository.Repository
	cfg  config.Config
}

type FeedbackService struct {
	repo *repository.Repository
}

type NotificationService struct {
	repo *repository.Repository
}

type SupportService struct {
	repo *repository.Repository
}

type AdminService struct {
	repo *repository.Repository
	cfg  config.Config
}

func NewServices(repo *repository.Repository, cfg config.Config, manager *auth.Manager) *Services {
	return &Services{
		Auth:         &AuthService{repo: repo, cfg: cfg, manager: manager},
		Rider:        &RiderService{repo: repo},
		Shift:        &ShiftService{repo: repo, cfg: cfg},
		Dispatch:     &DispatchService{repo: repo, cfg: cfg},
		Finance:      &FinanceService{repo: repo, cfg: cfg},
		Feedback:     &FeedbackService{repo: repo},
		Notification: &NotificationService{repo: repo},
		Support:      &SupportService{repo: repo},
		Admin:        &AdminService{repo: repo, cfg: cfg},
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest, clientIP string) (*AuthResult, error) {
	user, rider, err := s.repo.FindUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid credentials")
	}
	if !auth.ComparePassword(user.PasswordHash, req.Password) {
		return nil, apperrors.Unauthorized("invalid credentials")
	}
	if user.Status != domain.UserStatusActive {
		return nil, apperrors.Forbidden("user account is not active")
	}
	tokens, err := s.issueTokens(ctx, user, req.DeviceID, req.DeviceName, clientIP)
	if err != nil {
		return nil, err
	}
	return &AuthResult{User: user, Rider: rider, Tokens: tokens}, nil
}

func (s *AuthService) SendOTP(ctx context.Context, req dto.OTPRequest) (map[string]any, error) {
	user, _, err := s.repo.FindUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, apperrors.NotFound("USER_NOT_FOUND", "user not found")
	}
	code := generateOTPCode()
	s.repo.SaveOTP(ctx, model.OTP{UserID: user.ID, Purpose: domain.OTPPurposeLogin, Code: code, Status: domain.OTPStatusPending, ExpiresAt: time.Now().UTC().Add(s.cfg.Auth.OTPExpiry)})
	response := map[string]any{"expires_in_seconds": int(s.cfg.Auth.OTPExpiry.Seconds()), "channel": "SMS"}
	if s.cfg.Auth.ExposeOTPInDevMode {
		response["otp"] = code
	}
	return response, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, req dto.OTPVerifyRequest, clientIP string) (*AuthResult, error) {
	user, rider, err := s.repo.FindUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, apperrors.NotFound("USER_NOT_FOUND", "user not found")
	}
	otp, err := s.repo.GetLatestOTP(ctx, domain.OTPPurposeLogin, user.ID)
	if err != nil {
		return nil, apperrors.NotFound("OTP_NOT_FOUND", "otp not found")
	}
	if otp.Status == domain.OTPStatusLocked {
		return nil, apperrors.Forbidden("otp is locked")
	}
	if time.Now().UTC().After(otp.ExpiresAt) {
		_, _ = s.repo.UpdateOTP(ctx, domain.OTPPurposeLogin, user.ID, func(item *model.OTP) { item.Status = domain.OTPStatusExpired })
		return nil, apperrors.New(400, "OTP_EXPIRED", "otp expired")
	}
	if otp.Code != req.OTP {
		updated, _ := s.repo.UpdateOTP(ctx, domain.OTPPurposeLogin, user.ID, func(item *model.OTP) {
			item.RetryCount++
			if item.RetryCount >= s.cfg.Auth.OTPMaxRetries {
				item.Status = domain.OTPStatusLocked
			}
		})
		if updated != nil && updated.Status == domain.OTPStatusLocked {
			return nil, apperrors.New(429, "OTP_LOCKED", "otp retry limit exceeded")
		}
		return nil, apperrors.New(400, "INVALID_OTP", "invalid otp")
	}
	_, _ = s.repo.UpdateOTP(ctx, domain.OTPPurposeLogin, user.ID, func(item *model.OTP) {
		now := time.Now().UTC()
		item.Status = domain.OTPStatusVerified
		item.VerifiedAt = &now
	})
	tokens, err := s.issueTokens(ctx, user, req.DeviceID, req.DeviceName, clientIP)
	if err != nil {
		return nil, err
	}
	return &AuthResult{User: user, Rider: rider, Tokens: tokens}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (TokenBundle, error) {
	refreshHash := shaToken(req.RefreshToken)
	session, err := s.repo.GetSession(ctx, refreshHash)
	if err != nil || session.RevokedAt != nil || time.Now().UTC().After(session.ExpiresAt) {
		return TokenBundle{}, apperrors.Unauthorized("invalid refresh token")
	}
	if session.DeviceID != req.DeviceID {
		return TokenBundle{}, apperrors.Forbidden("refresh token does not belong to this device")
	}
	user, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return TokenBundle{}, apperrors.Unauthorized("user not found for refresh token")
	}
	accessToken, err := s.manager.Issue(auth.Claims{Subject: user.ID, Email: user.Email, Roles: user.Roles, DeviceID: req.DeviceID, IssuedAt: time.Now().UTC().Unix(), ExpiresAt: time.Now().UTC().Add(s.cfg.Auth.AccessTokenTTL).Unix()})
	if err != nil {
		return TokenBundle{}, apperrors.New(500, "TOKEN_ISSUE_FAILED", "failed to issue access token")
	}
	return TokenBundle{AccessToken: accessToken, RefreshToken: req.RefreshToken, ExpiresIn: int64(s.cfg.Auth.AccessTokenTTL.Seconds()), TokenType: "Bearer"}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) {
	s.repo.RevokeSession(ctx, shaToken(refreshToken))
}

func (s *AuthService) LogoutAll(ctx context.Context, userID string) {
	s.repo.RevokeAllSessions(ctx, userID)
}

func (s *AuthService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (map[string]any, error) {
	return s.SendOTP(ctx, dto.OTPRequest{Login: req.Login})
}

func (s *AuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	user, _, err := s.repo.FindUserByLogin(ctx, req.Login)
	if err != nil {
		return apperrors.NotFound("USER_NOT_FOUND", "user not found")
	}
	otp, err := s.repo.GetLatestOTP(ctx, domain.OTPPurposeLogin, user.ID)
	if err != nil || otp.Code != req.OTP || time.Now().UTC().After(otp.ExpiresAt) {
		return apperrors.New(400, "INVALID_OTP", "invalid or expired otp")
	}
	passwordHash, err := auth.HashPassword(req.NewPassword, s.cfg.Auth.PasswordHashCost)
	if err != nil {
		return apperrors.New(500, "PASSWORD_HASH_FAILED", "failed to hash password")
	}
	return s.repo.UpdatePassword(ctx, user.ID, passwordHash)
}

func (s *AuthService) Me(ctx context.Context, userID string) (map[string]any, error) {
	user, rider, err := s.repo.GetRiderProfile(ctx, userID)
	if err == nil {
		return map[string]any{"user": user, "rider": rider}, nil
	}
	userOnly, userErr := s.repo.GetUserByID(ctx, userID)
	if userErr != nil {
		return nil, apperrors.NotFound("USER_NOT_FOUND", "user not found")
	}
	return map[string]any{"user": userOnly}, nil
}

func (s *AuthService) issueTokens(ctx context.Context, user *model.User, deviceID, deviceName, clientIP string) (TokenBundle, error) {
	accessToken, err := s.manager.Issue(auth.Claims{Subject: user.ID, Email: user.Email, Roles: user.Roles, DeviceID: deviceID, IssuedAt: time.Now().UTC().Unix(), ExpiresAt: time.Now().UTC().Add(s.cfg.Auth.AccessTokenTTL).Unix()})
	if err != nil {
		return TokenBundle{}, apperrors.New(500, "TOKEN_ISSUE_FAILED", "failed to issue access token")
	}
	refreshToken, refreshHash := auth.GenerateOpaqueToken()
	s.repo.SaveSession(ctx, model.Session{UserID: user.ID, DeviceID: deviceID, DeviceName: deviceName, RefreshHash: refreshHash, LastIP: clientIP, ExpiresAt: time.Now().UTC().Add(s.cfg.Auth.RefreshTokenTTL)})
	return TokenBundle{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: int64(s.cfg.Auth.AccessTokenTTL.Seconds()), TokenType: "Bearer"}, nil
}

func (s *RiderService) GetProfile(ctx context.Context, userID string) (map[string]any, error) {
	user, rider, err := s.repo.GetRiderProfile(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	vehicle := s.repo.RiderVehicles[rider.ID]
	bank := s.repo.RiderBankAccounts[rider.ID]
	return map[string]any{"user": user, "rider": rider, "vehicle": vehicle, "bank_account": bank, "documents": s.repo.ListDocuments(ctx, rider.ID)}, nil
}

func (s *RiderService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateRiderProfileRequest) (*model.User, error) {
	user, err := s.repo.UpdateUserProfile(ctx, userID, req.FirstName, req.LastName, req.Email, req.Phone)
	if err != nil {
		return nil, apperrors.NotFound("USER_NOT_FOUND", "user not found")
	}
	return user, nil
}

func (s *RiderService) UpdatePhoto(ctx context.Context, userID, photoURL string) (*model.User, error) {
	user, err := s.repo.UpdatePhoto(ctx, userID, photoURL)
	if err != nil {
		return nil, apperrors.NotFound("USER_NOT_FOUND", "user not found")
	}
	return user, nil
}

func (s *RiderService) UpdateVehicle(ctx context.Context, userID string, req dto.UpdateVehicleRequest) (*model.RiderVehicle, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	insuranceExpiry, parseErr := time.Parse("2006-01-02", req.InsuranceExpiry)
	if parseErr != nil {
		return nil, apperrors.New(400, "INVALID_DATE", "insurance_expiry must be YYYY-MM-DD")
	}
	return s.repo.UpdateVehicle(ctx, rider.ID, model.RiderVehicle{VehicleType: req.VehicleType, RegistrationNo: req.RegistrationNo, Color: req.Color, Capacity: req.Capacity, InsuranceExpiry: insuranceExpiry})
}

func (s *RiderService) UpdateDocuments(ctx context.Context, userID string, req dto.UpdateDocumentsRequest) ([]model.RiderDocument, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	docs := make([]model.RiderDocument, 0, len(req.Documents))
	for _, item := range req.Documents {
		docs = append(docs, model.RiderDocument{DocumentType: item.DocumentType, DocumentURL: item.DocumentURL})
	}
	return s.repo.ReplaceDocuments(ctx, rider.ID, docs), nil
}

func (s *RiderService) GetDocuments(ctx context.Context, userID string) ([]model.RiderDocument, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.ListDocuments(ctx, rider.ID), nil
}

func (s *RiderService) UpsertBankAccount(ctx context.Context, userID string, req dto.UpdateBankAccountRequest) (*model.RiderBankAccount, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.UpsertBankAccount(ctx, rider.ID, model.RiderBankAccount{AccountHolder: req.AccountHolder, BankName: req.BankName, AccountNumber: req.AccountNumber, IFSCCode: req.IFSCCode, UPIID: req.UPIID}), nil
}

func (s *RiderService) GetStatus(ctx context.Context, userID string) (map[string]any, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	shift, _ := s.repo.GetActiveShift(ctx, rider.ID)
	return map[string]any{"status": rider.Status, "availability_status": rider.AvailabilityStatus, "kyc_status": rider.KYCStatus, "approval_status": rider.ApprovalStatus, "active_shift": shift}, nil
}

func (s *ShiftService) GoOnline(ctx context.Context, userID, requestID string) (map[string]any, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	if rider.Status != domain.UserStatusActive || rider.KYCStatus != "APPROVED" || rider.ApprovalStatus != "APPROVED" {
		return nil, apperrors.Forbidden("rider is not eligible to go online")
	}
	assignments := s.repo.GetRestaurantAssignments(ctx, rider.ID)
	if len(assignments) == 0 {
		return nil, apperrors.Forbidden("rider is not assigned to any restaurant")
	}
	shift, shiftErr := s.repo.GetActiveShift(ctx, rider.ID)
	if shiftErr != nil {
		if !s.cfg.Dispatch.ShiftAutoStart {
			return nil, apperrors.Conflict("SHIFT_REQUIRED", "an active shift is required")
		}
		started := s.repo.StartShift(ctx, rider.ID, userID)
		shift = &started
	}
	updated, err := s.repo.SetRiderAvailability(ctx, rider.ID, domain.RiderAvailabilityOnline)
	if err != nil {
		return nil, apperrors.New(500, "RIDER_STATUS_UPDATE_FAILED", "failed to update rider status")
	}
	s.repo.AppendAudit(ctx, model.AuditLog{ActorID: userID, ActorRole: domain.RoleRider, Action: "RIDER_WENT_ONLINE", EntityType: "RIDER", EntityID: rider.ID, RequestID: requestID, Metadata: map[string]any{"shift_id": shift.ID}})
	return map[string]any{"rider": updated, "shift": shift, "available_for_assignment": true}, nil
}

func (s *ShiftService) GoOffline(ctx context.Context, userID, requestID string) (map[string]any, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	updated, err := s.repo.SetRiderAvailability(ctx, rider.ID, domain.RiderAvailabilityOffline)
	if err != nil {
		return nil, apperrors.New(500, "RIDER_STATUS_UPDATE_FAILED", "failed to update rider status")
	}
	s.repo.AppendAudit(ctx, model.AuditLog{ActorID: userID, ActorRole: domain.RoleRider, Action: "RIDER_WENT_OFFLINE", EntityType: "RIDER", EntityID: rider.ID, RequestID: requestID})
	return map[string]any{"rider": updated, "available_for_assignment": false}, nil
}

func (s *ShiftService) StartShift(ctx context.Context, userID string) (*model.RiderShift, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	if _, err := s.repo.GetActiveShift(ctx, rider.ID); err == nil {
		return nil, apperrors.Conflict("SHIFT_ALREADY_ACTIVE", "active shift already exists")
	}
	shift := s.repo.StartShift(ctx, rider.ID, userID)
	return &shift, nil
}

func (s *ShiftService) EndShift(ctx context.Context, userID string) (*model.RiderShift, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	shift, err := s.repo.EndShift(ctx, rider.ID, userID)
	if err != nil {
		return nil, apperrors.Conflict("NO_ACTIVE_SHIFT", err.Error())
	}
	_, _ = s.repo.SetRiderAvailability(ctx, rider.ID, domain.RiderAvailabilityOffline)
	return shift, nil
}

func (s *ShiftService) StartBreak(ctx context.Context, userID string, req dto.BreakRequest) (*model.RiderBreak, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	shift, err := s.repo.GetActiveShift(ctx, rider.ID)
	if err != nil {
		return nil, apperrors.Conflict("SHIFT_REQUIRED", "active shift required")
	}
	brk := s.repo.StartBreak(ctx, rider.ID, shift.ID, req.Reason)
	return &brk, nil
}

func (s *ShiftService) EndBreak(ctx context.Context, userID string) (*model.RiderBreak, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	brk, err := s.repo.EndBreak(ctx, rider.ID)
	if err != nil {
		return nil, apperrors.Conflict("NO_ACTIVE_BREAK", err.Error())
	}
	return brk, nil
}

func (s *ShiftService) GetTodayShift(ctx context.Context, userID string) ([]model.RiderShift, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	today := time.Now().UTC().Format("2006-01-02")
	items := make([]model.RiderShift, 0)
	for _, shift := range s.repo.ListShifts(ctx, rider.ID) {
		if shift.StartedAt.UTC().Format("2006-01-02") == today {
			items = append(items, shift)
		}
	}
	return items, nil
}

func (s *ShiftService) GetShiftHistory(ctx context.Context, userID string) ([]model.RiderShift, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.ListShifts(ctx, rider.ID), nil
}

func shaToken(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func generateOTPCode() string {
	return fmt.Sprintf("%06d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000))
}
func (s *DispatchService) ListOrderRequests(ctx context.Context, userID string) ([]OrderRequestView, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	assignments := s.repo.ListPendingAssignments(ctx, rider.ID)
	views := make([]OrderRequestView, 0, len(assignments))
	for _, assignment := range assignments {
		order, items, orderErr := s.repo.GetOrder(ctx, assignment.OrderID)
		if orderErr == nil {
			views = append(views, OrderRequestView{Assignment: assignment, Order: *order, Items: items})
		}
	}
	return views, nil
}

func (s *DispatchService) GetOrderRequest(ctx context.Context, userID, assignmentID string) (*OrderRequestView, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	assignment, err := s.repo.GetAssignment(ctx, assignmentID)
	if err != nil || assignment.RiderID != rider.ID {
		return nil, apperrors.NotFound("ASSIGNMENT_NOT_FOUND", "order request not found")
	}
	order, items, err := s.repo.GetOrder(ctx, assignment.OrderID)
	if err != nil {
		return nil, apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	return &OrderRequestView{Assignment: *assignment, Order: *order, Items: items}, nil
}

func (s *DispatchService) AcceptOrderRequest(ctx context.Context, userID, assignmentID, requestID string) (*OrderRequestView, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	assignment, err := s.repo.GetAssignment(ctx, assignmentID)
	if err != nil || assignment.RiderID != rider.ID {
		return nil, apperrors.NotFound("ASSIGNMENT_NOT_FOUND", "order request not found")
	}
	if time.Now().UTC().After(assignment.DecisionDeadlineAt) {
		return nil, apperrors.Conflict("ORDER_REQUEST_EXPIRED", "order request expired")
	}
	if _, err := s.repo.GetActiveShift(ctx, rider.ID); err != nil {
		return nil, apperrors.Conflict("SHIFT_REQUIRED", "active shift required to accept order")
	}
	accepted, err := s.repo.AcceptAssignment(ctx, assignmentID)
	if err != nil {
		return nil, apperrors.Conflict("ORDER_ACCEPT_FAILED", err.Error())
	}
	_, _ = s.repo.UpdateOrderStatus(ctx, accepted.OrderID, domain.DeliveryStatusAccepted)
	s.repo.AddOrderHistory(ctx, model.DeliveryStatusHistory{OrderID: accepted.OrderID, AssignmentID: accepted.ID, Status: domain.DeliveryStatusAccepted, ActorID: userID, ActorRole: domain.RoleRider, Comment: "Rider accepted order request"})
	order, items, _ := s.repo.GetOrder(ctx, accepted.OrderID)
	if order.PickupOTPRequired {
		_, _ = s.issueOrderOTP(ctx, order.ID, domain.OTPPurposePickup)
	}
	if order.DeliveryOTPRequired {
		_, _ = s.issueOrderOTP(ctx, order.ID, domain.OTPPurposeDelivery)
	}
	s.repo.AppendAudit(ctx, model.AuditLog{ActorID: userID, ActorRole: domain.RoleRider, Action: "ORDER_ACCEPTED", EntityType: "ORDER", EntityID: order.ID, RequestID: requestID, Metadata: map[string]any{"assignment_id": accepted.ID}})
	return &OrderRequestView{Assignment: *accepted, Order: *order, Items: items}, nil
}

func (s *DispatchService) RejectOrderRequest(ctx context.Context, userID, assignmentID, reason, requestID string) (*model.DeliveryAssignment, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	assignment, err := s.repo.GetAssignment(ctx, assignmentID)
	if err != nil || assignment.RiderID != rider.ID {
		return nil, apperrors.NotFound("ASSIGNMENT_NOT_FOUND", "order request not found")
	}
	rejected, err := s.repo.RejectAssignment(ctx, assignmentID, reason)
	if err != nil {
		return nil, apperrors.Conflict("ORDER_REJECT_FAILED", err.Error())
	}
	s.repo.AddOrderHistory(ctx, model.DeliveryStatusHistory{OrderID: rejected.OrderID, AssignmentID: rejected.ID, Status: domain.DeliveryStatusReadyForAssignment, ActorID: userID, ActorRole: domain.RoleRider, Comment: reason})
	s.repo.AppendAudit(ctx, model.AuditLog{ActorID: userID, ActorRole: domain.RoleRider, Action: "ORDER_REJECTED", EntityType: "ORDER", EntityID: rejected.OrderID, RequestID: requestID, Metadata: map[string]any{"assignment_id": rejected.ID, "reason": reason}})
	return rejected, nil
}

func (s *DispatchService) GetActiveOrder(ctx context.Context, userID string) (map[string]any, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	assignment, err := s.repo.GetActiveAssignmentForRider(ctx, rider.ID)
	if err != nil {
		return nil, apperrors.NotFound("ACTIVE_ORDER_NOT_FOUND", "active order not found")
	}
	order, items, _ := s.repo.GetOrder(ctx, assignment.OrderID)
	return map[string]any{"assignment": assignment, "order": order, "items": items}, nil
}

func (s *DispatchService) ListAssignedOrders(ctx context.Context, userID string) ([]model.Order, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.ListAssignedOrders(ctx, rider.ID), nil
}

func (s *DispatchService) ListOrderHistory(ctx context.Context, userID string) ([]model.Order, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.ListOrderHistoryByRider(ctx, rider.ID), nil
}

func (s *DispatchService) GetOrder(ctx context.Context, userID, orderID string) (map[string]any, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	assignment, err := s.repo.FindAssignmentByOrderAndRider(ctx, orderID, rider.ID)
	if err != nil {
		return nil, apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	order, items, _ := s.repo.GetOrder(ctx, orderID)
	return map[string]any{"assignment": assignment, "order": order, "items": items}, nil
}

func (s *DispatchService) MarkArrivedRestaurant(ctx context.Context, userID, orderID string) (*model.Order, error) {
	return s.transitionOrder(ctx, userID, orderID, domain.DeliveryStatusReachedRestaurant, "Reached restaurant")
}

func (s *DispatchService) MarkPickedUp(ctx context.Context, userID, orderID string) (*model.Order, error) {
	order, _, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	if order.PickupOTPRequired && order.Status != domain.DeliveryStatusPickupVerified {
		return nil, apperrors.Conflict("PICKUP_OTP_REQUIRED", "pickup otp verification required before pickup")
	}
	return s.transitionOrder(ctx, userID, orderID, domain.DeliveryStatusPickedUp, "Order picked up from restaurant")
}

func (s *DispatchService) MarkOnTheWay(ctx context.Context, userID, orderID string) (*model.Order, error) {
	return s.transitionOrder(ctx, userID, orderID, domain.DeliveryStatusOnTheWay, "Rider started delivery route")
}

func (s *DispatchService) MarkArrivedCustomer(ctx context.Context, userID, orderID string) (*model.Order, error) {
	return s.transitionOrder(ctx, userID, orderID, domain.DeliveryStatusReachedCustomer, "Rider reached customer location")
}
func (s *DispatchService) MarkDelivered(ctx context.Context, userID, orderID string) (*model.Order, error) {
	order, _, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	if order.DeliveryOTPRequired && order.Status != domain.DeliveryStatusDeliveryVerified {
		return nil, apperrors.Conflict("DELIVERY_OTP_REQUIRED", "delivery otp verification required before completing order")
	}
	delivered, err := s.transitionOrder(ctx, userID, orderID, domain.DeliveryStatusDelivered, "Order delivered successfully")
	if err != nil {
		return nil, err
	}
	rider, _ := s.repo.GetRiderByUserID(ctx, userID)
	assignment, _ := s.repo.FindAssignmentByOrderAndRider(ctx, orderID, rider.ID)
	_, _ = s.repo.UpdateAssignmentStatus(ctx, assignment.ID, domain.AssignmentStatusCompleted, "")
	_, _ = s.repo.SetRiderAvailability(ctx, rider.ID, domain.RiderAvailabilityOnline)
	s.finalizeEarnings(ctx, rider.ID, delivered)
	return delivered, nil
}

func (s *DispatchService) MarkFailed(ctx context.Context, userID, orderID, reason string) (*model.Order, error) {
	failed, err := s.transitionOrder(ctx, userID, orderID, domain.DeliveryStatusFailed, reason)
	if err != nil {
		return nil, err
	}
	rider, _ := s.repo.GetRiderByUserID(ctx, userID)
	assignment, _ := s.repo.FindAssignmentByOrderAndRider(ctx, orderID, rider.ID)
	_, _ = s.repo.UpdateAssignmentStatus(ctx, assignment.ID, domain.AssignmentStatusFailed, reason)
	_, _ = s.repo.SetRiderAvailability(ctx, rider.ID, domain.RiderAvailabilityOnline)
	return failed, nil
}

func (s *DispatchService) RequestCancel(ctx context.Context, userID, orderID, reason string) (map[string]any, error) {
	rider, _ := s.repo.GetRiderByUserID(ctx, userID)
	assignment, err := s.repo.FindAssignmentByOrderAndRider(ctx, orderID, rider.ID)
	if err != nil {
		return nil, apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	s.repo.AddOrderHistory(ctx, model.DeliveryStatusHistory{OrderID: orderID, AssignmentID: assignment.ID, Status: domain.DeliveryStatusCancelled, ActorID: userID, ActorRole: domain.RoleRider, Comment: reason})
	return map[string]any{"order_id": orderID, "requested": true, "reason": reason}, nil
}

func (s *DispatchService) GetTimeline(ctx context.Context, orderID string) []model.DeliveryStatusHistory {
	return s.repo.ListOrderHistory(ctx, orderID)
}

func (s *DispatchService) VerifyPickupOTP(ctx context.Context, userID, orderID, otpCode string) (*model.Order, error) {
	return s.verifyOrderOTP(ctx, userID, orderID, domain.OTPPurposePickup, otpCode, domain.DeliveryStatusPickupVerified, "Pickup OTP verified")
}

func (s *DispatchService) VerifyDeliveryOTP(ctx context.Context, userID, orderID, otpCode string) (*model.Order, error) {
	return s.verifyOrderOTP(ctx, userID, orderID, domain.OTPPurposeDelivery, otpCode, domain.DeliveryStatusDeliveryVerified, "Delivery OTP verified")
}

func (s *DispatchService) ResendPickupOTP(ctx context.Context, orderID string) (map[string]any, error) {
	return s.issueOrderOTP(ctx, orderID, domain.OTPPurposePickup)
}

func (s *DispatchService) ResendDeliveryOTP(ctx context.Context, orderID string) (map[string]any, error) {
	return s.issueOrderOTP(ctx, orderID, domain.OTPPurposeDelivery)
}

func (s *DispatchService) GetOrderTracking(ctx context.Context, orderID string) (map[string]any, error) {
	order, _, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	var assignment *model.DeliveryAssignment
	for _, candidate := range s.repo.Assignments {
		if candidate.OrderID == orderID {
			clone := *candidate
			assignment = &clone
			break
		}
	}
	if assignment == nil {
		return nil, apperrors.NotFound("ASSIGNMENT_NOT_FOUND", "assignment not found")
	}
	location, _ := s.repo.GetLatestLocation(ctx, assignment.RiderID)
	return map[string]any{"order": order, "assignment": assignment, "latest_location": location, "timeline": s.repo.ListOrderHistory(ctx, orderID)}, nil
}

func (s *DispatchService) GetRouteHistory(ctx context.Context, userID string) ([]model.LocationLog, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.GetRouteHistory(ctx, rider.ID), nil
}

func (s *DispatchService) transitionOrder(ctx context.Context, userID, orderID, nextStatus, comment string) (*model.Order, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	assignment, err := s.repo.FindAssignmentByOrderAndRider(ctx, orderID, rider.ID)
	if err != nil {
		return nil, apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	order, _, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	if !canTransition(order.Status, nextStatus) {
		return nil, apperrors.Conflict("INVALID_ORDER_TRANSITION", fmt.Sprintf("cannot move order from %s to %s", order.Status, nextStatus))
	}
	updated, err := s.repo.UpdateOrderStatus(ctx, orderID, nextStatus)
	if err != nil {
		return nil, apperrors.New(500, "ORDER_STATUS_UPDATE_FAILED", "failed to update order status")
	}
	s.repo.AddOrderHistory(ctx, model.DeliveryStatusHistory{OrderID: orderID, AssignmentID: assignment.ID, Status: nextStatus, ActorID: userID, ActorRole: domain.RoleRider, Comment: comment})
	return updated, nil
}

func (s *DispatchService) verifyOrderOTP(ctx context.Context, userID, orderID, purpose, code, nextStatus, comment string) (*model.Order, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	assignment, err := s.repo.FindAssignmentByOrderAndRider(ctx, orderID, rider.ID)
	if err != nil {
		return nil, apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	otp, err := s.repo.GetLatestOTP(ctx, purpose, orderID)
	if err != nil {
		return nil, apperrors.NotFound("OTP_NOT_FOUND", "otp not found")
	}
	if time.Now().UTC().After(otp.ExpiresAt) {
		return nil, apperrors.New(400, "OTP_EXPIRED", "otp expired")
	}
	if otp.Code != code {
		return nil, apperrors.New(400, "INVALID_OTP", "invalid otp")
	}
	_, _ = s.repo.UpdateOTP(ctx, purpose, orderID, func(item *model.OTP) {
		now := time.Now().UTC()
		item.Status = domain.OTPStatusVerified
		item.VerifiedAt = &now
	})
	updated, err := s.repo.UpdateOrderStatus(ctx, orderID, nextStatus)
	if err != nil {
		return nil, apperrors.New(500, "ORDER_STATUS_UPDATE_FAILED", "failed to update order status")
	}
	s.repo.AddOrderHistory(ctx, model.DeliveryStatusHistory{OrderID: orderID, AssignmentID: assignment.ID, Status: nextStatus, ActorID: userID, ActorRole: domain.RoleRider, Comment: comment})
	return updated, nil
}

func (s *DispatchService) issueOrderOTP(ctx context.Context, orderID, purpose string) (map[string]any, error) {
	latest, err := s.repo.GetLatestOTP(ctx, purpose, orderID)
	if err == nil && latest.ResendCount >= s.repo.GetSystemConfig(ctx).OTPMaxResends {
		return nil, apperrors.New(429, "OTP_RESEND_LIMIT_REACHED", "otp resend limit reached")
	}
	code := generateOTPCode()
	saved := s.repo.SaveOTP(ctx, model.OTP{OrderID: orderID, Purpose: purpose, Code: code, Status: domain.OTPStatusPending, ExpiresAt: time.Now().UTC().Add(s.cfg.Auth.OTPExpiry)})
	if err == nil {
		_, _ = s.repo.UpdateOTP(ctx, purpose, orderID, func(item *model.OTP) { item.ResendCount = latest.ResendCount + 1 })
	}
	response := map[string]any{"order_id": orderID, "purpose": purpose, "expires_in_seconds": int(s.cfg.Auth.OTPExpiry.Seconds())}
	if s.cfg.Auth.ExposeOTPInDevMode {
		response["otp"] = saved.Code
	}
	return response, nil
}

func (s *DispatchService) finalizeEarnings(ctx context.Context, riderID string, order *model.Order) {
	incentive := 0.0
	if order.DistanceKM > 5 {
		incentive = 12
	}
	net := order.BasePayout + order.DistancePayout + order.WaitingCharges + order.SurgeBonus + order.TipAmount + incentive
	earning := s.repo.AddEarning(ctx, riderID, model.RiderEarning{OrderID: order.ID, BasePayout: order.BasePayout, DistancePayout: order.DistancePayout, WaitingCharges: order.WaitingCharges, SurgeBonus: order.SurgeBonus, TipAmount: order.TipAmount, IncentiveAmount: incentive, NetEarning: net})
	s.repo.AddWalletTransaction(ctx, riderID, model.WalletTransaction{Type: domain.WalletTxnTypeCredit, Amount: earning.NetEarning, ReferenceID: order.ID, ReferenceType: "ORDER", Description: "Wallet credit after successful delivery"})
}

func canTransition(current, next string) bool {
	transitions := map[string][]string{
		domain.DeliveryStatusAssigned:          {domain.DeliveryStatusAccepted},
		domain.DeliveryStatusAccepted:          {domain.DeliveryStatusReachedRestaurant, domain.DeliveryStatusFailed},
		domain.DeliveryStatusReachedRestaurant: {domain.DeliveryStatusPickupVerified, domain.DeliveryStatusPickedUp, domain.DeliveryStatusFailed},
		domain.DeliveryStatusPickupVerified:    {domain.DeliveryStatusPickedUp},
		domain.DeliveryStatusPickedUp:          {domain.DeliveryStatusOnTheWay},
		domain.DeliveryStatusOnTheWay:          {domain.DeliveryStatusReachedCustomer, domain.DeliveryStatusFailed},
		domain.DeliveryStatusReachedCustomer:   {domain.DeliveryStatusDeliveryVerified, domain.DeliveryStatusDelivered, domain.DeliveryStatusFailed},
		domain.DeliveryStatusDeliveryVerified:  {domain.DeliveryStatusDelivered},
	}
	for _, status := range transitions[current] {
		if status == next {
			return true
		}
	}
	return false
}
func (s *FinanceService) GetSlice(ctx context.Context, riderID string, since time.Time) EarningsSlice {
	records := s.repo.ListEarnings(ctx, riderID)
	filtered := make([]model.RiderEarning, 0)
	total := 0.0
	for _, item := range records {
		if item.CreatedAt.After(since) || item.CreatedAt.Equal(since) {
			filtered = append(filtered, item)
			total += item.NetEarning
		}
	}
	return EarningsSlice{Total: total, Count: len(filtered), Records: filtered}
}

func (s *FinanceService) GetToday(ctx context.Context, userID string) (EarningsSlice, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return EarningsSlice{}, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.GetSlice(ctx, rider.ID, time.Now().UTC().Truncate(24*time.Hour)), nil
}

func (s *FinanceService) GetWeekly(ctx context.Context, userID string) (EarningsSlice, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return EarningsSlice{}, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.GetSlice(ctx, rider.ID, time.Now().UTC().AddDate(0, 0, -7)), nil
}

func (s *FinanceService) GetMonthly(ctx context.Context, userID string) (EarningsSlice, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return EarningsSlice{}, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.GetSlice(ctx, rider.ID, time.Now().UTC().AddDate(0, -1, 0)), nil
}

func (s *FinanceService) GetSummary(ctx context.Context, userID string) (map[string]any, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	wallet, _ := s.repo.GetWallet(ctx, rider.ID)
	weekly := s.GetSlice(ctx, rider.ID, time.Now().UTC().AddDate(0, 0, -7))
	monthly := s.GetSlice(ctx, rider.ID, time.Now().UTC().AddDate(0, -1, 0))
	return map[string]any{"wallet": wallet, "weekly_total": weekly.Total, "monthly_total": monthly.Total, "acceptance_rate": rider.AcceptanceRate, "completion_rate": rider.CompletionRate}, nil
}

func (s *FinanceService) GetHistory(ctx context.Context, userID string) ([]model.RiderEarning, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.ListEarnings(ctx, rider.ID), nil
}

func (s *FinanceService) GetIncentives(ctx context.Context, userID string) ([]model.RiderEarning, error) {
	return s.getBonuses(ctx, userID)
}

func (s *FinanceService) GetBonusHistory(ctx context.Context, userID string) ([]model.RiderEarning, error) {
	return s.getBonuses(ctx, userID)
}

func (s *FinanceService) getBonuses(ctx context.Context, userID string) ([]model.RiderEarning, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	items := make([]model.RiderEarning, 0)
	for _, item := range s.repo.ListEarnings(ctx, rider.ID) {
		if item.IncentiveAmount > 0 {
			items = append(items, item)
		}
	}
	return items, nil
}

func (s *FinanceService) GetWallet(ctx context.Context, userID string) (*model.Wallet, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.GetWallet(ctx, rider.ID)
}

func (s *FinanceService) GetWalletTransactions(ctx context.Context, userID string) ([]model.WalletTransaction, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.ListWalletTransactions(ctx, rider.ID), nil
}

func (s *FinanceService) GetPayouts(ctx context.Context, userID string) ([]model.PayoutRequest, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.ListPayoutRequests(ctx, rider.ID), nil
}

func (s *FinanceService) GetPayout(ctx context.Context, userID, payoutID string) (*model.PayoutRequest, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.GetPayoutRequest(ctx, rider.ID, payoutID)
}

func (s *FinanceService) GetBankAccount(ctx context.Context, userID string) (*model.RiderBankAccount, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.GetBankAccount(ctx, rider.ID)
}

func (s *FinanceService) RequestPayout(ctx context.Context, userID string, amount float64) (*model.PayoutRequest, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	wallet, err := s.repo.GetWallet(ctx, rider.ID)
	if err != nil {
		return nil, apperrors.NotFound("WALLET_NOT_FOUND", "wallet not found")
	}
	if amount < s.cfg.Wallet.MinimumPayoutAmount {
		return nil, apperrors.New(400, "PAYOUT_BELOW_MINIMUM", "requested payout is below minimum threshold")
	}
	if wallet.Balance < amount {
		return nil, apperrors.Conflict("INSUFFICIENT_BALANCE", "insufficient wallet balance")
	}
	payout := s.repo.CreatePayoutRequest(ctx, rider.ID, model.PayoutRequest{Amount: amount, Status: domain.PayoutStatusPending})
	return &payout, nil
}

func (s *FeedbackService) GetSummary(ctx context.Context, userID string) (map[string]any, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	reviews := s.repo.ListRatings(ctx, rider.ID)
	customerCount := 0
	restaurantCount := 0
	for _, review := range reviews {
		if review.Source == domain.RatingSourceCustomer {
			customerCount++
		} else {
			restaurantCount++
		}
	}
	return map[string]any{"average_rating": rider.AvgRating, "rating_count": rider.RatingCount, "customer_reviews": customerCount, "restaurant_reviews": restaurantCount}, nil
}

func (s *FeedbackService) GetReviews(ctx context.Context, userID string) ([]model.RatingReview, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.ListRatings(ctx, rider.ID), nil
}

func (s *FeedbackService) GetPerformanceScore(ctx context.Context, userID string) (map[string]any, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	score := math.Round(((rider.AcceptanceRate*0.25)+(rider.CompletionRate*0.35)+((rider.AvgRating/5)*100*0.25)+((100-rider.CancellationRate)*0.15))*100) / 100
	return map[string]any{"score": score, "acceptance_rate": rider.AcceptanceRate, "completion_rate": rider.CompletionRate, "rating_average": rider.AvgRating}, nil
}

func (s *NotificationService) List(ctx context.Context, userID string) ([]model.Notification, error) {
	return s.repo.ListNotifications(ctx, userID), nil
}
func (s *NotificationService) MarkRead(ctx context.Context, userID, notificationID string) error {
	return s.repo.MarkNotificationRead(ctx, userID, notificationID)
}
func (s *NotificationService) MarkAllRead(ctx context.Context, userID string) {
	s.repo.MarkAllNotificationsRead(ctx, userID)
}
func (s *NotificationService) RegisterDeviceToken(ctx context.Context, userID string, req dto.DeviceTokenRequest) model.DeviceToken {
	return s.repo.UpsertDeviceToken(ctx, userID, model.DeviceToken{DeviceID: req.DeviceID, Platform: req.Platform, Token: req.Token})
}
func (s *NotificationService) DeleteDeviceToken(ctx context.Context, userID, deviceID string) {
	s.repo.DeleteDeviceToken(ctx, userID, deviceID)
}
func (s *SupportService) CreateTicket(ctx context.Context, userID string, req dto.SupportTicketRequest) (*model.SupportTicket, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	ticket := s.repo.CreateSupportTicket(ctx, model.SupportTicket{RiderID: rider.ID, Subject: req.Subject, Category: req.Category, Priority: req.Priority, Status: domain.SupportTicketStatusOpen, OrderID: req.OrderID, Description: req.Description})
	s.repo.AddSupportMessage(ctx, model.SupportTicketMessage{TicketID: ticket.ID, ActorID: userID, ActorRole: domain.RoleRider, Message: req.Description})
	return &ticket, nil
}

func (s *SupportService) ListTickets(ctx context.Context, userID string) ([]model.SupportTicket, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.ListSupportTickets(ctx, rider.ID), nil
}

func (s *SupportService) GetTicket(ctx context.Context, ticketID string) (map[string]any, error) {
	ticket, messages, err := s.repo.GetSupportTicket(ctx, ticketID)
	if err != nil {
		return nil, apperrors.NotFound("TICKET_NOT_FOUND", "ticket not found")
	}
	return map[string]any{"ticket": ticket, "messages": messages}, nil
}

func (s *SupportService) ReplyTicket(ctx context.Context, ticketID, userID string, req dto.SupportReplyRequest) (map[string]any, error) {
	message := s.repo.AddSupportMessage(ctx, model.SupportTicketMessage{TicketID: ticketID, ActorID: userID, ActorRole: domain.RoleRider, Message: req.Message})
	ticket, messages, _ := s.repo.GetSupportTicket(ctx, ticketID)
	return map[string]any{"ticket": ticket, "latest_message": message, "messages": messages}, nil
}

func (s *AdminService) ListRiders(ctx context.Context) ([]map[string]any, error) {
	items := s.repo.ListRiders(ctx)
	output := make([]map[string]any, 0, len(items))
	for _, item := range items {
		output = append(output, map[string]any{"user": item.User, "rider": item.Rider})
	}
	return output, nil
}

func (s *AdminService) GetRider(ctx context.Context, riderID string) (map[string]any, error) {
	rider, user, err := s.repo.GetRiderByID(ctx, riderID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return map[string]any{"user": user, "rider": rider, "vehicle": s.repo.RiderVehicles[rider.ID], "bank_account": s.repo.RiderBankAccounts[rider.ID], "documents": s.repo.ListDocuments(ctx, rider.ID)}, nil
}

func (s *AdminService) CreateRider(ctx context.Context, req dto.CreateRiderRequest) (map[string]any, error) {
	hash, err := auth.HashPassword(req.Password, s.cfg.Auth.PasswordHashCost)
	if err != nil {
		return nil, apperrors.New(500, "PASSWORD_HASH_FAILED", "failed to hash password")
	}
	user, rider := s.repo.CreateRider(ctx, model.User{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Phone: req.Phone, PasswordHash: hash, Status: domain.UserStatusActive, Roles: []string{domain.RoleRider}, PrimaryRole: domain.RoleRider, EmailVerified: true, PhoneVerified: true}, model.Rider{Status: domain.UserStatusActive, AvailabilityStatus: domain.RiderAvailabilityOffline, KYCStatus: "PENDING", ApprovalStatus: "PENDING"}, req.RestaurantID)
	return map[string]any{"user": user, "rider": rider}, nil
}

func (s *AdminService) UpdateRiderStatus(ctx context.Context, riderID, status string) (*model.Rider, error) {
	return s.repo.UpdateRiderStatus(ctx, riderID, status)
}

func (s *AdminService) ListUnassignedOrders(ctx context.Context) ([]model.Order, error) {
	return s.repo.ListUnassignedOrders(ctx), nil
}
func (s *AdminService) AssignOrder(ctx context.Context, orderID, riderID, assignedBy string) (*model.DeliveryAssignment, error) {
	return s.repo.AssignOrder(ctx, orderID, riderID, assignedBy)
}
func (s *AdminService) ReassignOrder(ctx context.Context, orderID, riderID, assignedBy string) (*model.DeliveryAssignment, error) {
	return s.repo.ReassignOrder(ctx, orderID, riderID, assignedBy)
}
func (s *AdminService) ListLiveOrders(ctx context.Context) ([]model.Order, error) {
	return s.repo.ListLiveOrders(ctx), nil
}

func (s *AdminService) ListLiveRiderStatus(ctx context.Context) ([]map[string]any, error) {
	items := s.repo.ListRiders(ctx)
	output := make([]map[string]any, 0, len(items))
	for _, item := range items {
		activeShift, _ := s.repo.GetActiveShift(ctx, item.Rider.ID)
		activeAssignment, _ := s.repo.GetActiveAssignmentForRider(ctx, item.Rider.ID)
		output = append(output, map[string]any{"user": item.User, "rider": item.Rider, "active_shift": activeShift, "active_assignment": activeAssignment})
	}
	return output, nil
}

func (s *AdminService) GetAnalytics(ctx context.Context) (map[string]any, error) {
	riders := s.repo.ListRiders(ctx)
	totals := map[string]float64{"acceptance_rate": 0, "completion_rate": 0, "cancellation_rate": 0, "rating_average": 0}
	for _, item := range riders {
		totals["acceptance_rate"] += item.Rider.AcceptanceRate
		totals["completion_rate"] += item.Rider.CompletionRate
		totals["cancellation_rate"] += item.Rider.CancellationRate
		totals["rating_average"] += item.Rider.AvgRating
	}
	count := math.Max(float64(len(riders)), 1)
	return map[string]any{"rider_count": len(riders), "avg_acceptance_rate": math.Round((totals["acceptance_rate"]/count)*100) / 100, "avg_completion_rate": math.Round((totals["completion_rate"]/count)*100) / 100, "avg_cancellation_rate": math.Round((totals["cancellation_rate"]/count)*100) / 100, "avg_rating": math.Round((totals["rating_average"]/count)*100) / 100, "live_orders": len(s.repo.ListLiveOrders(ctx))}, nil
}

func (s *AdminService) GetSystemConfig(ctx context.Context) repository.SystemConfig {
	return s.repo.GetSystemConfig(ctx)
}

func (s *AdminService) UpdateSystemConfig(ctx context.Context, req dto.UpdateConfigRequest) repository.SystemConfig {
	return s.repo.UpdateSystemConfig(ctx, repository.SystemConfig{OrderAcceptTimeoutSeconds: req.OrderAcceptTimeoutSeconds, PickupOTPRequired: req.PickupOTPRequired, DeliveryOTPRequired: req.DeliveryOTPRequired, OTPMaxRetries: req.OTPMaxRetries, OTPMaxResends: req.OTPMaxResends, RiderMaxActiveOrders: req.RiderMaxActiveOrders, SurgeMultiplierDefault: req.SurgeMultiplierDefault, MinimumPayoutAmount: req.MinimumPayoutAmount})
}

func (s *AdminService) ListPayoutRequests(ctx context.Context) ([]model.PayoutRequest, error) {
	return s.repo.ListAllPayoutRequests(ctx), nil
}

func (s *AdminService) ApprovePayout(ctx context.Context, payoutID, reviewer string) (*model.PayoutRequest, error) {
	return s.repo.ReviewPayout(ctx, payoutID, domain.PayoutStatusApproved, reviewer, "")
}

func (s *AdminService) RejectPayout(ctx context.Context, payoutID, reviewer, reason string) (*model.PayoutRequest, error) {
	return s.repo.ReviewPayout(ctx, payoutID, domain.PayoutStatusRejected, reviewer, reason)
}
func (s *AdminService) UpdateRider(ctx context.Context, riderID string, req dto.UpdateRiderProfileRequest) (map[string]any, error) {
	rider, user, err := s.repo.GetRiderByID(ctx, riderID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	updatedUser, err := s.repo.UpdateUserProfile(ctx, user.ID, req.FirstName, req.LastName, req.Email, req.Phone)
	if err != nil {
		return nil, apperrors.New(500, "RIDER_UPDATE_FAILED", "failed to update rider profile")
	}
	return map[string]any{"user": updatedUser, "rider": rider}, nil
}
func (s *DispatchService) UpdateLocation(ctx context.Context, userID string, req dto.LocationUpdateRequest) (*model.RiderLocationLatest, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	location := s.repo.SaveLatestLocation(ctx, model.RiderLocationLatest{RiderID: rider.ID, OrderID: req.OrderID, Latitude: req.Latitude, Longitude: req.Longitude, AccuracyMeters: req.AccuracyMeters, SpeedKPH: req.SpeedKPH, HeadingDegrees: req.HeadingDegrees, BatteryLevel: req.BatteryLevel, Source: req.Source})
	return &location, nil
}

func (s *DispatchService) BulkUpdateLocation(ctx context.Context, userID string, req dto.BulkLocationUpdateRequest) ([]model.RiderLocationLatest, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	points := make([]model.RiderLocationLatest, 0, len(req.Points))
	for _, point := range req.Points {
		points = append(points, model.RiderLocationLatest{OrderID: point.OrderID, Latitude: point.Latitude, Longitude: point.Longitude, AccuracyMeters: point.AccuracyMeters, SpeedKPH: point.SpeedKPH, HeadingDegrees: point.HeadingDegrees, BatteryLevel: point.BatteryLevel, Source: point.Source})
	}
	return s.repo.BulkSaveLocations(ctx, rider.ID, points), nil
}

func (s *DispatchService) GetLatestLocation(ctx context.Context, userID string) (*model.RiderLocationLatest, error) {
	rider, err := s.repo.GetRiderByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("RIDER_NOT_FOUND", "rider not found")
	}
	return s.repo.GetLatestLocation(ctx, rider.ID)
}
