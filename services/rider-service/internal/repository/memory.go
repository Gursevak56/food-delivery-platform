package repository

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/domain"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/model"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/utils"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/auth"
)

type SystemConfig struct {
	OrderAcceptTimeoutSeconds int     `json:"order_accept_timeout_seconds"`
	PickupOTPRequired         bool    `json:"pickup_otp_required"`
	DeliveryOTPRequired       bool    `json:"delivery_otp_required"`
	OTPMaxRetries             int     `json:"otp_max_retries"`
	OTPMaxResends             int     `json:"otp_max_resends"`
	RiderMaxActiveOrders      int     `json:"rider_max_active_orders"`
	SurgeMultiplierDefault    float64 `json:"surge_multiplier_default"`
	MinimumPayoutAmount       float64 `json:"minimum_payout_amount"`
}

type Repository struct {
	mu sync.RWMutex

	Users                map[string]*model.User
	Restaurants          map[string]*model.Restaurant
	Riders               map[string]*model.Rider
	RiderDocuments       map[string][]model.RiderDocument
	RiderVehicles        map[string]*model.RiderVehicle
	RiderBankAccounts    map[string]*model.RiderBankAccount
	RiderRestaurantLinks map[string][]model.RiderRestaurantAssignment
	Orders               map[string]*model.Order
	OrderItems           map[string][]model.OrderItem
	Assignments          map[string]*model.DeliveryAssignment
	OrderStatusHistory   map[string][]model.DeliveryStatusHistory
	OTPs                 map[string][]model.OTP
	Shifts               map[string][]model.RiderShift
	Breaks               map[string][]model.RiderBreak
	LatestLocations      map[string]*model.RiderLocationLatest
	LocationLogs         map[string][]model.LocationLog
	Wallets              map[string]*model.Wallet
	WalletTransactions   map[string][]model.WalletTransaction
	Earnings             map[string][]model.RiderEarning
	PayoutRequests       map[string][]model.PayoutRequest
	Ratings              map[string][]model.RatingReview
	Notifications        map[string][]model.Notification
	DeviceTokens         map[string][]model.DeviceToken
	SupportTickets       map[string]*model.SupportTicket
	SupportMessages      map[string][]model.SupportTicketMessage
	Sessions             map[string]*model.Session
	AuditLogs            []model.AuditLog
	Config               SystemConfig
}

func NewMemoryRepository(hashCost int) *Repository {
	repo := &Repository{
		Users:                map[string]*model.User{},
		Restaurants:          map[string]*model.Restaurant{},
		Riders:               map[string]*model.Rider{},
		RiderDocuments:       map[string][]model.RiderDocument{},
		RiderVehicles:        map[string]*model.RiderVehicle{},
		RiderBankAccounts:    map[string]*model.RiderBankAccount{},
		RiderRestaurantLinks: map[string][]model.RiderRestaurantAssignment{},
		Orders:               map[string]*model.Order{},
		OrderItems:           map[string][]model.OrderItem{},
		Assignments:          map[string]*model.DeliveryAssignment{},
		OrderStatusHistory:   map[string][]model.DeliveryStatusHistory{},
		OTPs:                 map[string][]model.OTP{},
		Shifts:               map[string][]model.RiderShift{},
		Breaks:               map[string][]model.RiderBreak{},
		LatestLocations:      map[string]*model.RiderLocationLatest{},
		LocationLogs:         map[string][]model.LocationLog{},
		Wallets:              map[string]*model.Wallet{},
		WalletTransactions:   map[string][]model.WalletTransaction{},
		Earnings:             map[string][]model.RiderEarning{},
		PayoutRequests:       map[string][]model.PayoutRequest{},
		Ratings:              map[string][]model.RatingReview{},
		Notifications:        map[string][]model.Notification{},
		DeviceTokens:         map[string][]model.DeviceToken{},
		SupportTickets:       map[string]*model.SupportTicket{},
		SupportMessages:      map[string][]model.SupportTicketMessage{},
		Sessions:             map[string]*model.Session{},
		Config: SystemConfig{
			OrderAcceptTimeoutSeconds: 45,
			PickupOTPRequired:         true,
			DeliveryOTPRequired:       true,
			OTPMaxRetries:             5,
			OTPMaxResends:             3,
			RiderMaxActiveOrders:      1,
			SurgeMultiplierDefault:    1.0,
			MinimumPayoutAmount:       500,
		},
	}
	repo.seed(hashCost)
	return repo
}

