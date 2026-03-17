package domain

const (
	RoleSuperAdmin        = "SUPER_ADMIN"
	RoleRestaurantOwner   = "RESTAURANT_OWNER"
	RoleRestaurantAdmin   = "RESTAURANT_ADMIN"
	RoleRestaurantStaff   = "RESTAURANT_STAFF"
	RoleDispatcherManager = "DISPATCHER_MANAGER"
	RoleRider             = "RIDER"
)

const (
	UserStatusActive    = "ACTIVE"
	UserStatusInactive  = "INACTIVE"
	UserStatusSuspended = "SUSPENDED"
)

const (
	RiderAvailabilityOffline = "OFFLINE"
	RiderAvailabilityOnline  = "ONLINE"
	RiderAvailabilityBreak   = "BREAK"
	RiderAvailabilityBusy    = "BUSY"
)

const (
	ShiftStatusScheduled = "SCHEDULED"
	ShiftStatusActive    = "ACTIVE"
	ShiftStatusEnded     = "ENDED"
)

const (
	BreakStatusActive = "ACTIVE"
	BreakStatusEnded  = "ENDED"
)

const (
	AssignmentStatusOffered   = "OFFERED"
	AssignmentStatusAccepted  = "ACCEPTED"
	AssignmentStatusRejected  = "REJECTED"
	AssignmentStatusTimedOut  = "TIMED_OUT"
	AssignmentStatusCompleted = "COMPLETED"
	AssignmentStatusFailed    = "FAILED"
	AssignmentStatusCancelled = "CANCELLED"
)

const (
	DeliveryStatusCreated            = "CREATED"
	DeliveryStatusReadyForAssignment = "READY_FOR_ASSIGNMENT"
	DeliveryStatusAssigned           = "ASSIGNED"
	DeliveryStatusAccepted           = "ACCEPTED"
	DeliveryStatusReachedRestaurant  = "REACHED_RESTAURANT"
	DeliveryStatusPickupVerified     = "PICKUP_VERIFIED"
	DeliveryStatusPickedUp           = "PICKED_UP"
	DeliveryStatusOnTheWay           = "ON_THE_WAY"
	DeliveryStatusReachedCustomer    = "REACHED_CUSTOMER"
	DeliveryStatusDeliveryVerified   = "DELIVERY_VERIFIED"
	DeliveryStatusDelivered          = "DELIVERED"
	DeliveryStatusFailed             = "FAILED"
	DeliveryStatusCancelled          = "CANCELLED"
)

const (
	OTPPurposeLogin    = "LOGIN"
	OTPPurposePickup   = "PICKUP"
	OTPPurposeDelivery = "DELIVERY"
)

const (
	OTPStatusPending  = "PENDING"
	OTPStatusVerified = "VERIFIED"
	OTPStatusExpired  = "EXPIRED"
	OTPStatusLocked   = "LOCKED"
)

const (
	WalletTxnTypeCredit = "CREDIT"
	WalletTxnTypeDebit  = "DEBIT"
)

const (
	PayoutStatusPending    = "PENDING"
	PayoutStatusApproved   = "APPROVED"
	PayoutStatusProcessing = "PROCESSING"
	PayoutStatusPaid       = "PAID"
	PayoutStatusRejected   = "REJECTED"
)

const (
	SupportTicketStatusOpen       = "OPEN"
	SupportTicketStatusInProgress = "IN_PROGRESS"
	SupportTicketStatusResolved   = "RESOLVED"
	SupportTicketStatusClosed     = "CLOSED"
)

const (
	NotificationStatusUnread = "UNREAD"
	NotificationStatusRead   = "READ"
)

const (
	RatingSourceCustomer   = "CUSTOMER"
	RatingSourceRestaurant = "RESTAURANT"
)
