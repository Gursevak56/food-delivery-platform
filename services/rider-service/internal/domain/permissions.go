package domain

var RolePermissions = map[string][]string{
	RoleSuperAdmin: {
		"riders:read", "riders:write", "orders:assign", "orders:reassign", "orders:track", "analytics:read", "config:write", "payouts:review",
	},
	RoleRestaurantOwner: {
		"riders:read", "orders:assign", "orders:reassign", "orders:track", "analytics:read", "payouts:review",
	},
	RoleRestaurantAdmin: {
		"riders:read", "orders:assign", "orders:reassign", "orders:track", "analytics:read",
	},
	RoleRestaurantStaff: {
		"orders:track",
	},
	RoleDispatcherManager: {
		"riders:read", "orders:assign", "orders:reassign", "orders:track", "analytics:read",
	},
	RoleRider: {
		"auth:self", "profile:self", "shift:self", "orders:self", "wallet:self", "support:self", "notifications:self",
	},
}