func (r *Repository) seed(hashCost int) {
	now := time.Now().UTC()
	restaurant := &model.Restaurant{ID: "rest_mg_001", Name: "Mangaale Downtown", Code: "MG-DT", Status: domain.UserStatusActive, City: "Pune", Area: "Baner", CreatedAt: now, UpdatedAt: now}
	r.Restaurants[restaurant.ID] = restaurant

	adminHash, _ := auth.HashPassword("Admin@123", hashCost)
	riderHash, _ := auth.HashPassword("Rider@123", hashCost)
	dispatcherHash, _ := auth.HashPassword("Dispatch@123", hashCost)

	admin := &model.User{ID: "usr_admin_001", FirstName: "Platform", LastName: "Admin", Email: "admin@rider.local", Phone: "+910000000001", PasswordHash: adminHash, Status: domain.UserStatusActive, Roles: []string{domain.RoleSuperAdmin}, PrimaryRole: domain.RoleSuperAdmin, EmailVerified: true, PhoneVerified: true, CreatedAt: now, UpdatedAt: now}
	dispatcher := &model.User{ID: "usr_dispatch_001", FirstName: "Dispatch", LastName: "Lead", Email: "dispatch@rider.local", Phone: "+910000000002", PasswordHash: dispatcherHash, Status: domain.UserStatusActive, Roles: []string{domain.RoleDispatcherManager}, PrimaryRole: domain.RoleDispatcherManager, EmailVerified: true, PhoneVerified: true, CreatedAt: now, UpdatedAt: now}
	riderUser := &model.User{ID: "usr_rider_001", FirstName: "Ravi", LastName: "Kumar", Email: "rider@rider.local", Phone: "+910000000010", PasswordHash: riderHash, Status: domain.UserStatusActive, Roles: []string{domain.RoleRider}, PrimaryRole: domain.RoleRider, PhotoURL: "https://cdn.example.com/riders/ravi.jpg", EmailVerified: true, PhoneVerified: true, CreatedAt: now, UpdatedAt: now}

	r.Users[admin.ID] = admin
	r.Users[dispatcher.ID] = dispatcher
	r.Users[riderUser.ID] = riderUser

	rider := &model.Rider{ID: "rdr_001", UserID: riderUser.ID, Status: domain.UserStatusActive, AvailabilityStatus: domain.RiderAvailabilityOffline, KYCStatus: "APPROVED", ApprovalStatus: "APPROVED", VehicleType: "BIKE", VehicleNumber: "MH12AB1234", AvgRating: 4.8, RatingCount: 128, AcceptanceRate: 92.3, CompletionRate: 95.6, CancellationRate: 1.4, RejectionRate: 4.1, CreatedAt: now, UpdatedAt: now}
	r.Riders[rider.ID] = rider
	r.RiderVehicles[rider.ID] = &model.RiderVehicle{ID: "veh_001", RiderID: rider.ID, VehicleType: "BIKE", RegistrationNo: "MH12AB1234", Color: "Red", Capacity: 25, InsuranceExpiry: now.AddDate(0, 6, 0), UpdatedAt: now}
	r.RiderBankAccounts[rider.ID] = &model.RiderBankAccount{ID: "bank_001", RiderID: rider.ID, AccountHolder: "Ravi Kumar", BankName: "State Bank of India", AccountNumber: "XXXXXX9012", IFSCCode: "SBIN0000456", IsVerified: true, UpdatedAt: now}
	r.RiderDocuments[rider.ID] = []model.RiderDocument{
		{ID: "doc_001", RiderID: rider.ID, DocumentType: "DRIVING_LICENSE", DocumentURL: "https://cdn.example.com/docs/dl.jpg", Status: "APPROVED", UploadedAt: now},
		{ID: "doc_002", RiderID: rider.ID, DocumentType: "AADHAAR", DocumentURL: "https://cdn.example.com/docs/aadhaar.jpg", Status: "APPROVED", UploadedAt: now},
	}
	r.RiderRestaurantLinks[rider.ID] = []model.RiderRestaurantAssignment{{ID: "rr_001", RiderID: rider.ID, RestaurantID: restaurant.ID, IsPrimary: true, IsActive: true, AssignedBy: admin.ID, AssignedAt: now}}
	r.Wallets[rider.ID] = &model.Wallet{ID: "wal_001", RiderID: rider.ID, Balance: 1240, HoldBalance: 0, Currency: "INR", UpdatedAt: now}
	r.Ratings[rider.ID] = []model.RatingReview{
		{ID: "rev_001", RiderID: rider.ID, OrderID: "ord_001", Source: domain.RatingSourceCustomer, Rating: 5, Review: "Very polite and fast delivery", CreatedAt: now.Add(-24 * time.Hour)},
		{ID: "rev_002", RiderID: rider.ID, OrderID: "ord_002", Source: domain.RatingSourceRestaurant, Rating: 4, Review: "Reached pickup station quickly", CreatedAt: now.Add(-48 * time.Hour)},
	}
	r.Notifications[riderUser.ID] = []model.Notification{{ID: "ntf_001", UserID: riderUser.ID, Type: "ORDER_REQUEST", Title: "New delivery request", Body: "Order ORD-1001 is waiting for your response", Status: domain.NotificationStatusUnread, Channel: "PUSH", OrderID: "ord_001", CreatedAt: now.Add(-2 * time.Minute)}}

	order := &model.Order{ID: "ord_001", OrderNumber: "ORD-1001", RestaurantID: restaurant.ID, CustomerName: "Ananya Singh", CustomerPhone: "+919876543210", DeliveryAddress: "Baner Road, Pune", DeliveryLatitude: 18.559, DeliveryLongitude: 73.786, Status: domain.DeliveryStatusAssigned, DistanceKM: 5.8, BasePayout: 40, DistancePayout: 26, WaitingCharges: 0, SurgeBonus: 10, TipAmount: 15, TotalAmount: 780, PickupOTPRequired: true, DeliveryOTPRequired: true, CreatedAt: now.Add(-10 * time.Minute), UpdatedAt: now}
	r.Orders[order.ID] = order
	r.OrderItems[order.ID] = []model.OrderItem{{ID: "itm_001", OrderID: order.ID, Name: "Paneer Bowl", Quantity: 2, UnitPrice: 220}, {ID: "itm_002", OrderID: order.ID, Name: "Fresh Lime Soda", Quantity: 1, UnitPrice: 120}}
	assignment := &model.DeliveryAssignment{ID: "asg_001", OrderID: order.ID, RiderID: rider.ID, RestaurantID: restaurant.ID, Status: domain.AssignmentStatusOffered, AssignedBy: dispatcher.ID, AssignedAt: now.Add(-3 * time.Minute), DecisionDeadlineAt: now.Add(42 * time.Second)}
	r.Assignments[assignment.ID] = assignment
	r.OrderStatusHistory[order.ID] = []model.DeliveryStatusHistory{{ID: "hst_001", OrderID: order.ID, AssignmentID: assignment.ID, Status: domain.DeliveryStatusAssigned, ActorID: dispatcher.ID, ActorRole: domain.RoleDispatcherManager, Comment: "Order assigned to rider queue", CreatedAt: now.Add(-3 * time.Minute)}}

	secondaryOrder := &model.Order{ID: "ord_002", OrderNumber: "ORD-1002", RestaurantID: restaurant.ID, CustomerName: "Karthik Rao", CustomerPhone: "+919988776655", DeliveryAddress: "Aundh, Pune", DeliveryLatitude: 18.562, DeliveryLongitude: 73.807, Status: domain.DeliveryStatusDelivered, DistanceKM: 3.4, BasePayout: 40, DistancePayout: 18, WaitingCharges: 5, SurgeBonus: 0, TipAmount: 20, TotalAmount: 540, PickupOTPRequired: false, DeliveryOTPRequired: true, CreatedAt: now.Add(-30 * time.Hour), UpdatedAt: now.Add(-29 * time.Hour), DeliveredAt: ptrTime(now.Add(-29 * time.Hour))}
	r.Orders[secondaryOrder.ID] = secondaryOrder
	r.Assignments["asg_002"] = &model.DeliveryAssignment{ID: "asg_002", OrderID: secondaryOrder.ID, RiderID: rider.ID, RestaurantID: restaurant.ID, Status: domain.AssignmentStatusCompleted, AssignedBy: dispatcher.ID, AssignedAt: now.Add(-31 * time.Hour), AcceptedAt: ptrTime(now.Add(-30*time.Hour + 5*time.Minute)), DecisionDeadlineAt: now.Add(-30 * time.Hour)}
	r.Earnings[rider.ID] = []model.RiderEarning{{ID: "earn_001", RiderID: rider.ID, OrderID: secondaryOrder.ID, BasePayout: 40, DistancePayout: 18, WaitingCharges: 5, SurgeBonus: 0, TipAmount: 20, IncentiveAmount: 10, PenaltyAmount: 0, CancellationCompensation: 0, NetEarning: 93, CreatedAt: now.Add(-29 * time.Hour)}}
	r.WalletTransactions[rider.ID] = []model.WalletTransaction{{ID: "wtx_001", WalletID: "wal_001", RiderID: rider.ID, Type: domain.WalletTxnTypeCredit, Amount: 93, ReferenceID: secondaryOrder.ID, ReferenceType: "ORDER", Description: "Wallet credit for delivery ORD-1002", CreatedAt: now.Add(-29 * time.Hour)}}
	r.PayoutRequests[rider.ID] = []model.PayoutRequest{{ID: "pay_001", RiderID: rider.ID, WalletID: "wal_001", Amount: 750, Status: domain.PayoutStatusPaid, RequestedAt: now.Add(-7 * 24 * time.Hour), ReviewedAt: ptrTime(now.Add(-6 * 24 * time.Hour)), ReviewedBy: admin.ID}}
}

