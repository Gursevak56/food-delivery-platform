package dto

type LoginRequest struct {
	Login      string `json:"login" validate:"required"`
	Password   string `json:"password" validate:"required"`
	DeviceID   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
}

type OTPRequest struct {
	Login string `json:"login" validate:"required"`
}

type OTPVerifyRequest struct {
	Login      string `json:"login" validate:"required"`
	OTP        string `json:"otp" validate:"required,len=6,numeric"`
	DeviceID   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	DeviceID     string `json:"device_id" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Login string `json:"login" validate:"required"`
}

type ResetPasswordRequest struct {
	Login       string `json:"login" validate:"required"`
	OTP         string `json:"otp" validate:"required,len=6,numeric"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type UpdateRiderProfileRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required"`
}

type UpdatePhotoRequest struct {
	PhotoURL string `json:"photo_url" validate:"required,url"`
}

type UpdateVehicleRequest struct {
	VehicleType     string  `json:"vehicle_type" validate:"required"`
	RegistrationNo  string  `json:"registration_no" validate:"required"`
	Color           string  `json:"color" validate:"required"`
	Capacity        float64 `json:"capacity" validate:"required,gt=0"`
	InsuranceExpiry string  `json:"insurance_expiry" validate:"required"`
}

type RiderDocumentInput struct {
	DocumentType string `json:"document_type" validate:"required"`
	DocumentURL  string `json:"document_url" validate:"required,url"`
}

type UpdateDocumentsRequest struct {
	Documents []RiderDocumentInput `json:"documents" validate:"required,min=1,dive"`
}

type UpdateBankAccountRequest struct {
	AccountHolder string `json:"account_holder" validate:"required"`
	BankName      string `json:"bank_name" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	IFSCCode      string `json:"ifsc_code" validate:"required"`
	UPIID         string `json:"upi_id"`
}

type BreakRequest struct {
	Reason string `json:"reason"`
}

type ShiftEndRequest struct {
	Reason string `json:"reason"`
}

type RejectOrderRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type VerifyOrderOTPRequest struct {
	OTP string `json:"otp" validate:"required,len=6,numeric"`
}

type OrderFailureRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type LocationUpdateRequest struct {
	OrderID        string  `json:"order_id"`
	Latitude       float64 `json:"latitude" validate:"required"`
	Longitude      float64 `json:"longitude" validate:"required"`
	AccuracyMeters float64 `json:"accuracy_meters"`
	SpeedKPH       float64 `json:"speed_kph"`
	HeadingDegrees float64 `json:"heading_degrees"`
	BatteryLevel   int     `json:"battery_level"`
	Source         string  `json:"source" validate:"required"`
}

type BulkLocationUpdateRequest struct {
	Points []LocationUpdateRequest `json:"points" validate:"required,min=1,dive"`
}

type RequestPayoutRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type DeviceTokenRequest struct {
	DeviceID string `json:"device_id" validate:"required"`
	Platform string `json:"platform" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

type SupportTicketRequest struct {
	Subject     string `json:"subject" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Priority    string `json:"priority" validate:"required"`
	OrderID     string `json:"order_id"`
	Description string `json:"description" validate:"required"`
}

type SupportReplyRequest struct {
	Message string `json:"message" validate:"required"`
}

type CreateRiderRequest struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required"`
	Password     string `json:"password" validate:"required,min=8"`
	RestaurantID string `json:"restaurant_id" validate:"required"`
}

type UpdateRiderStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type AssignRiderRequest struct {
	RiderID string `json:"rider_id" validate:"required"`
}

type UpdateConfigRequest struct {
	OrderAcceptTimeoutSeconds int     `json:"order_accept_timeout_seconds" validate:"required,gte=15"`
	PickupOTPRequired         bool    `json:"pickup_otp_required"`
	DeliveryOTPRequired       bool    `json:"delivery_otp_required"`
	RiderMaxActiveOrders      int     `json:"rider_max_active_orders" validate:"required,gte=1"`
	OTPMaxRetries             int     `json:"otp_max_retries" validate:"required,gte=1"`
	OTPMaxResends             int     `json:"otp_max_resends" validate:"required,gte=1"`
	SurgeMultiplierDefault    float64 `json:"surge_multiplier_default" validate:"required,gte=1"`
	MinimumPayoutAmount       float64 `json:"minimum_payout_amount" validate:"required,gte=0"`
}
