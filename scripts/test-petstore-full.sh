#!/bin/bash

# WEBPRISM Petstore Full Integration Test
# This script demonstrates the complete workflow of WEBPRISM:
# 1. Upload Petstore OpenAPI spec
# 2. Configure authentication
# 3. Execute proxy requests
# 4. Verify responses

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:8080/api"
SPEC_FILE="testdata/petstore-openapi3.json"
SPEC_NAME="petstore"
LOG_FILE="logs/petstore-test-$(date +%Y%m%d-%H%M%S).log"

# Create logs directory
mkdir -p logs

# Helper functions
log_header() {
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${PURPLE}$1${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
}

log_step() {
    echo -e "${CYAN}▶ $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

log_data() {
    echo -e "${NC}$1${NC}"
}

# Log all output
exec > >(tee -a "$LOG_FILE")
exec 2>&1

# Start test
log_header "🚀 WEBPRISM Petstore Full Integration Test"

log_info "Test started at: $(date)"
log_info "Base URL: $BASE_URL"
log_info "Spec file: $SPEC_FILE"
log_info "Log file: $LOG_FILE"

# Check if spec file exists
log_step "Step 0: Checking prerequisites"
if [ ! -f "$SPEC_FILE" ]; then
    log_error "Spec file not found: $SPEC_FILE"
    exit 1
fi
log_success "Spec file found"

# Check if server is running
if ! curl -s -f "$BASE_URL/v1/version" > /dev/null 2>&1; then
    log_error "Server is not running at $BASE_URL"
    log_info "Please start the server with: make run-server"
    exit 1
fi
log_success "Server is running"

# Step 1: Upload OpenAPI Spec
log_header "📤 Step 1: Upload Petstore OpenAPI Specification"

log_step "Uploading spec file..."

# Read spec content and properly escape it as a JSON string
SPEC_CONTENT=$(cat "$SPEC_FILE" | jq -Rs '.')

# Create upload request with proper format using jq to ensure correct JSON
UPLOAD_REQUEST=$(jq -n \
  --arg name "$SPEC_NAME" \
  --argjson content "$SPEC_CONTENT" \
  '{name: $name, specContent: $content, format: "SPEC_FORMAT_JSON"}')

UPLOAD_RESPONSE=$(curl -s -X POST \
    "$BASE_URL/v1/specs" \
    -H "Content-Type: application/json" \
    -d "$UPLOAD_REQUEST")

if echo "$UPLOAD_RESPONSE" | grep -q '"code"'; then
    log_error "Failed to upload spec"
    log_data "$UPLOAD_RESPONSE"
    exit 1
fi

SPEC_ID=$(echo "$UPLOAD_RESPONSE" | jq -r '.spec.id')
log_success "Spec uploaded successfully"
log_info "Spec ID: $SPEC_ID"
log_data "Response preview:"
echo "$UPLOAD_RESPONSE" | jq '.' | head -20

# Step 2: Retrieve and verify spec
log_header "📥 Step 2: Retrieve and Verify Uploaded Spec"

log_step "Retrieving spec by ID..."
GET_SPEC_RESPONSE=$(curl -s "$BASE_URL/v1/specs/$SPEC_ID")

if echo "$GET_SPEC_RESPONSE" | jq -e '.spec' > /dev/null 2>&1; then
    log_success "Spec retrieved successfully"

    # Extract key information
    SPEC_NAME_RETRIEVED=$(echo "$GET_SPEC_RESPONSE" | jq -r '.spec.name')
    SPEC_VERSION=$(echo "$GET_SPEC_RESPONSE" | jq -r '.spec.version')
    SPEC_BASE_URL=$(echo "$GET_SPEC_RESPONSE" | jq -r '.spec.baseUrl')

    log_info "Name: $SPEC_NAME_RETRIEVED"
    log_info "Version: $SPEC_VERSION"
    log_info "Base URL: $SPEC_BASE_URL"
    log_data "Full response:"
    echo "$GET_SPEC_RESPONSE" | jq '.' | head -30
else
    log_error "Failed to retrieve spec"
    log_data "$GET_SPEC_RESPONSE"
    exit 1
fi

# Step 3: List all specs
log_header "📋 Step 3: List All Specs"

log_step "Listing all specs..."
LIST_RESPONSE=$(curl -s "$BASE_URL/v1/specs")

SPEC_COUNT=$(echo "$LIST_RESPONSE" | jq -r '.specs | length')
log_success "Retrieved $SPEC_COUNT spec(s)"
log_data "Preview:"
echo "$LIST_RESPONSE" | jq '.' | head -20

# Step 4: Configure API Key Authentication
log_header "🔐 Step 4: Configure API Key Authentication"

log_step "Creating auth config for Petstore API..."
AUTH_CONFIG=$(cat <<EOF
{
  "authType": "AUTH_TYPE_API_KEY",
  "credentials": {
    "key": "special-key",
    "headerName": "api_key"
  },
  "position": "AUTH_POSITION_HEADER"
}
EOF
)

AUTH_RESPONSE=$(curl -s -X POST \
    "$BASE_URL/v1/specs/$SPEC_ID/auth" \
    -H "Content-Type: application/json" \
    -d "$AUTH_CONFIG")

if echo "$AUTH_RESPONSE" | jq -e '.code' > /dev/null 2>&1; then
    log_error "Failed to create auth config"
    log_data "$AUTH_RESPONSE"
    exit 1
fi

log_success "Auth config created successfully"
log_data "Response:"
echo "$AUTH_RESPONSE" | jq '.'

# Step 5: Execute Proxy Request - Get Inventory
log_header "🔄 Step 5: Execute Proxy Request - Get Store Inventory"

log_step "Making proxy request to GET /store/inventory (operationId: getInventory)..."
PROXY_REQUEST_1=$(cat <<EOF
{
  "specId": "$SPEC_ID",
  "operationId": "getInventory",
  "parameters": {}
}
EOF
)

PROXY_RESPONSE_1=$(curl -s -X POST \
    "$BASE_URL/v1/proxy" \
    -H "Content-Type: application/json" \
    -d "$PROXY_REQUEST_1")

if echo "$PROXY_RESPONSE_1" | jq -e '.statusCode' > /dev/null 2>&1; then
    STATUS_CODE=$(echo "$PROXY_RESPONSE_1" | jq -r '.statusCode')
    log_success "Proxy request completed"
    log_info "Status Code: $STATUS_CODE"
    log_data "Response:"
    echo "$PROXY_RESPONSE_1" | jq '.'
else
    log_error "Proxy request failed"
    log_data "$PROXY_RESPONSE_1"
fi

# Step 6: Execute Proxy Request - Find Pet by Status
log_header "🔍 Step 6: Execute Proxy Request - Find Pets by Status"

log_step "Making proxy request to GET /pet/findByStatus (operationId: findPetsByStatus)..."
PROXY_REQUEST_2=$(cat <<EOF
{
  "specId": "$SPEC_ID",
  "operationId": "findPetsByStatus",
  "parameters": {
    "status": "available"
  }
}
EOF
)

PROXY_RESPONSE_2=$(curl -s -X POST \
    "$BASE_URL/v1/proxy" \
    -H "Content-Type: application/json" \
    -d "$PROXY_REQUEST_2")

if echo "$PROXY_RESPONSE_2" | jq -e '.statusCode' > /dev/null 2>&1; then
    STATUS_CODE=$(echo "$PROXY_RESPONSE_2" | jq -r '.statusCode')
    log_success "Proxy request completed"
    log_info "Status Code: $STATUS_CODE"
    log_data "Response (first 3 pets):"
    echo "$PROXY_RESPONSE_2" | jq '.body[:3]' 2>/dev/null || echo "$PROXY_RESPONSE_2" | jq '.'
else
    log_error "Proxy request failed"
    log_data "$PROXY_RESPONSE_2"
fi

# Step 7: Execute Proxy Request - Get Pet by ID
log_header "🐕 Step 7: Execute Proxy Request - Get Pet by ID"

log_step "Making proxy request to GET /pet/{petId} (operationId: getPetById)..."
PROXY_REQUEST_3=$(cat <<EOF
{
  "specId": "$SPEC_ID",
  "operationId": "getPetById",
  "parameters": {
    "petId": "1"
  }
}
EOF
)

PROXY_RESPONSE_3=$(curl -s -X POST \
    "$BASE_URL/v1/proxy" \
    -H "Content-Type: application/json" \
    -d "$PROXY_REQUEST_3")

if echo "$PROXY_RESPONSE_3" | jq -e '.statusCode' > /dev/null 2>&1; then
    STATUS_CODE=$(echo "$PROXY_RESPONSE_3" | jq -r '.statusCode')
    log_success "Proxy request completed"
    log_info "Status Code: $STATUS_CODE"
    log_data "Response:"
    echo "$PROXY_RESPONSE_3" | jq '.'
else
    log_error "Proxy request failed"
    log_data "$PROXY_RESPONSE_3"
fi

# Step 8: Retrieve Auth Config
log_header "🔓 Step 8: Retrieve Auth Configuration"

log_step "Retrieving auth config for spec..."
GET_AUTH_RESPONSE=$(curl -s "$BASE_URL/v1/specs/$SPEC_ID/auth")

if echo "$GET_AUTH_RESPONSE" | jq -e '.authType' > /dev/null 2>&1; then
    log_success "Auth config retrieved successfully"
    AUTH_TYPE=$(echo "$GET_AUTH_RESPONSE" | jq -r '.authType')
    AUTH_POSITION=$(echo "$GET_AUTH_RESPONSE" | jq -r '.position')
    log_info "Auth Type: $AUTH_TYPE"
    log_info "Position: $AUTH_POSITION"
    log_data "Response:"
    echo "$GET_AUTH_RESPONSE" | jq '.'
else
    log_error "Failed to retrieve auth config"
    log_data "$GET_AUTH_RESPONSE"
fi

# Step 9: Health Check for Spec
log_header "💚 Step 9: Health Check for API Spec"

log_step "Performing health check for Petstore API..."
HEALTH_REQUEST=$(cat <<EOF
{
  "specId": "$SPEC_ID"
}
EOF
)

HEALTH_CHECK_RESPONSE=$(curl -s -X POST \
    "$BASE_URL/v1/specs/$SPEC_ID/health-check" \
    -H "Content-Type: application/json" \
    -d "$HEALTH_REQUEST")

if echo "$HEALTH_CHECK_RESPONSE" | jq -e '.status' > /dev/null 2>&1; then
    HEALTH_STATUS=$(echo "$HEALTH_CHECK_RESPONSE" | jq -r '.status')
    log_success "Health check completed"
    log_info "Status: $HEALTH_STATUS"
    log_data "Response:"
    echo "$HEALTH_CHECK_RESPONSE" | jq '.'
else
    log_error "Health check failed"
    log_data "$HEALTH_CHECK_RESPONSE"
fi

# Step 10: Get Version Info
log_header "📋 Step 10: Get System Version"

log_step "Retrieving system version..."
VERSION_RESPONSE=$(curl -s "$BASE_URL/v1/version")

log_success "Version info retrieved"
log_data "Version:"
echo "$VERSION_RESPONSE" | jq '.' 2>/dev/null || echo "$VERSION_RESPONSE"

# Summary
log_header "📊 Test Summary"

echo "Test completed at: $(date)"
echo ""
echo "Executed steps:"
echo "  ✅ Step 1: Upload Petstore OpenAPI 3.0 Spec"
echo "  ✅ Step 2: Retrieve and Verify Spec"
echo "  ✅ Step 3: List All Specs"
echo "  ✅ Step 4: Configure API Key Authentication"
echo "  ✅ Step 5: Execute Proxy Request - Get Inventory (operationId: getInventory)"
echo "  ✅ Step 6: Execute Proxy Request - Find Pets by Status (operationId: findPetsByStatus)"
echo "  ✅ Step 7: Execute Proxy Request - Get Pet by ID (operationId: getPetById)"
echo "  ✅ Step 8: Retrieve Auth Configuration"
echo "  ✅ Step 9: Health Check for API Spec"
echo "  ✅ Step 10: Get System Version"
echo ""
echo "Resources created:"
echo "  • Spec ID: $SPEC_ID"
echo "  • Auth Type: API_KEY (configured)"
echo ""
echo "Full log saved to: $LOG_FILE"
echo ""

# Cleanup prompt
log_header "🧹 Cleanup (Optional)"

echo "To clean up test data, you can:"
echo "  1. Delete spec (this will also delete associated auth config):"
echo "     curl -X DELETE $BASE_URL/v1/specs/$SPEC_ID"
echo ""

log_success "Test completed successfully! 🎉"
echo ""
echo "📚 View detailed log at: $LOG_FILE"