func ptrTime(value time.Time) *time.Time { return &value }

func (r *Repository) FindUserByLogin(_ context.Context, login string) (*model.User, *model.Rider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	login = strings.TrimSpace(strings.ToLower(login))
	for _, user := range r.Users {
		if strings.ToLower(user.Email) == login || strings.ToLower(user.Phone) == login {
			rider := r.getRiderByUserIDLocked(user.ID)
			return cloneUser(user), cloneRider(rider), nil
		}
	}
	return nil, nil, errors.New("user not found")
}

func (r *Repository) GetUserByID(_ context.Context, userID string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.Users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return cloneUser(user), nil
}

func (r *Repository) GetRiderByUserID(_ context.Context, userID string) (*model.Rider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rider := r.getRiderByUserIDLocked(userID)
	if rider == nil {
		return nil, errors.New("rider not found")
	}
	return cloneRider(rider), nil
}

func (r *Repository) GetRiderProfile(_ context.Context, userID string) (*model.User, *model.Rider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.Users[userID]
	if !ok {
		return nil, nil, errors.New("user not found")
	}
	rider := r.getRiderByUserIDLocked(userID)
	if rider == nil {
		return nil, nil, errors.New("rider not found")
	}
	return cloneUser(user), cloneRider(rider), nil
}

func (r *Repository) GetRiderByID(_ context.Context, riderID string) (*model.Rider, *model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rider, ok := r.Riders[riderID]
	if !ok {
		return nil, nil, errors.New("rider not found")
	}
	return cloneRider(rider), cloneUser(r.Users[rider.UserID]), nil
}

func (r *Repository) UpdateUserProfile(_ context.Context, userID, firstName, lastName, email, phone string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.Users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email
	user.Phone = phone
	user.UpdatedAt = time.Now().UTC()
	return cloneUser(user), nil
}

func (r *Repository) UpdatePhoto(_ context.Context, userID, photoURL string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.Users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	user.PhotoURL = photoURL
	user.UpdatedAt = time.Now().UTC()
	return cloneUser(user), nil
}

func (r *Repository) UpdateVehicle(_ context.Context, riderID string, vehicle model.RiderVehicle) (*model.RiderVehicle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	vehicle.ID = coalesceID(vehicle.ID, utils.NewID("veh"))
	vehicle.RiderID = riderID
	vehicle.UpdatedAt = time.Now().UTC()
	r.RiderVehicles[riderID] = &vehicle
	if rider, ok := r.Riders[riderID]; ok {
		rider.VehicleType = vehicle.VehicleType
		rider.VehicleNumber = vehicle.RegistrationNo
		rider.UpdatedAt = vehicle.UpdatedAt
	}
	clone := vehicle
	return &clone, nil
}

func (r *Repository) ListDocuments(_ context.Context, riderID string) []model.RiderDocument {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]model.RiderDocument(nil), r.RiderDocuments[riderID]...)
}

func (r *Repository) ReplaceDocuments(_ context.Context, riderID string, docs []model.RiderDocument) []model.RiderDocument {
	r.mu.Lock()
	defer r.mu.Unlock()
	for index := range docs {
		docs[index].ID = utils.NewID("doc")
		docs[index].RiderID = riderID
		docs[index].Status = "PENDING_REVIEW"
		docs[index].UploadedAt = time.Now().UTC()
	}
	r.RiderDocuments[riderID] = append([]model.RiderDocument(nil), docs...)
	return append([]model.RiderDocument(nil), docs...)
}

func (r *Repository) UpsertBankAccount(_ context.Context, riderID string, account model.RiderBankAccount) *model.RiderBankAccount {
	r.mu.Lock()
	defer r.mu.Unlock()
	account.ID = coalesceID(account.ID, utils.NewID("bank"))
	account.RiderID = riderID
	account.UpdatedAt = time.Now().UTC()
	stored := account
	r.RiderBankAccounts[riderID] = &stored
	clone := stored
	return &clone
}

