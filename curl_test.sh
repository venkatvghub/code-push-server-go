#!/bin/bash

# Base URL
BASE_URL="http://localhost:8080"

# Variables to store tokens and IDs
JWT_TOKEN="L4uLBuy7v2X9wMatEDwZFnJLk8RYnPlYV0lcKIZ76uj12S0N732uzg6pJSVHD3Q"
SOURCE_DEPLOYMENT_KEY=""
DEST_DEPLOYMENT_KEY=""
APP_NAME="TestApp"
SOURCE_DEPLOYMENT="Source"
DEST_DEPLOYMENT="Dest"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to check response and print status
check_response() {
    local response="$1"
    local expected_status="$2"
    if echo "$response" | grep -q "$expected_status"; then
        echo -e "${GREEN}SUCCESS${NC}: $3"
    else
        echo -e "${RED}FAILURE${NC}: $3"
        echo "Response: $response"
    fi
}

# Create a test zip file
echo "Creating test.zip for package upload..."
echo "test content" > test.txt
zip test.zip test.txt

# 1. Auth Routes
echo "=== Testing Auth Routes ==="

# Login (POST /auth/login)
echo "Testing POST /auth/login"
response=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "account=aaa@bbb.com&password=123456")
check_response "$response" '"status":"OK"' "Login should succeed"
JWT_TOKEN=$(echo "$response" | grep -o '"tokens":"[^"]*"' | cut -d'"' -f4)

# Logout (POST /auth/logout)
echo "Testing POST /auth/logout"
response=$(curl -s -X POST "$BASE_URL/auth/logout")
check_response "$response" "ok" "Logout should return 'ok'"

# Register (POST /auth/register)
echo "Testing POST /auth/register"
response=$(curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "email=aaa@bbb.com&password=123456")
check_response "$response" '"status":"OK"' "Registration should succeed"

# 2. Index Routes
echo "=== Testing Index Routes ==="

# Index (GET /)
echo "Testing GET /"
response=$(curl -s "$BASE_URL/")
check_response "$response" "CodePushServer" "Index page should render"

# Tokens (GET /tokens)
echo "Testing GET /tokens"
response=$(curl -s "$BASE_URL/tokens")
check_response "$response" "Obtain token" "Tokens page should render"

# Authenticated (GET /authenticated)
echo "Testing GET /authenticated"
response=$(curl -s -H "Authorization: Bearer $JWT_TOKEN" "$BASE_URL/authenticated")
check_response "$response" '"authenticated":true' "Authenticated should return true with valid token"

# Update Check (GET /updateCheck)
echo "Testing GET /updateCheck (will fail without deployment key)"
response=$(curl -s "$BASE_URL/updateCheck?deploymentKey=invalid&appVersion=1.0&label=&packageHash=&clientUniqueId=xyz")
check_response "$response" "Invalid deployment key" "Update check with invalid key should fail"

# Report Status Download (POST /reportStatus/download)
echo "Testing POST /reportStatus/download (will fail without deployment key)"
response=$(curl -s -X POST "$BASE_URL/reportStatus/download" \
    -H "Content-Type: application/json" \
    -d '{"clientUniqueId":"xyz","label":"v1","deploymentKey":"invalid"}')
check_response "$response" '"OK"' "Report status download should return OK even if invalid"

# Report Status Deploy (POST /reportStatus/deploy)
echo "Testing POST /reportStatus/deploy (will fail without deployment key)"
response=$(curl -s -X POST "$BASE_URL/reportStatus/deploy" \
    -H "Content-Type: application/json" \
    -d '{"clientUniqueId":"xyz","label":"v1","deploymentKey":"invalid","status":1}')
check_response "$response" '"OK"' "Report status deploy should return OK even if invalid"

# 3. Users Routes
echo "=== Testing Users Routes ==="

# Change Password (PATCH /users/password)
echo "Testing PATCH /users/password"
response=$(curl -s -X PATCH "$BASE_URL/users/password" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"oldPassword":"123456","newPassword":"newpass123"}')
check_response "$response" '"status":"OK"' "Password change should succeed"

