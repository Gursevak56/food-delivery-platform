package model

import "time"

type User struct {
	ID            string     `json:"id"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	PasswordHash  string     `json:"-"`
	Status        string     `json:"status"`
	Roles         []string   `json:"roles"`
	PrimaryRole   string     `json:"primary_role"`
	PhotoURL      string     `json:"photo_url,omitempty"`
	LastLoginAt   time.Time  `json:"last_login_at,omitempty"`
	EmailVerified bool       `json:"email_verified"`
	PhoneVerified bool       `json:"phone_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type Restaurant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Status    string    `json:"status"`
	City      string    `json:"city"`
	Area      string    `json:"area"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Rider struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	Status             string    `json:"status"`
	AvailabilityStatus string    `json:"availability_status"`
	KYCStatus          string    `json:"kyc_status"`
	ApprovalStatus     string    `json:"approval_status"`
	CurrentShiftID     string    `json:"current_shift_id,omitempty"`
	CurrentBreakID     string    `json:"current_break_id,omitempty"`
	VehicleType        string    `json:"vehicle_type,omitempty"`
	VehicleNumber      string    `json:"vehicle_number,omitempty"`
	AvgRating          float64   `json:"avg_rating"`
	RatingCount        int       `json:"rating_count"`
	AcceptanceRate     float64   `json:"acceptance_rate"`
	CompletionRate     float64   `json:"completion_rate"`
	CancellationRate   float64   `json:"cancellation_rate"`
	RejectionRate      float64   `json:"rejection_rate"`
	ActiveHoursToday   float64   `json:"active_hours_today"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type RiderDocument struct {
	ID           string     `json:"id"`
	RiderID      string     `json:"rider_id"`
	DocumentType string     `json:"document_type"`
	DocumentURL  string     `json:"document_url"`
	Status       string     `json:"status"`
	UploadedAt   time.Time  `json:"uploaded_at"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
}

type RiderVehicle struct {
	ID              string    `json:"id"`
	RiderID         string    `json:"rider_id"`
	VehicleType     string    `json:"vehicle_type"`
	RegistrationNo  string    `json:"registration_no"`
	Color           string    `json:"color"`
	Capacity        float64   `json:"capacity"`
	InsuranceExpiry time.Time `json:"insurance_expiry"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type RiderBankAccount struct {
	ID            string     `json:"id"`
	RiderID       string     `json:"rider_id"`
	AccountHolder string     `json:"account_holder"`
	BankName      string     `json:"bank_name"`
	AccountNumber string     `json:"account_number"`
	IFSCCode      string     `json:"ifsc_code"`
	UPIID         string     `json:"upi_id,omitempty"`
	IsVerified    bool       `json:"is_verified"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type RiderRestaurantAssignment struct {
	ID           string    `json:"id"`
	RiderID      string    `json:"rider_id"`
	RestaurantID string    `json:"restaurant_id"`
	IsPrimary    bool      `json:"is_primary"`
	IsActive     bool      `json:"is_active"`
	AssignedBy   string    `json:"assigned_by"`
	AssignedAt   time.Time `json:"assigned_at"`
}

type Order struct {
	ID                  string     `json:"id"`
	OrderNumber         string     `json:"order_number"`
	RestaurantID        string     `json:"restaurant_id"`
	CustomerName        string     `json:"customer_name"`
	CustomerPhone       string     `json:"customer_phone"`
	DeliveryAddress     string     `json:"delivery_address"`
	DeliveryLatitude    float64    `json:"delivery_latitude"`
	DeliveryLongitude   float64    `json:"delivery_longitude"`
	Status              string     `json:"status"`
	DistanceKM          float64    `json:"distance_km"`
	BasePayout          float64    `json:"base_payout"`
	DistancePayout      float64    `json:"distance_payout"`
	WaitingCharges      float64    `json:"waiting_charges"`
	SurgeBonus          float64    `json:"surge_bonus"`
	TipAmount           float64    `json:"tip_amount"`
	TotalAmount         float64    `json:"total_amount"`
	DeliveryNote        string     `json:"delivery_note,omitempty"`
	PickupOTPRequired   bool       `json:"pickup_otp_required"`
	DeliveryOTPRequired bool       `json:"delivery_otp_required"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeliveredAt         *time.Time `json:"delivered_at,omitempty"`
}

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type DeliveryAssignment struct {
	ID                 string     `json:"id"`
	OrderID            string     `json:"order_id"`
	RiderID            string     `json:"rider_id"`
	RestaurantID       string     `json:"restaurant_id"`
	Status             string     `json:"status"`
	AssignedBy         string     `json:"assigned_by"`
	AssignedAt         time.Time  `json:"assigned_at"`
	AcceptedAt         *time.Time `json:"accepted_at,omitempty"`
	RejectedAt         *time.Time `json:"rejected_at,omitempty"`
	DecisionDeadlineAt time.Time  `json:"decision_deadline_at"`
	RejectReason       string     `json:"reject_reason,omitempty"`
	FailureReason      string     `json:"failure_reason,omitempty"`
}

type DeliveryStatusHistory struct {
	ID           string    `json:"id"`
	OrderID      string    `json:"order_id"`
	AssignmentID string    `json:"assignment_id"`
	Status       string    `json:"status"`
	ActorID      string    `json:"actor_id"`
	ActorRole    string    `json:"actor_role"`
	Comment      string    `json:"comment,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type OTP struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id,omitempty"`
	OrderID     string     `json:"order_id,omitempty"`
	Purpose     string     `json:"purpose"`
	Code        string     `json:"-"`
	Status      string     `json:"status"`
	RetryCount  int        `json:"retry_count"`
	ResendCount int        `json:"resend_count"`
	ExpiresAt   time.Time  `json:"expires_at"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type RiderShift struct {
	ID            string     `json:"id"`
	RiderID       string     `json:"rider_id"`
	Status        string     `json:"status"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at,omitempty"`
	BreakMinutes  int        `json:"break_minutes"`
	ActiveMinutes int        `json:"active_minutes"`
	StartedBy     string     `json:"started_by"`
	EndedBy       string     `json:"ended_by,omitempty"`
}

type RiderBreak struct {
	ID        string     `json:"id"`
	RiderID   string     `json:"rider_id"`
	ShiftID   string     `json:"shift_id"`
	Reason    string     `json:"reason,omitempty"`
	Status    string     `json:"status"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

type RiderLocationLatest struct {
	RiderID        string    `json:"rider_id"`
	OrderID        string    `json:"order_id,omitempty"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	AccuracyMeters float64   `json:"accuracy_meters"`
	SpeedKPH       float64   `json:"speed_kph"`
	HeadingDegrees float64   `json:"heading_degrees"`
	BatteryLevel   int       `json:"battery_level"`
	Source         string    `json:"source"`
	RecordedAt     time.Time `json:"recorded_at"`
}

type LocationLog struct {
	ID         string    `json:"id"`
	RiderID    string    `json:"rider_id"`
	OrderID    string    `json:"order_id,omitempty"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	RecordedAt time.Time `json:"recorded_at"`
	Source     string    `json:"source"`
}

type Wallet struct {
	ID          string    `json:"id"`
	RiderID     string    `json:"rider_id"`
	Balance     float64   `json:"balance"`
	HoldBalance float64   `json:"hold_balance"`
	Currency    string    `json:"currency"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WalletTransaction struct {
	ID            string    `json:"id"`
	WalletID      string    `json:"wallet_id"`
	RiderID       string    `json:"rider_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	ReferenceID   string    `json:"reference_id"`
	ReferenceType string    `json:"reference_type"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

type RiderEarning struct {
	ID                       string    `json:"id"`
	RiderID                  string    `json:"rider_id"`
	OrderID                  string    `json:"order_id"`
	BasePayout               float64   `json:"base_payout"`
	DistancePayout           float64   `json:"distance_payout"`
	WaitingCharges           float64   `json:"waiting_charges"`
	SurgeBonus               float64   `json:"surge_bonus"`
	TipAmount                float64   `json:"tip_amount"`
	IncentiveAmount          float64   `json:"incentive_amount"`
	PenaltyAmount            float64   `json:"penalty_amount"`
	CancellationCompensation float64   `json:"cancellation_compensation"`
	NetEarning               float64   `json:"net_earning"`
	CreatedAt                time.Time `json:"created_at"`
}

type PayoutRequest struct {
	ID              string     `json:"id"`
	RiderID         string     `json:"rider_id"`
	WalletID        string     `json:"wallet_id"`
	Amount          float64    `json:"amount"`
	Status          string     `json:"status"`
	RequestedAt     time.Time  `json:"requested_at"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy      string     `json:"reviewed_by,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
}

type RatingReview struct {
	ID        string    `json:"id"`
	RiderID   string    `json:"rider_id"`
	OrderID   string    `json:"order_id"`
	Source    string    `json:"source"`
	Rating    int       `json:"rating"`
	Review    string    `json:"review"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Status    string     `json:"status"`
	Channel   string     `json:"channel"`
	OrderID   string     `json:"order_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

type DeviceToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	DeviceID  string    `json:"device_id"`
	Platform  string    `json:"platform"`
	Token     string    `json:"token"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SupportTicket struct {
	ID          string    `json:"id"`
	RiderID     string    `json:"rider_id"`
	Subject     string    `json:"subject"`
	Category    string    `json:"category"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	OrderID     string    `json:"order_id,omitempty"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SupportTicketMessage struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	ActorID   string    `json:"actor_id"`
	ActorRole string    `json:"actor_role"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	DeviceID    string     `json:"device_id"`
	DeviceName  string     `json:"device_name"`
	RefreshHash string     `json:"-"`
	LastIP      string     `json:"last_ip,omitempty"`
	LastSeenAt  time.Time  `json:"last_seen_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type AuditLog struct {
	ID         string         `json:"id"`
	ActorID    string         `json:"actor_id"`
	ActorRole  string         `json:"actor_role"`
	Action     string         `json:"action"`
	EntityType string         `json:"entity_type"`
	EntityID   string         `json:"entity_id"`
	RequestID  string         `json:"request_id"`
	Metadata   map[string]any `json:"metadata"`
	CreatedAt  time.Time      `json:"created_at"`
}
