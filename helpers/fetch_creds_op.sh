#!/bin/bash

# ---- TENANT ACCOUNT ----
TENANT_USERNAME=$(op item get $TENANT_ACCOUNT_ID --vault "$OP_VAULT" --fields TENANT_USERNAME --reveal)
TENANT_PRIVATE_SSH_KEY=$(op item get $TENANT_ACCOUNT_ID --vault "$OP_VAULT" --fields TENANT_PRIVATE_SSH_KEY | tr -d '"')
TENANT_ADDRESS=$(op item get $TENANT_ACCOUNT_ID --vault "$OP_VAULT" --fields TENANT_ADDRESS)

# ---- BACKEND ACCOUNT ----
LOG_LEVEL=$(op item get $SERVICE_ACCOUNT_ID --vault "$OP_VAULT" --fields LOG_LEVEL --reveal)


MIST_BACKEND_APP_URL=$(op item get $SERVICE_ACCOUNT_ID --vault "$OP_VAULT" --fields MIST_BACKEND_APP_URL --reveal)

MIST_API_JWT_SECRET_KEY=$(op item get $SERVICE_ACCOUNT_ID --vault "$OP_VAULT" --fields MIST_API_JWT_SECRET_KEY --reveal)
MIST_API_JWT_AUDIENCE=$(op item get $SERVICE_ACCOUNT_ID --vault "$OP_VAULT" --fields MIST_API_JWT_AUDIENCE --reveal)
MIST_API_JWT_ISSUER=$(op item get $SERVICE_ACCOUNT_ID --vault "$OP_VAULT" --fields MIST_API_JWT_ISSUER --reveal)

APP_PORT=$(op item get $SERVICE_ACCOUNT_ID --vault "$OP_VAULT" --fields APP_PORT --reveal)

# Define file paths
KEY_FILE="key.pem"
INVENTORY_FILE="ansible/inventory/hosts.ini"

# Create ansible inventory directory and temporary file
mkdir -p ansible/inventory
touch $INVENTORY_FILE
echo "[mist-service]" >> "$INVENTORY_FILE"
echo $TENANT_ADDRESS >> "$INVENTORY_FILE"

# Create PRIVATE_SSH_KEY temporary file
touch $KEY_FILE
chmod 600 "$KEY_FILE"
echo -e "$TENANT_PRIVATE_SSH_KEY" >> "$KEY_FILE"


# Create tmporary environment variables file
touch .tmpenvs
echo "export TENANT_USERNAME=$TENANT_USERNAME" >> ".tmpenvs"

echo "export APP_PORT=$APP_PORT" >> ".tmpenvs"
echo "export MIST_BACKEND_APP_URL=\"$MIST_BACKEND_APP_URL\"" >> ".tmpenvs"

echo "export MIST_API_JWT_SECRET_KEY=\"$MIST_API_JWT_SECRET_KEY\"" >> ".tmpenvs"
echo "export MIST_API_JWT_AUDIENCE=\"$MIST_API_JWT_AUDIENCE\"" >> ".tmpenvs"
echo "export MIST_API_JWT_ISSUER=\"$MIST_API_JWT_ISSUER\"" >> ".tmpenvs"

