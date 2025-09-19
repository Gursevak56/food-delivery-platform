#!/bin/bash

BASE_URL="http://localhost:8085"
CONTENT_TYPE="Content-Type: application/json"

echo "=== Create Restaurant ==="
CREATE_REST=$(curl -s -X POST $BASE_URL/restaurants \
-H "$CONTENT_TYPE" \
-d '{
  "ownerId": 1,
  "name": "Pizza Hut",
  "description": "Famous for pizzas",
  "cuisine_type": "Italian",
  "phone": "9876543210",
  "email": "pizzahut@example.com",
  "street_address": "123 Main Street",
  "city": "New York",
  "state": "NY",
  "postal_code": "10001",
  "latitude": 40.7128,
  "longitude": -74.0060,
  "image_url": "http://example.com/pizza.png",
  "delivery_fee": 2.5,
  "minimum_order": 10.0,
  "delivery_time_min": 20,
  "delivery_time_max": 40
}')
echo $CREATE_REST | jq .
REST_ID=$(echo $CREATE_REST | jq -r '.restaurant_id')

echo -e "\n=== Get Restaurant by ID ($REST_ID) ==="
curl -s -X GET $BASE_URL/restaurants/$REST_ID | jq .

echo -e "\n=== Get All Restaurants ==="
curl -s -X GET $BASE_URL/restaurants | jq .

echo -e "\n=== Update Restaurant ($REST_ID) ==="
curl -s -X PUT $BASE_URL/restaurants/$REST_ID \
-H "$CONTENT_TYPE" \
-d '{
  "ownerId": 1,
  "name": "Pizza Hut Updated",
  "description": "Updated description",
  "cuisine_type": "Italian",
  "phone": "9876543210",
  "email": "pizzahut@example.com",
  "street_address": "123 Main Street",
  "city": "New York",
  "state": "NY",
  "postal_code": "10001",
  "latitude": 40.7128,
  "longitude": -74.0060,
  "image_url": "http://example.com/pizza.png",
  "delivery_fee": 3.0,
  "minimum_order": 12.0,
  "delivery_time_min": 25,
  "delivery_time_max": 45
}' | jq .

echo -e "\n=== Create Category ==="
CREATE_CAT=$(curl -s -X POST $BASE_URL/menu/categories \
-H "$CONTENT_TYPE" \
-d '{
  "restaurantId": '"$REST_ID"',
  "name": "Starters",
  "description": "Appetizers and snacks",
  "display_order": 1
}')
echo $CREATE_CAT | jq .
CAT_ID=$(echo $CREATE_CAT | jq -r '.category_id')

echo -e "\n=== Create Menu Item ==="
CREATE_ITEM=$(curl -s -X POST $BASE_URL/menu/items \
-H "$CONTENT_TYPE" \
-d '{
  "restaurantId": '"$REST_ID"',
  "category_id": '"$CAT_ID"',
  "name": "Margherita Pizza",
  "description": "Classic cheese pizza",
  "price": 9.99,
  "image_url": "http://example.com/margherita.png",
  "is_vegetarian": true,
  "is_vegan": false,
  "is_gluten_free": false,
  "calories": 250,
  "preparation_time": 15,
  "is_available": true,
  "display_order": 1
}')
echo $CREATE_ITEM | jq .

echo -e "\n=== Get Menu Items by Restaurant ($REST_ID) ==="
curl -s -X GET $BASE_URL/menu/items/$REST_ID | jq .

echo -e "\n=== Generate QR Code for Restaurant ($REST_ID) ==="
curl -s -X GET $BASE_URL/restaurants/$REST_ID/qr --output qr_$REST_ID.png
echo "QR Code saved as qr_$REST_ID.png"

echo -e "\n=== Delete Restaurant ($REST_ID) ==="
curl -s -X DELETE $BASE_URL/restaurants/$REST_ID | jq .

echo -e "\n=== Done ==="
