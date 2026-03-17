# Database Design

## PostgreSQL tables and purpose

### Identity and RBAC

- `users`: core identity, login fields, user lifecycle, verification flags
- `roles`: role catalog
- `user_roles`: many-to-many role mapping
- `refresh_tokens`: refresh sessions and device metadata

### Rider operations

- `riders`: rider operational profile and performance counters
- `rider_documents`: KYC and compliance docs
- `rider_vehicle_details`: operational vehicle metadata
- `rider_bank_accounts`: payout destination
- `rider_restaurant_assignments`: which restaurants a rider can serve
- `rider_shifts`: shift windows and active minutes
- `rider_breaks`: intra-shift breaks
- `rider_locations_latest`: latest authoritative rider point

### Delivery and dispatch

- `restaurants`: restaurant dimension
- `orders`: delivery order header
- `order_items`: line items snapshot
- `delivery_assignments`: rider offer / accept / reject / timeout state
- `delivery_status_history`: audit trail of delivery state transitions
- `order_otps`: pickup and delivery verification OTPs

### Finance

- `wallets`: rider ledger balance
- `wallet_transactions`: immutable wallet movements
- `rider_earnings`: order-level earnings breakdown
- `payout_requests`: payout lifecycle
- `ratings_reviews`: customer or restaurant review entries

### Communication and support

- `notifications`: latest user-visible notification records
- `support_tickets`: support case header
- `support_ticket_messages`: conversation thread in relational form for quick case lookup

### Governance

- `audit_logs`: critical action trail
- `system_configurations`: runtime configurable dispatch and payout knobs

## MongoDB collections and purpose

- `rider_location_logs`: raw GPS ingestion, frequent writes, optional metadata expansion
- `order_tracking_events`: public tracking and internal movement stream
- `notification_event_logs`: provider payloads, retry traces, delivery receipts
- `rider_activity_stream`: online/offline/break/app-open/action events
- `support_conversation_events`: free-form support event stream, attachments, bot/system inserts
- `device_event_logs`: app version, crash, battery, network, and device telemetry

## Core index strategy

### PostgreSQL

- `users(email)`, `users(phone)` unique partial indexes on active users
- `riders(user_id)` unique
- `delivery_assignments(rider_id, status, decision_deadline_at)`
- `orders(status, restaurant_id, created_at desc)`
- `delivery_status_history(order_id, created_at)`
- `order_otps(order_id, purpose, created_at desc)`
- `rider_shifts(rider_id, started_at desc)`
- `rider_locations_latest(rider_id)` primary key or unique
- `wallet_transactions(wallet_id, created_at desc)`
- `rider_earnings(rider_id, created_at desc)`
- `payout_requests(rider_id, status, requested_at desc)`
- `notifications(user_id, is_read, created_at desc)`
- `support_tickets(rider_id, status, updated_at desc)`
- `audit_logs(actor_id, created_at desc)`

### MongoDB

- `rider_location_logs`: compound `(rider_id, recorded_at desc)` with TTL option if raw logs have retention
- `order_tracking_events`: `(order_id, event_time desc)`
- `notification_event_logs`: `(user_id, created_at desc)` and `(provider_message_id)`
- `rider_activity_stream`: `(rider_id, event_time desc)`
- `support_conversation_events`: `(ticket_id, created_at asc)`
- `device_event_logs`: `(user_id, created_at desc)` and `(device_id, created_at desc)`

## State models

### Delivery statuses

- `CREATED`
- `READY_FOR_ASSIGNMENT`
- `ASSIGNED`
- `ACCEPTED`
- `REACHED_RESTAURANT`
- `PICKUP_VERIFIED`
- `PICKED_UP`
- `ON_THE_WAY`
- `REACHED_CUSTOMER`
- `DELIVERY_VERIFIED`
- `DELIVERED`
- `FAILED`
- `CANCELLED`

### Payout statuses

- `PENDING`
- `APPROVED`
- `PROCESSING`
- `PAID`
- `REJECTED`

### Soft delete guidance

Use `deleted_at` on:

- `users`
- `riders`
- `rider_bank_accounts`
- `rider_vehicle_details`
- `restaurants`
- `notifications`
- `support_tickets`

High-volume ledger, assignment, audit, and timeline tables should stay append-only instead of soft-deleted.
