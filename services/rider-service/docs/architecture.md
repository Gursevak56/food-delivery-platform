# Architecture

## Service shape

The rider backend is implemented as a clean modular monolith with explicit module boundaries:

- Auth and Session Management
- Rider Profile and Compliance
- Shift and Availability
- Dispatch and Assignment
- Active Delivery Lifecycle
- OTP Verification
- Live Location and Tracking
- Earnings, Wallet, and Payouts
- Ratings and Performance
- Notifications
- Support
- Admin and Dispatcher Controls
- Audit and Activity Tracking

This shape keeps deployment and ownership simple for one team while remaining easy to split into independent services later. Natural future service boundaries are:

- `auth-service`
- `dispatch-service`
- `delivery-tracking-service`
- `rider-finance-service`
- `notification-service`
- `support-service`

## Request flow

1. HTTP request enters Gin router.
2. Request ID, structured access log, recovery, rate limiting, auth, and RBAC middleware run.
3. Handler validates DTOs and delegates to service layer.
4. Service layer enforces business rules and coordinates repository writes.
5. Repository persists transactional records to PostgreSQL and flexible event/log payloads to MongoDB.
6. Response helper emits consistent success or error payloads.

## Why PostgreSQL vs MongoDB

### PostgreSQL

Use PostgreSQL for entities that need:

- strong consistency
- relational joins
- transactional integrity
- aggregates and reporting
- constraints and referential integrity
- auditable financial or assignment state

Tables that belong in PostgreSQL:

- `users`, `roles`, `user_roles`
- `restaurants`
- `riders`, `rider_documents`, `rider_vehicle_details`, `rider_bank_accounts`
- `rider_restaurant_assignments`
- `orders`, `order_items`
- `delivery_assignments`, `delivery_status_history`
- `order_otps`
- `rider_shifts`, `rider_breaks`
- `rider_locations_latest`
- `wallets`, `wallet_transactions`, `rider_earnings`, `payout_requests`
- `ratings_reviews`
- `notifications`
- `support_tickets`, `support_ticket_messages`
- `refresh_tokens`
- `audit_logs`
- `system_configurations`

### MongoDB

Use MongoDB for append-heavy, semi-structured, high-volume event streams where schema changes are common and document access is more important than relational joins.

Collections that belong in MongoDB:

- `rider_location_logs`
- `order_tracking_events`
- `notification_event_logs`
- `rider_activity_stream`
- `support_conversation_events`
- `device_event_logs`

## Real-time strategy

Recommended production setup:

- SSE for rider order-request stream and lightweight admin live boards
- WebSocket option later if two-way presence/chat volume becomes high
- Redis Pub/Sub or Streams as the fan-out layer between dispatch decisions and connected riders

## Background jobs

Recommended workers:

- order-request timeout and reassignment worker
- OTP resend / cleanup worker
- push notification delivery worker
- payout processing worker
- earnings aggregation worker
- performance score rollup worker
- stale location cleanup worker
- support SLA escalation worker

## Scaling notes

- Use idempotency keys for assignment and payout APIs.
- Use row-level locks for assignment and wallet mutation paths.
- Partition high-volume timeline and audit tables by month if volume grows.
- Cache system configuration and rider live state in Redis.
- Move location ingest to a write-optimized ingestion worker if update throughput becomes high.