func (r *Repository) GetBankAccount(_ context.Context, riderID string) (*model.RiderBankAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	account, ok := r.RiderBankAccounts[riderID]
	if !ok {
		return nil, errors.New("bank account not found")
	}
	clone := *account
	return &clone, nil
}

func (r *Repository) SaveSession(_ context.Context, session model.Session) model.Session {
	r.mu.Lock()
	defer r.mu.Unlock()
	session.ID = coalesceID(session.ID, utils.NewID("ses"))
	session.CreatedAt = time.Now().UTC()
	session.LastSeenAt = session.CreatedAt
	stored := session
	r.Sessions[session.RefreshHash] = &stored
	return stored
}

func (r *Repository) GetSession(_ context.Context, refreshHash string) (*model.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.Sessions[refreshHash]
	if !ok {
		return nil, errors.New("session not found")
	}
	clone := *session
	return &clone, nil
}

func (r *Repository) RevokeSession(_ context.Context, refreshHash string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if session, ok := r.Sessions[refreshHash]; ok {
		now := time.Now().UTC()
		session.RevokedAt = &now
		session.LastSeenAt = now
	}
}

func (r *Repository) RevokeAllSessions(_ context.Context, userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	for _, session := range r.Sessions {
		if session.UserID == userID {
			session.RevokedAt = &now
			session.LastSeenAt = now
		}
	}
}

func (r *Repository) SaveOTP(_ context.Context, otp model.OTP) model.OTP {
	r.mu.Lock()
	defer r.mu.Unlock()
	otp.ID = coalesceID(otp.ID, utils.NewID("otp"))
	otp.CreatedAt = time.Now().UTC()
	key := otp.Purpose + ":" + firstNonEmpty(otp.UserID, otp.OrderID)
	r.OTPs[key] = append(r.OTPs[key], otp)
	return otp
}

func (r *Repository) GetLatestOTP(_ context.Context, purpose, targetID string) (*model.OTP, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := purpose + ":" + targetID
	values := r.OTPs[key]
	if len(values) == 0 {
		return nil, errors.New("otp not found")
	}
	otp := values[len(values)-1]
	return &otp, nil
}

func (r *Repository) UpdateOTP(_ context.Context, purpose, targetID string, updater func(*model.OTP)) (*model.OTP, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := purpose + ":" + targetID
	values := r.OTPs[key]
	if len(values) == 0 {
		return nil, errors.New("otp not found")
	}
	updater(&values[len(values)-1])
	r.OTPs[key] = values
	otp := values[len(values)-1]
	return &otp, nil
}

func (r *Repository) UpdatePassword(_ context.Context, userID, passwordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.Users[userID]
	if !ok {
		return errors.New("user not found")
	}
	user.PasswordHash = passwordHash
	user.UpdatedAt = time.Now().UTC()
	return nil
}
func (r *Repository) SetRiderAvailability(_ context.Context, riderID, status string) (*model.Rider, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rider, ok := r.Riders[riderID]
	if !ok {
		return nil, errors.New("rider not found")
	}
	rider.AvailabilityStatus = status
	rider.UpdatedAt = time.Now().UTC()
	return cloneRider(rider), nil
}

func (r *Repository) StartShift(_ context.Context, riderID, startedBy string) model.RiderShift {
	r.mu.Lock()
	defer r.mu.Unlock()
	shift := model.RiderShift{ID: utils.NewID("shf"), RiderID: riderID, Status: domain.ShiftStatusActive, StartedAt: time.Now().UTC(), StartedBy: startedBy}
	r.Shifts[riderID] = append(r.Shifts[riderID], shift)
	if rider, ok := r.Riders[riderID]; ok {
		rider.CurrentShiftID = shift.ID
	}
	return shift
}

func (r *Repository) EndShift(_ context.Context, riderID, endedBy string) (*model.RiderShift, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	shifts := r.Shifts[riderID]
	if len(shifts) == 0 {
		return nil, errors.New("shift not found")
	}
	last := &shifts[len(shifts)-1]
	if last.Status != domain.ShiftStatusActive {
		return nil, errors.New("no active shift")
	}
	now := time.Now().UTC()
	last.Status = domain.ShiftStatusEnded
	last.EndedAt = &now
	last.EndedBy = endedBy
	last.ActiveMinutes = int(now.Sub(last.StartedAt).Minutes())
	r.Shifts[riderID] = shifts
	if rider, ok := r.Riders[riderID]; ok {
		rider.CurrentShiftID = ""
	}
	clone := *last
	return &clone, nil
}

func (r *Repository) GetActiveShift(_ context.Context, riderID string) (*model.RiderShift, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	shifts := r.Shifts[riderID]
	for i := len(shifts) - 1; i >= 0; i-- {
		if shifts[i].Status == domain.ShiftStatusActive {
			shift := shifts[i]
			return &shift, nil
		}
	}
	return nil, errors.New("active shift not found")
}

func (r *Repository) ListShifts(_ context.Context, riderID string) []model.RiderShift {
	r.mu.RLock()
	defer r.mu.RUnlock()
	shifts := append([]model.RiderShift(nil), r.Shifts[riderID]...)
	sort.Slice(shifts, func(i, j int) bool { return shifts[i].StartedAt.After(shifts[j].StartedAt) })
	return shifts
}

func (r *Repository) StartBreak(_ context.Context, riderID, shiftID, reason string) model.RiderBreak {
	r.mu.Lock()
	defer r.mu.Unlock()
	brk := model.RiderBreak{ID: utils.NewID("brk"), RiderID: riderID, ShiftID: shiftID, Reason: reason, Status: domain.BreakStatusActive, StartedAt: time.Now().UTC()}
	r.Breaks[riderID] = append(r.Breaks[riderID], brk)
	if rider, ok := r.Riders[riderID]; ok {
		rider.CurrentBreakID = brk.ID
		rider.AvailabilityStatus = domain.RiderAvailabilityBreak
	}
	return brk
}