# Reset password back to original for further tests
response=$(curl -s -X PATCH "$BASE_URL/users/password" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"oldPassword":"newpass123","newPassword":"123456"}')
check_response "$response" '"status":"OK"' "Password reset should succeed"

# 4. AccessKeys Routes
echo "=== Testing AccessKeys Routes ==="

# Create Access Key (POST /accessKeys)
echo "Testing POST /accessKeys"
response=$(curl -s -X POST "$BASE_URL/accessKeys" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"createdBy":"test","friendlyName":"TestKey","ttl":86400000,"description":"Test Key","isSession":false}')
check_response "$response" '"accessKey"' "Access key creation should succeed"

# 5. Account Routes
echo "=== Testing Account Routes ==="

# Get Access Keys (GET /account/accessKeys)
echo "Testing GET /account/accessKeys"
response=$(curl -s -H "Authorization: Bearer $JWT_TOKEN" "$BASE_URL/account/accessKeys")
check_response "$response" '"accessKeys"' "Get access keys should succeed"

# 6. Apps Routes
echo "=== Testing Apps Routes ==="

# Add App (POST /apps)
echo "Testing POST /apps"
response=$(curl -s -X POST "$BASE_URL/apps" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"'$APP_NAME'","os":"iOS","platform":"React-Native","manuallyProvisionDeployments":false}')
check_response "$response" '"app"' "App creation should succeed"

# List Collaborators (GET /apps/:appName/collaborators)
echo "Testing GET /apps/$APP_NAME/collaborators"
response=$(curl -s -H "Authorization: Bearer $JWT_TOKEN" "$BASE_URL/apps/$APP_NAME/collaborators")
check_response "$response" '"collaborators"' "List collaborators should succeed"

# Add Collaborator (POST /apps/:appName/collaborators/:email)
echo "Testing POST /apps/$APP_NAME/collaborators/test@example.com"
response=$(curl -s -X POST -H "Authorization: Bearer $JWT_TOKEN" "$BASE_URL/apps/$APP_NAME/collaborators/test@example.com")
check_response "$response" "{}" "Add collaborator should succeed"

# Add Deployment (POST /apps/:appName/deployments)
echo "Testing POST /apps/$APP_NAME/deployments (Source)"
response=$(curl -s -X POST "$BASE_URL/apps/$APP_NAME/deployments" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"'$SOURCE_DEPLOYMENT'"}')
check_response "$response" '"deployment"' "Add Source deployment should succeed"
SOURCE_DEPLOYMENT_KEY=$(echo "$response" | grep -o '"key":"[^"]*"' | cut -d'"' -f4)

echo "Testing POST /apps/$APP_NAME/deployments (Dest)"
response=$(curl -s -X POST "$BASE_URL/apps/$APP_NAME/deployments" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"'$DEST_DEPLOYMENT'"}')
check_response "$response" '"deployment"' "Add Dest deployment should succeed"
DEST_DEPLOYMENT_KEY=$(echo "$response" | grep -o '"key":"[^"]*"' | cut -d'"' -f4)

# Release Package (POST /apps/:appName/deployments/:deploymentName/release)
echo "Testing POST /apps/$APP_NAME/deployments/$SOURCE_DEPLOYMENT/release"
response=$(curl -s -X POST "$BASE_URL/apps/$APP_NAME/deployments/$SOURCE_DEPLOYMENT/release" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -F "file=@test.zip" \
    -F "description=Test release" \
    -F "isMandatory=true")
check_response "$response" '"msg":"succeed"' "Package release should succeed"

# Promote Package (POST /apps/:appName/deployments/promote)
echo "Testing POST /apps/$APP_NAME/deployments/promote"
response=$(curl -s -X POST "$BASE_URL/apps/$APP_NAME/deployments/promote" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"sourceDeploymentName":"'$SOURCE_DEPLOYMENT'","destDeploymentName":"'$DEST_DEPLOYMENT'"}')
check_response "$response" '"package"' "Package promotion should succeed"

