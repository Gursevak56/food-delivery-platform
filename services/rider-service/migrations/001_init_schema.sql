CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE role_code AS ENUM ('SUPER_ADMIN','RESTAURANT_OWNER','RESTAURANT_ADMIN','RESTAURANT_STAFF','DISPATCHER_MANAGER','RIDER');
CREATE TYPE user_status AS ENUM ('ACTIVE','INACTIVE','SUSPENDED');
CREATE TYPE rider_availability AS ENUM ('OFFLINE','ONLINE','BREAK','BUSY');
CREATE TYPE rider_doc_status AS ENUM ('PENDING_REVIEW','APPROVED','REJECTED','EXPIRED');
CREATE TYPE assignment_status AS ENUM ('OFFERED','ACCEPTED','REJECTED','TIMED_OUT','COMPLETED','FAILED','CANCELLED');
CREATE TYPE delivery_status AS ENUM ('CREATED','READY_FOR_ASSIGNMENT','ASSIGNED','ACCEPTED','REACHED_RESTAURANT','PICKUP_VERIFIED','PICKED_UP','ON_THE_WAY','REACHED_CUSTOMER','DELIVERY_VERIFIED','DELIVERED','FAILED','CANCELLED');
CREATE TYPE otp_purpose AS ENUM ('LOGIN','PICKUP','DELIVERY');
CREATE TYPE otp_status AS ENUM ('PENDING','VERIFIED','EXPIRED','LOCKED');
CREATE TYPE wallet_txn_type AS ENUM ('CREDIT','DEBIT');
CREATE TYPE payout_status AS ENUM ('PENDING','APPROVED','PROCESSING','PAID','REJECTED');
CREATE TYPE rating_source AS ENUM ('CUSTOMER','RESTAURANT');
CREATE TYPE notification_status AS ENUM ('UNREAD','READ');
CREATE TYPE support_ticket_status AS ENUM ('OPEN','IN_PROGRESS','RESOLVED','CLOSED');
CREATE TYPE support_ticket_priority AS ENUM ('LOW','MEDIUM','HIGH','CRITICAL');

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code role_code NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    primary_role role_code NOT NULL,
    status user_status NOT NULL DEFAULT 'ACTIVE',
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    photo_url TEXT,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id),
    role_id UUID NOT NULL REFERENCES roles(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE restaurants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    city TEXT,
    area TEXT,
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE riders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id),
    status user_status NOT NULL DEFAULT 'ACTIVE',
    availability_status rider_availability NOT NULL DEFAULT 'OFFLINE',
    approval_status TEXT NOT NULL DEFAULT 'PENDING',
    kyc_status TEXT NOT NULL DEFAULT 'PENDING',
    acceptance_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    rejection_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    completion_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    cancellation_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    avg_rating NUMERIC(3,2) NOT NULL DEFAULT 0,
    rating_count INTEGER NOT NULL DEFAULT 0,
    active_hours_today NUMERIC(10,2) NOT NULL DEFAULT 0,
    current_shift_id UUID,
    current_break_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE rider_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL REFERENCES riders(id),
    document_type TEXT NOT NULL,
    document_url TEXT NOT NULL,
    status rider_doc_status NOT NULL DEFAULT 'PENDING_REVIEW',
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rider_vehicle_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL UNIQUE REFERENCES riders(id),
    vehicle_type TEXT NOT NULL,
    registration_no TEXT NOT NULL,
    color TEXT,
    capacity_kg NUMERIC(10,2),
    insurance_expiry DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE rider_bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL UNIQUE REFERENCES riders(id),
    account_holder TEXT NOT NULL,
    bank_name TEXT NOT NULL,
    account_number TEXT NOT NULL,
    ifsc_code TEXT NOT NULL,
    upi_id TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE rider_restaurant_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL REFERENCES riders(id),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id),
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (rider_id, restaurant_id)
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number TEXT NOT NULL UNIQUE,
    restaurant_id UUID NOT NULL REFERENCES restaurants(id),
    customer_name TEXT NOT NULL,
    customer_phone TEXT NOT NULL,
    delivery_address TEXT NOT NULL,
    delivery_latitude NUMERIC(10,7),
    delivery_longitude NUMERIC(10,7),
    status delivery_status NOT NULL DEFAULT 'CREATED',
    pickup_otp_required BOOLEAN NOT NULL DEFAULT FALSE,
    delivery_otp_required BOOLEAN NOT NULL DEFAULT TRUE,
    distance_km NUMERIC(10,2) NOT NULL DEFAULT 0,
    base_payout NUMERIC(10,2) NOT NULL DEFAULT 0,
    distance_payout NUMERIC(10,2) NOT NULL DEFAULT 0,
    waiting_charges NUMERIC(10,2) NOT NULL DEFAULT 0,
    surge_bonus NUMERIC(10,2) NOT NULL DEFAULT 0,
    tip_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    delivered_at TIMESTAMPTZ
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE delivery_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id),
    rider_id UUID NOT NULL REFERENCES riders(id),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id),
    status assignment_status NOT NULL DEFAULT 'OFFERED',
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    accepted_at TIMESTAMPTZ,
    rejected_at TIMESTAMPTZ,
    decision_deadline_at TIMESTAMPTZ NOT NULL,
    reject_reason TEXT,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE delivery_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    assignment_id UUID REFERENCES delivery_assignments(id),
    status delivery_status NOT NULL,
    actor_id UUID REFERENCES users(id),
    actor_role role_code,
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_otps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    purpose otp_purpose NOT NULL,
    code_hash TEXT NOT NULL,
    status otp_status NOT NULL DEFAULT 'PENDING',
    retry_count INTEGER NOT NULL DEFAULT 0,
    resend_count INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE rider_shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL REFERENCES riders(id),
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    started_by UUID REFERENCES users(id),
    ended_by UUID REFERENCES users(id),
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    break_minutes INTEGER NOT NULL DEFAULT 0,
    active_minutes INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rider_breaks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL REFERENCES riders(id),
    shift_id UUID NOT NULL REFERENCES rider_shifts(id) ON DELETE CASCADE,
    reason TEXT,
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rider_locations_latest (
    rider_id UUID PRIMARY KEY REFERENCES riders(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id),
    latitude NUMERIC(10,7) NOT NULL,
    longitude NUMERIC(10,7) NOT NULL,
    accuracy_meters NUMERIC(10,2),
    speed_kph NUMERIC(10,2),
    heading_degrees NUMERIC(10,2),
    battery_level INTEGER,
    source TEXT NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL UNIQUE REFERENCES riders(id),
    balance NUMERIC(12,2) NOT NULL DEFAULT 0,
    hold_balance NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency_code TEXT NOT NULL DEFAULT 'INR',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    rider_id UUID NOT NULL REFERENCES riders(id),
    transaction_type wallet_txn_type NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    reference_type TEXT,
    reference_id UUID,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rider_earnings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL REFERENCES riders(id),
    order_id UUID NOT NULL REFERENCES orders(id),
    base_payout NUMERIC(10,2) NOT NULL DEFAULT 0,
    distance_payout NUMERIC(10,2) NOT NULL DEFAULT 0,
    waiting_charges NUMERIC(10,2) NOT NULL DEFAULT 0,
    surge_bonus NUMERIC(10,2) NOT NULL DEFAULT 0,
    tip_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
    incentive_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
    penalty_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
    cancellation_compensation NUMERIC(10,2) NOT NULL DEFAULT 0,
    net_earning NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE payout_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL REFERENCES riders(id),
    wallet_id UUID NOT NULL REFERENCES wallets(id),
    amount NUMERIC(12,2) NOT NULL,
    status payout_status NOT NULL DEFAULT 'PENDING',
    reviewed_by UUID REFERENCES users(id),
    rejection_reason TEXT,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ratings_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL REFERENCES riders(id),
    order_id UUID REFERENCES orders(id),
    source rating_source NOT NULL,
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    review TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    order_id UUID REFERENCES orders(id),
    notification_type TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    channel TEXT NOT NULL DEFAULT 'PUSH',
    status notification_status NOT NULL DEFAULT 'UNREAD',
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE support_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id UUID NOT NULL REFERENCES riders(id),
    order_id UUID REFERENCES orders(id),
    subject TEXT NOT NULL,
    category TEXT NOT NULL,
    priority support_ticket_priority NOT NULL DEFAULT 'MEDIUM',
    status support_ticket_status NOT NULL DEFAULT 'OPEN',
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE support_ticket_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    actor_id UUID REFERENCES users(id),
    actor_role role_code,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    device_id TEXT NOT NULL,
    device_name TEXT NOT NULL,
    refresh_token_hash TEXT NOT NULL UNIQUE,
    last_ip TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES users(id),
    actor_role role_code,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id UUID,
    request_id TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE system_configurations (
    config_key TEXT PRIMARY KEY,
    config_value JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uidx_users_email_active ON users (LOWER(email)) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX uidx_users_phone_active ON users (phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_riders_status_availability ON riders (status, availability_status);
CREATE INDEX idx_rider_documents_lookup ON rider_documents (rider_id, status);
CREATE INDEX idx_rider_restaurant_assignments_lookup ON rider_restaurant_assignments (restaurant_id, rider_id) WHERE is_active = TRUE;
CREATE INDEX idx_orders_status_restaurant_time ON orders (status, restaurant_id, created_at DESC);
CREATE INDEX idx_delivery_assignments_rider_status_deadline ON delivery_assignments (rider_id, status, decision_deadline_at);
CREATE INDEX idx_delivery_assignments_order_status ON delivery_assignments (order_id, status);
CREATE INDEX idx_delivery_status_history_order_time ON delivery_status_history (order_id, created_at);
CREATE INDEX idx_order_otps_lookup ON order_otps (order_id, purpose, created_at DESC);
CREATE INDEX idx_rider_shifts_rider_time ON rider_shifts (rider_id, started_at DESC);
CREATE INDEX idx_rider_breaks_shift_time ON rider_breaks (shift_id, started_at DESC);
CREATE INDEX idx_wallet_transactions_wallet_time ON wallet_transactions (wallet_id, created_at DESC);
CREATE INDEX idx_rider_earnings_rider_time ON rider_earnings (rider_id, created_at DESC);
CREATE INDEX idx_payout_requests_rider_status_time ON payout_requests (rider_id, status, requested_at DESC);
CREATE INDEX idx_ratings_reviews_rider_time ON ratings_reviews (rider_id, created_at DESC);
CREATE INDEX idx_notifications_user_status_time ON notifications (user_id, status, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_support_tickets_rider_status_time ON support_tickets (rider_id, status, updated_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_support_ticket_messages_ticket_time ON support_ticket_messages (ticket_id, created_at);
CREATE INDEX idx_refresh_tokens_user_device ON refresh_tokens (user_id, device_id);
CREATE INDEX idx_audit_logs_actor_time ON audit_logs (actor_id, created_at DESC);