func (r *Repository) EndBreak(_ context.Context, riderID string) (*model.RiderBreak, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	breaks := r.Breaks[riderID]
	if len(breaks) == 0 {
		return nil, errors.New("break not found")
	}
	last := &breaks[len(breaks)-1]
	if last.Status != domain.BreakStatusActive {
		return nil, errors.New("no active break")
	}
	now := time.Now().UTC()
	last.Status = domain.BreakStatusEnded
	last.EndedAt = &now
	r.Breaks[riderID] = breaks
	if rider, ok := r.Riders[riderID]; ok {
		rider.CurrentBreakID = ""
		rider.AvailabilityStatus = domain.RiderAvailabilityOnline
	}
	clone := *last
	return &clone, nil
}

func (r *Repository) ListPendingAssignments(_ context.Context, riderID string) []model.DeliveryAssignment {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	list := make([]model.DeliveryAssignment, 0)
	for _, assignment := range r.Assignments {
		if assignment.RiderID == riderID && assignment.Status == domain.AssignmentStatusOffered {
			if now.After(assignment.DecisionDeadlineAt) {
				assignment.Status = domain.AssignmentStatusTimedOut
				if order, ok := r.Orders[assignment.OrderID]; ok {
					order.Status = domain.DeliveryStatusReadyForAssignment
					order.UpdatedAt = now
				}
				continue
			}
			list = append(list, *assignment)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].AssignedAt.After(list[j].AssignedAt) })
	return list
}

func (r *Repository) GetAssignment(_ context.Context, assignmentID string) (*model.DeliveryAssignment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	assignment, ok := r.Assignments[assignmentID]
	if !ok {
		return nil, errors.New("assignment not found")
	}
	clone := *assignment
	return &clone, nil
}

func (r *Repository) AcceptAssignment(_ context.Context, assignmentID string) (*model.DeliveryAssignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	assignment, ok := r.Assignments[assignmentID]
	if !ok {
		return nil, errors.New("assignment not found")
	}
	now := time.Now().UTC()
	assignment.Status = domain.AssignmentStatusAccepted
	assignment.AcceptedAt = &now
	if order, ok := r.Orders[assignment.OrderID]; ok {
		order.Status = domain.DeliveryStatusAccepted
		order.UpdatedAt = now
	}
	if rider, ok := r.Riders[assignment.RiderID]; ok {
		rider.AvailabilityStatus = domain.RiderAvailabilityBusy
		rider.UpdatedAt = now
	}
	clone := *assignment
	return &clone, nil
}

func (r *Repository) RejectAssignment(_ context.Context, assignmentID, reason string) (*model.DeliveryAssignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	assignment, ok := r.Assignments[assignmentID]
	if !ok {
		return nil, errors.New("assignment not found")
	}
	now := time.Now().UTC()
	assignment.Status = domain.AssignmentStatusRejected
	assignment.RejectReason = reason
	assignment.RejectedAt = &now
	if order, ok := r.Orders[assignment.OrderID]; ok {
		order.Status = domain.DeliveryStatusReadyForAssignment
		order.UpdatedAt = now
	}
	clone := *assignment
	return &clone, nil
}

func (r *Repository) GetOrder(_ context.Context, orderID string) (*model.Order, []model.OrderItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.Orders[orderID]
	if !ok {
		return nil, nil, errors.New("order not found")
	}
	clone := *order
	items := append([]model.OrderItem(nil), r.OrderItems[orderID]...)
	return &clone, items, nil
}

func (r *Repository) UpdateOrderStatus(_ context.Context, orderID, status string) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.Orders[orderID]
	if !ok {
		return nil, errors.New("order not found")
	}
	now := time.Now().UTC()
	order.Status = status
	order.UpdatedAt = now
	if status == domain.DeliveryStatusDelivered {
		order.DeliveredAt = &now
	}
	clone := *order
	return &clone, nil
}

func (r *Repository) AddOrderHistory(_ context.Context, event model.DeliveryStatusHistory) model.DeliveryStatusHistory {
	r.mu.Lock()
	defer r.mu.Unlock()
	event.ID = coalesceID(event.ID, utils.NewID("hst"))
	event.CreatedAt = time.Now().UTC()
	r.OrderStatusHistory[event.OrderID] = append(r.OrderStatusHistory[event.OrderID], event)
	return event
}

func (r *Repository) ListOrderHistory(_ context.Context, orderID string) []model.DeliveryStatusHistory {
	r.mu.RLock()
	defer r.mu.RUnlock()
	history := append([]model.DeliveryStatusHistory(nil), r.OrderStatusHistory[orderID]...)
	sort.Slice(history, func(i, j int) bool { return history[i].CreatedAt.Before(history[j].CreatedAt) })
	return history
}

func (r *Repository) GetActiveAssignmentForRider(_ context.Context, riderID string) (*model.DeliveryAssignment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, assignment := range r.Assignments {
		if assignment.RiderID == riderID && (assignment.Status == domain.AssignmentStatusAccepted || assignment.Status == domain.AssignmentStatusOffered) {
			clone := *assignment
			return &clone, nil
		}
	}
	return nil, errors.New("active assignment not found")
}

func (r *Repository) ListAssignedOrders(_ context.Context, riderID string) []model.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orders := make([]model.Order, 0)
	for _, assignment := range r.Assignments {
		if assignment.RiderID == riderID && assignment.Status != domain.AssignmentStatusCompleted && assignment.Status != domain.AssignmentStatusRejected && assignment.Status != domain.AssignmentStatusTimedOut {
			if order, ok := r.Orders[assignment.OrderID]; ok {
				orders = append(orders, *order)
			}
		}
	}
	sort.Slice(orders, func(i, j int) bool { return orders[i].CreatedAt.After(orders[j].CreatedAt) })
	return orders
}

