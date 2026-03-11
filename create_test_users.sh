#!/bin/bash

echo "Creating test users for the Writer Portal..."

# Function to create a user
create_user() {
    local username=$1
    local email=$2
    local password=$3
    local role=$4

    echo "Creating $role: $username"

    curl -s -X POST http://localhost:5004/signup \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "username=$username&email=$email&password=$password&role=$role" \
        | echo "Response: $(cat)"

    echo ""
}

# Create test users
echo "=========================================="
echo "Creating Test Users"
echo "=========================================="

create_user "testwriter" "writer@example.com" "password123" "writer"
create_user "testeditor" "editor@example.com" "password123" "editor"
create_user "demowriter" "demo@example.com" "password123" "writer"

echo "=========================================="
echo "Test Users Created!"
echo ""
echo "You can now test with these accounts:"
echo ""
echo "WRITER ACCOUNTS:"
echo "  Username: testwriter"
echo "  Password: password123"
echo "  Email: writer@example.com"
echo ""
echo "  Username: demowriter"
echo "  Password: password123"
echo "  Email: demo@example.com"
echo ""
echo "EDITOR ACCOUNT:"
echo "  Username: testeditor"
echo "  Password: password123"
echo "  Email: editor@example.com"
echo ""
echo "NEW USER (no account needed):"
echo "  Just visit the site without logging in"
echo "  Can view all pages including character pages"
echo "  Cannot access /writing or /payment"
echo ""
echo "Login at: http://localhost:5004/login"
echo "=========================================="