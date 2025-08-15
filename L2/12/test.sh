#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Setup ---
echo "--- Building the grep program ---"
go build -o mygrep main.go
echo "Build complete. Executable 'mygrep' created."
echo ""

# Create a test file
cat > input.txt << EOL
Hello World
hello world
HELLO WORLD
Go is fun
Golang is great
This is a test file
Another line
Case matters
Another world
Final line
EOL

echo "--- Starting tests ---"

# --- Test Cases ---

echo ""
echo "1. Basic search for 'world'"
./mygrep world input.txt

echo ""
echo "2. Case-insensitive search for 'world' (-i)"
./mygrep -i world input.txt

echo ""
echo "3. Inverted search for 'world' with line numbers (-v -n)"
./mygrep -v -n world input.txt

echo ""
echo "4. Context After: 2 lines after 'fun' (-A 2)"
./mygrep -A 2 fun input.txt

echo ""
echo "5. Context Before: 1 line before 'line' (-B 1)"
./mygrep -B 1 line input.txt

echo ""
echo "6. Context Around: 1 line around 'great' (-C 1)"
./mygrep -C 1 great input.txt

echo ""
echo "7. Count matches for 'world' (case-insensitive) (-c -i)"
./mygrep -c -i world input.txt

echo ""
echo "8. Fixed string search for 'is fun' (-F)"
./mygrep -F "is fun" input.txt