func (r *Repository) ListOrderHistoryByRider(_ context.Context, riderID string) []model.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orders := make([]model.Order, 0)
	for _, assignment := range r.Assignments {
		if assignment.RiderID == riderID && assignment.Status == domain.AssignmentStatusCompleted {
			if order, ok := r.Orders[assignment.OrderID]; ok {
				orders = append(orders, *order)
			}
		}
	}
	sort.Slice(orders, func(i, j int) bool { return orders[i].UpdatedAt.After(orders[j].UpdatedAt) })
	return orders
}

func (r *Repository) SaveLatestLocation(_ context.Context, location model.RiderLocationLatest) model.RiderLocationLatest {
	r.mu.Lock()
	defer r.mu.Unlock()
	location.RecordedAt = time.Now().UTC()
	stored := location
	r.LatestLocations[location.RiderID] = &stored
	log := model.LocationLog{ID: utils.NewID("loc"), RiderID: location.RiderID, OrderID: location.OrderID, Latitude: location.Latitude, Longitude: location.Longitude, RecordedAt: location.RecordedAt, Source: location.Source}
	r.LocationLogs[location.RiderID] = append(r.LocationLogs[location.RiderID], log)
	return stored
}

func (r *Repository) BulkSaveLocations(_ context.Context, riderID string, points []model.RiderLocationLatest) []model.RiderLocationLatest {
	output := make([]model.RiderLocationLatest, 0, len(points))
	for _, point := range points {
		point.RiderID = riderID
		output = append(output, r.SaveLatestLocation(context.Background(), point))
	}
	return output
}

func (r *Repository) GetLatestLocation(_ context.Context, riderID string) (*model.RiderLocationLatest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	location, ok := r.LatestLocations[riderID]
	if !ok {
		return nil, errors.New("location not found")
	}
	clone := *location
	return &clone, nil
}

func (r *Repository) GetRouteHistory(_ context.Context, riderID string) []model.LocationLog {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]model.LocationLog(nil), r.LocationLogs[riderID]...)
}
func (r *Repository) GetWallet(_ context.Context, riderID string) (*model.Wallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	wallet, ok := r.Wallets[riderID]
	if !ok {
		return nil, errors.New("wallet not found")
	}
	clone := *wallet
	return &clone, nil
}

func (r *Repository) AddWalletTransaction(_ context.Context, riderID string, txn model.WalletTransaction) model.WalletTransaction {
	r.mu.Lock()
	defer r.mu.Unlock()
	txn.ID = coalesceID(txn.ID, utils.NewID("wtx"))
	txn.RiderID = riderID
	txn.CreatedAt = time.Now().UTC()
	if wallet := r.Wallets[riderID]; wallet != nil {
		if txn.Type == domain.WalletTxnTypeCredit {
			wallet.Balance += txn.Amount
		} else {
			wallet.Balance -= txn.Amount
		}
		wallet.UpdatedAt = txn.CreatedAt
		txn.WalletID = wallet.ID
	}
	r.WalletTransactions[riderID] = append(r.WalletTransactions[riderID], txn)
	return txn
}

func (r *Repository) ListWalletTransactions(_ context.Context, riderID string) []model.WalletTransaction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	txns := append([]model.WalletTransaction(nil), r.WalletTransactions[riderID]...)
	sort.Slice(txns, func(i, j int) bool { return txns[i].CreatedAt.After(txns[j].CreatedAt) })
	return txns
}

func (r *Repository) AddEarning(_ context.Context, riderID string, earning model.RiderEarning) model.RiderEarning {
	r.mu.Lock()
	defer r.mu.Unlock()
	earning.ID = coalesceID(earning.ID, utils.NewID("earn"))
	earning.RiderID = riderID
	earning.CreatedAt = time.Now().UTC()
	r.Earnings[riderID] = append(r.Earnings[riderID], earning)
	return earning
}

func (r *Repository) ListEarnings(_ context.Context, riderID string) []model.RiderEarning {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := append([]model.RiderEarning(nil), r.Earnings[riderID]...)
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items
}

func (r *Repository) CreatePayoutRequest(_ context.Context, riderID string, request model.PayoutRequest) model.PayoutRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	request.ID = coalesceID(request.ID, utils.NewID("pay"))
	request.RiderID = riderID
	request.RequestedAt = time.Now().UTC()
	if wallet := r.Wallets[riderID]; wallet != nil {
		request.WalletID = wallet.ID
		wallet.HoldBalance += request.Amount
		wallet.Balance -= request.Amount
		wallet.UpdatedAt = request.RequestedAt
	}
	r.PayoutRequests[riderID] = append(r.PayoutRequests[riderID], request)
	return request
}

func (r *Repository) ListPayoutRequests(_ context.Context, riderID string) []model.PayoutRequest {
	r.mu.RLock()
	defer r.mu.RUnlock()
	requests := append([]model.PayoutRequest(nil), r.PayoutRequests[riderID]...)
	sort.Slice(requests, func(i, j int) bool { return requests[i].RequestedAt.After(requests[j].RequestedAt) })
	return requests
}

func (r *Repository) GetPayoutRequest(_ context.Context, riderID, payoutID string) (*model.PayoutRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, payout := range r.PayoutRequests[riderID] {
		if payout.ID == payoutID {
			clone := payout
			return &clone, nil
		}
	}
	return nil, errors.New("payout not found")
}

func (r *Repository) ReviewPayout(_ context.Context, payoutID, status, reviewer, reason string) (*model.PayoutRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	for riderID, payouts := range r.PayoutRequests {
		for index := range payouts {
			if payouts[index].ID == payoutID {
				payouts[index].Status = status
				payouts[index].ReviewedAt = &now
				payouts[index].ReviewedBy = reviewer
				payouts[index].RejectionReason = reason
				if wallet := r.Wallets[riderID]; wallet != nil {
					if status == domain.PayoutStatusRejected {
						wallet.Balance += payouts[index].Amount
					}
					wallet.HoldBalance -= payouts[index].Amount
					wallet.UpdatedAt = now
				}
				r.PayoutRequests[riderID] = payouts
				clone := payouts[index]
				return &clone, nil
			}
		}
	}
	return nil, errors.New("payout not found")
}

