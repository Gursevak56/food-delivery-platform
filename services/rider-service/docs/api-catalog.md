# API Catalog

## Auth

- `POST /v1/auth/rider/login`
- `POST /v1/auth/rider/otp/send`
- `POST /v1/auth/rider/otp/verify`
- `POST /v1/auth/refresh-token`
- `POST /v1/auth/logout`
- `POST /v1/auth/logout-all`
- `POST /v1/auth/forgot-password`
- `POST /v1/auth/reset-password`
- `GET /v1/auth/me`

## Rider Profile

- `GET /v1/riders/me`
- `PUT /v1/riders/me`
- `PUT /v1/riders/me/photo`
- `PUT /v1/riders/me/vehicle`
- `PUT /v1/riders/me/documents`
- `GET /v1/riders/me/documents`
- `PUT /v1/riders/me/bank-account`
- `GET /v1/riders/me/status`

## Availability and Shift

- `POST /v1/riders/me/go-online`
- `POST /v1/riders/me/go-offline`
- `POST /v1/riders/me/break/start`
- `POST /v1/riders/me/break/end`
- `POST /v1/riders/me/shift/start`
- `POST /v1/riders/me/shift/end`
- `GET /v1/riders/me/shift/today`
- `GET /v1/riders/me/shift/history`

## Assignment and Order Requests

- `GET /v1/riders/me/order-requests`
- `GET /v1/riders/me/order-requests/{id}`
- `POST /v1/riders/me/order-requests/{id}/accept`
- `POST /v1/riders/me/order-requests/{id}/reject`
- `GET /v1/riders/me/active-order`
- `GET /v1/riders/me/orders/assigned`
- `GET /v1/riders/me/orders/history`

## Active Delivery Lifecycle

- `GET /v1/riders/me/orders/{orderId}`
- `POST /v1/riders/me/orders/{orderId}/arrived-at-restaurant`
- `POST /v1/riders/me/orders/{orderId}/picked-up`
- `POST /v1/riders/me/orders/{orderId}/start-delivery`
- `POST /v1/riders/me/orders/{orderId}/arrived-at-customer`
- `POST /v1/riders/me/orders/{orderId}/deliver`
- `POST /v1/riders/me/orders/{orderId}/failed`
- `POST /v1/riders/me/orders/{orderId}/cancel-request`
- `GET /v1/riders/me/orders/{orderId}/timeline`

## OTP and Verification

- `POST /v1/riders/me/orders/{orderId}/verify-pickup-otp`
- `POST /v1/riders/me/orders/{orderId}/verify-delivery-otp`
- `POST /v1/orders/{orderId}/resend-delivery-otp`
- `POST /v1/orders/{orderId}/resend-pickup-otp`

## Location and Tracking

- `POST /v1/riders/me/location/update`
- `POST /v1/riders/me/location/bulk-update`
- `GET /v1/riders/me/location/latest`
- `GET /v1/orders/{orderId}/tracking`
- `GET /v1/riders/me/routes/{orderId}`

## Earnings and Wallet

- `GET /v1/riders/me/earnings/today`
- `GET /v1/riders/me/earnings/weekly`
- `GET /v1/riders/me/earnings/monthly`
- `GET /v1/riders/me/earnings/summary`
- `GET /v1/riders/me/earnings/history`
- `GET /v1/riders/me/incentives`
- `GET /v1/riders/me/bonus-history`
- `GET /v1/riders/me/wallet`
- `GET /v1/riders/me/wallet/transactions`
- `GET /v1/riders/me/payouts`
- `GET /v1/riders/me/payouts/{id}`
- `POST /v1/riders/me/payouts/request`
- `GET /v1/riders/me/bank-account`

## Ratings, Notifications, Support

- `GET /v1/riders/me/ratings/summary`
- `GET /v1/riders/me/reviews`
- `GET /v1/riders/me/performance-score`
- `GET /v1/riders/me/notifications`
- `PUT /v1/riders/me/notifications/{id}/read`
- `PUT /v1/riders/me/notifications/read-all`
- `POST /v1/riders/me/device-token`
- `DELETE /v1/riders/me/device-token`
- `POST /v1/riders/me/support-tickets`
- `GET /v1/riders/me/support-tickets`
- `GET /v1/riders/me/support-tickets/{id}`
- `POST /v1/riders/me/support-tickets/{id}/reply`

## Admin and Dispatcher

- `GET /v1/admin/riders`
- `GET /v1/admin/riders/{id}`
- `POST /v1/admin/riders`
- `PUT /v1/admin/riders/{id}`
- `PUT /v1/admin/riders/{id}/status`
- `GET /v1/admin/orders/unassigned`
- `POST /v1/admin/orders/{orderId}/assign-rider`
- `POST /v1/admin/orders/{orderId}/reassign-rider`
- `GET /v1/admin/orders/live`
- `GET /v1/admin/riders/live-status`
- `GET /v1/admin/analytics/riders`
- `GET /v1/admin/config`
- `PUT /v1/admin/config`
- `GET /v1/admin/payouts`
- `POST /v1/admin/payouts/{id}/approve`
- `POST /v1/admin/payouts/{id}/reject`

## Response standard

Success envelope:

```json
{
  "success": true,
  "message": "Order accepted successfully",
  "data": {},
  "meta": {}
}
```

Error envelope:

```json
{
  "success": false,
  "message": "Invalid OTP",
  "error_code": "INVALID_OTP",
  "errors": {
    "otp": "numeric"
  }
}
```