# Rollback Package (POST /apps/:appName/deployments/:deploymentName/rollback)
echo "Testing POST /apps/$APP_NAME/deployments/$SOURCE_DEPLOYMENT/rollback"
response=$(curl -s -X POST "$BASE_URL/apps/$APP_NAME/deployments/$SOURCE_DEPLOYMENT/rollback" \
    -H "Authorization: Bearer $JWT_TOKEN")
check_response "$response" '"msg":"ok"' "Rollback without label should succeed"

# Rollback Package with Label (POST /apps/:appName/deployments/:deploymentName/rollback/:label)
echo "Testing POST /apps/$APP_NAME/deployments/$SOURCE_DEPLOYMENT/rollback/v1"
response=$(curl -s -X POST "$BASE_URL/apps/$APP_NAME/deployments/$SOURCE_DEPLOYMENT/rollback/v1" \
    -H "Authorization: Bearer $JWT_TOKEN")
check_response "$response" '"msg":"ok"' "Rollback with label should succeed or fail gracefully"

# Delete App (DELETE /apps/:appName)
echo "Testing DELETE /apps/$APP_NAME"
response=$(curl -s -X DELETE -H "Authorization: Bearer $JWT_TOKEN" "$BASE_URL/apps/$APP_NAME")
check_response "$response" "{}" "App deletion should succeed"

# 7. IndexV1 Routes (v0.1/public/codepush)
echo "=== Testing IndexV1 Routes (/v0.1/public/codepush) ==="

# Update Check (GET /v0.1/public/codepush/update_check)
echo "Testing GET /v0.1/public/codepush/update_check (will fail without deployment key)"
response=$(curl -s "$BASE_URL/v0.1/public/codepush/update_check?deployment_key=invalid&app_version=1.0&label=&package_hash=&client_unique_id=xyz")
check_response "$response" "Invalid deployment key" "Update check with invalid key should fail"

echo "Testing GET /v0.1/public/codepush/update_check with valid deployment key"
response=$(curl -s "$BASE_URL/v0.1/public/codepush/update_check?deployment_key=$SOURCE_DEPLOYMENT_KEY&app_version=1.0&label=&package_hash=&client_unique_id=xyz")
check_response "$response" '"update_info"' "Update check with valid key should succeed"

# Report Status Download (POST /v0.1/public/codepush/report_status/download)
echo "Testing POST /v0.1/public/codepush/report_status/download (will fail without deployment key)"
response=$(curl -s -X POST "$BASE_URL/v0.1/public/codepush/report_status/download" \
    -H "Content-Type: application/json" \
    -d '{"client_unique_id":"xyz","label":"v1","deployment_key":"invalid"}')
check_response "$response" '"OK"' "Report status download should return OK even if invalid"

echo "Testing POST /v0.1/public/codepush/report_status/download with valid deployment key"
response=$(curl -s -X POST "$BASE_URL/v0.1/public/codepush/report_status/download" \
    -H "Content-Type: application/json" \
    -d '{"client_unique_id":"xyz","label":"v1","deployment_key":"'$SOURCE_DEPLOYMENT_KEY'"}')
check_response "$response" '"OK"' "Report status download with valid key should succeed"

# Report Status Deploy (POST /v0.1/public/codepush/report_status/deploy)
echo "Testing POST /v0.1/public/codepush/report_status/deploy (will fail without deployment key)"
response=$(curl -s -X POST "$BASE_URL/v0.1/public/codepush/report_status/deploy" \
    -H "Content-Type: application/json" \
    -d '{"client_unique_id":"xyz","label":"v1","deployment_key":"invalid","status":1}')
check_response "$response" '"OK"' "Report status deploy should return OK even if invalid"

echo "Testing POST /v0.1/public/codepush/report_status/deploy with valid deployment key"
response=$(curl -s -X POST "$BASE_URL/v0.1/public/codepush/report_status/deploy" \
    -H "Content-Type: application/json" \
    -d '{"client_unique_id":"xyz","label":"v1","deployment_key":"'$SOURCE_DEPLOYMENT_KEY'","status":1}')
check_response "$response" '"OK"' "Report status deploy with valid key should succeed"

# Cleanup
echo "Cleaning up test files..."
rm -f test.txt test.zip

echo "=== Testing Complete ==="