func (r *Repository) ListRatings(_ context.Context, riderID string) []model.RatingReview {
	r.mu.RLock()
	defer r.mu.RUnlock()
	reviews := append([]model.RatingReview(nil), r.Ratings[riderID]...)
	sort.Slice(reviews, func(i, j int) bool { return reviews[i].CreatedAt.After(reviews[j].CreatedAt) })
	return reviews
}

func (r *Repository) ListNotifications(_ context.Context, userID string) []model.Notification {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := append([]model.Notification(nil), r.Notifications[userID]...)
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items
}

func (r *Repository) AddNotification(_ context.Context, notification model.Notification) model.Notification {
	r.mu.Lock()
	defer r.mu.Unlock()
	notification.ID = coalesceID(notification.ID, utils.NewID("ntf"))
	notification.CreatedAt = time.Now().UTC()
	r.Notifications[notification.UserID] = append(r.Notifications[notification.UserID], notification)
	return notification
}

func (r *Repository) MarkNotificationRead(_ context.Context, userID, notificationID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	notifications := r.Notifications[userID]
	now := time.Now().UTC()
	for index := range notifications {
		if notifications[index].ID == notificationID {
			notifications[index].Status = domain.NotificationStatusRead
			notifications[index].ReadAt = &now
			r.Notifications[userID] = notifications
			return nil
		}
	}
	return errors.New("notification not found")
}

func (r *Repository) MarkAllNotificationsRead(_ context.Context, userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	notifications := r.Notifications[userID]
	now := time.Now().UTC()
	for index := range notifications {
		notifications[index].Status = domain.NotificationStatusRead
		notifications[index].ReadAt = &now
	}
	r.Notifications[userID] = notifications
}

func (r *Repository) UpsertDeviceToken(_ context.Context, userID string, device model.DeviceToken) model.DeviceToken {
	r.mu.Lock()
	defer r.mu.Unlock()
	device.ID = utils.NewID("dev")
	device.UserID = userID
	device.UpdatedAt = time.Now().UTC()
	devices := r.DeviceTokens[userID]
	replaced := false
	for index := range devices {
		if devices[index].DeviceID == device.DeviceID {
			devices[index] = device
			replaced = true
			break
		}
	}
	if !replaced {
		devices = append(devices, device)
	}
	r.DeviceTokens[userID] = devices
	return device
}

func (r *Repository) DeleteDeviceToken(_ context.Context, userID, deviceID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	devices := r.DeviceTokens[userID]
	filtered := devices[:0]
	for _, device := range devices {
		if device.DeviceID != deviceID {
			filtered = append(filtered, device)
		}
	}
	r.DeviceTokens[userID] = filtered
}

func (r *Repository) CreateSupportTicket(_ context.Context, ticket model.SupportTicket) model.SupportTicket {
	r.mu.Lock()
	defer r.mu.Unlock()
	ticket.ID = coalesceID(ticket.ID, utils.NewID("sup"))
	ticket.CreatedAt = time.Now().UTC()
	ticket.UpdatedAt = ticket.CreatedAt
	r.SupportTickets[ticket.ID] = &ticket
	return ticket
}

func (r *Repository) ListSupportTickets(_ context.Context, riderID string) []model.SupportTicket {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.SupportTicket, 0)
	for _, ticket := range r.SupportTickets {
		if ticket.RiderID == riderID {
			items = append(items, *ticket)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
	return items
}

func (r *Repository) GetSupportTicket(_ context.Context, ticketID string) (*model.SupportTicket, []model.SupportTicketMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ticket, ok := r.SupportTickets[ticketID]
	if !ok {
		return nil, nil, errors.New("ticket not found")
	}
	clone := *ticket
	return &clone, append([]model.SupportTicketMessage(nil), r.SupportMessages[ticketID]...), nil
}

func (r *Repository) AddSupportMessage(_ context.Context, message model.SupportTicketMessage) model.SupportTicketMessage {
	r.mu.Lock()
	defer r.mu.Unlock()
	message.ID = coalesceID(message.ID, utils.NewID("msg"))
	message.CreatedAt = time.Now().UTC()
	r.SupportMessages[message.TicketID] = append(r.SupportMessages[message.TicketID], message)
	if ticket, ok := r.SupportTickets[message.TicketID]; ok {
		ticket.UpdatedAt = message.CreatedAt
	}
	return message
}

func (r *Repository) AppendAudit(_ context.Context, entry model.AuditLog) model.AuditLog {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry.ID = coalesceID(entry.ID, utils.NewID("aud"))
	entry.CreatedAt = time.Now().UTC()
	r.AuditLogs = append(r.AuditLogs, entry)
	return entry
}

func (r *Repository) ListRiders(_ context.Context) []struct {
	Rider model.Rider
	User  model.User
} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]struct {
		Rider model.Rider
		User  model.User
	}, 0, len(r.Riders))
	for _, rider := range r.Riders {
		user := r.Users[rider.UserID]
		if user == nil {
			continue
		}
		list = append(list, struct {
			Rider model.Rider
			User  model.User
		}{Rider: *cloneRider(rider), User: *cloneUser(user)})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].User.CreatedAt.After(list[j].User.CreatedAt) })
	return list
}

func (r *Repository) CreateRider(_ context.Context, user model.User, rider model.Rider, restaurantID string) (*model.User, *model.Rider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = coalesceID(user.ID, utils.NewID("usr"))
	rider.ID = coalesceID(rider.ID, utils.NewID("rdr"))
	rider.UserID = user.ID
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	rider.CreatedAt = user.CreatedAt
	rider.UpdatedAt = user.CreatedAt
	r.Users[user.ID] = &user
	r.Riders[rider.ID] = &rider
	r.RiderRestaurantLinks[rider.ID] = []model.RiderRestaurantAssignment{{ID: utils.NewID("rr"), RiderID: rider.ID, RestaurantID: restaurantID, IsPrimary: true, IsActive: true, AssignedBy: user.ID, AssignedAt: user.CreatedAt}}
	r.Wallets[rider.ID] = &model.Wallet{ID: utils.NewID("wal"), RiderID: rider.ID, Currency: "INR", UpdatedAt: user.CreatedAt}
	return cloneUser(&user), cloneRider(&rider)
}

func (r *Repository) UpdateRiderStatus(_ context.Context, riderID, status string) (*model.Rider, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rider, ok := r.Riders[riderID]
	if !ok {
		return nil, errors.New("rider not found")
	}
	rider.Status = status
	rider.UpdatedAt = time.Now().UTC()
	return cloneRider(rider), nil
}

func (r *Repository) ListUnassignedOrders(_ context.Context) []model.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]model.Order, 0)
	for _, order := range r.Orders {
		if order.Status == domain.DeliveryStatusReadyForAssignment {
			list = append(list, *order)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
	return list
}

func (r *Repository) AssignOrder(_ context.Context, orderID, riderID, assignedBy string) (*model.DeliveryAssignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.Orders[orderID]
	if !ok {
		return nil, errors.New("order not found")
	}
	if _, ok := r.Riders[riderID]; !ok {
		return nil, errors.New("rider not found")
	}
	assignedAt := time.Now().UTC()
	assignment := &model.DeliveryAssignment{ID: utils.NewID("asg"), OrderID: orderID, RiderID: riderID, RestaurantID: order.RestaurantID, Status: domain.AssignmentStatusOffered, AssignedBy: assignedBy, AssignedAt: assignedAt, DecisionDeadlineAt: assignedAt.Add(time.Duration(r.Config.OrderAcceptTimeoutSeconds) * time.Second)}
	r.Assignments[assignment.ID] = assignment
	order.Status = domain.DeliveryStatusAssigned
	order.UpdatedAt = assignedAt
	clone := *assignment
	return &clone, nil
}

func (r *Repository) ReassignOrder(_ context.Context, orderID, riderID, assignedBy string) (*model.DeliveryAssignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.Orders[orderID]
	if !ok {
		return nil, errors.New("order not found")
	}
	for _, assignment := range r.Assignments {
		if assignment.OrderID == orderID && assignment.Status == domain.AssignmentStatusOffered {
			assignment.Status = domain.AssignmentStatusCancelled
		}
	}
	assignedAt := time.Now().UTC()
	assignment := &model.DeliveryAssignment{ID: utils.NewID("asg"), OrderID: orderID, RiderID: riderID, RestaurantID: order.RestaurantID, Status: domain.AssignmentStatusOffered, AssignedBy: assignedBy, AssignedAt: assignedAt, DecisionDeadlineAt: assignedAt.Add(time.Duration(r.Config.OrderAcceptTimeoutSeconds) * time.Second)}
	r.Assignments[assignment.ID] = assignment
	order.Status = domain.DeliveryStatusAssigned
	order.UpdatedAt = assignedAt
	clone := *assignment
	return &clone, nil
}

func (r *Repository) ListLiveOrders(_ context.Context) []model.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.Order, 0)
	for _, order := range r.Orders {
		switch order.Status {
		case domain.DeliveryStatusAssigned, domain.DeliveryStatusAccepted, domain.DeliveryStatusReachedRestaurant, domain.DeliveryStatusPickupVerified, domain.DeliveryStatusPickedUp, domain.DeliveryStatusOnTheWay, domain.DeliveryStatusReachedCustomer, domain.DeliveryStatusDeliveryVerified:
			items = append(items, *order)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
	return items
}

func (r *Repository) GetSystemConfig(_ context.Context) SystemConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Config
}
func (r *Repository) UpdateSystemConfig(_ context.Context, cfg SystemConfig) SystemConfig {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Config = cfg
	return r.Config
}
func (r *Repository) GetRestaurantAssignments(_ context.Context, riderID string) []model.RiderRestaurantAssignment {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]model.RiderRestaurantAssignment(nil), r.RiderRestaurantLinks[riderID]...)
}

func (r *Repository) getRiderByUserIDLocked(userID string) *model.Rider {
	for _, rider := range r.Riders {
		if rider.UserID == userID {
			return rider
		}
	}
	return nil
}

func cloneUser(user *model.User) *model.User {
	if user == nil {
		return nil
	}
	clone := *user
	clone.Roles = append([]string(nil), user.Roles...)
	return &clone
}

func cloneRider(rider *model.Rider) *model.Rider {
	if rider == nil {
		return nil
	}
	clone := *rider
	return &clone
}

func coalesceID(current, fallback string) string {
	if strings.TrimSpace(current) != "" {
		return current
	}
	return fallback
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
func (r *Repository) FindAssignmentByOrderAndRider(_ context.Context, orderID, riderID string) (*model.DeliveryAssignment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, assignment := range r.Assignments {
		if assignment.OrderID == orderID && assignment.RiderID == riderID {
			clone := *assignment
			return &clone, nil
		}
	}
	return nil, errors.New("assignment not found")
}

func (r *Repository) UpdateAssignmentStatus(_ context.Context, assignmentID, status, reason string) (*model.DeliveryAssignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	assignment, ok := r.Assignments[assignmentID]
	if !ok {
		return nil, errors.New("assignment not found")
	}
	now := time.Now().UTC()
	assignment.Status = status
	switch status {
	case domain.AssignmentStatusAccepted:
		assignment.AcceptedAt = &now
	case domain.AssignmentStatusRejected:
		assignment.RejectedAt = &now
		assignment.RejectReason = reason
	case domain.AssignmentStatusFailed, domain.AssignmentStatusCancelled:
		assignment.FailureReason = reason
	}
	clone := *assignment
	return &clone, nil
}
func (r *Repository) ListAllPayoutRequests(_ context.Context) []model.PayoutRequest {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]model.PayoutRequest, 0)
	for _, payouts := range r.PayoutRequests {
		list = append(list, payouts...)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].RequestedAt.After(list[j].RequestedAt) })
	return list
}
