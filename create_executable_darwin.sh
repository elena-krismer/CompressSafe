#!/bin/bash

# Define variables
GO_SCRIPT="compresssafe.go"
EXECUTABLE="compresssafe"
INSTALL_DIR="/usr/local/bin"

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Go is not installed. Please install Go and try again."
    exit 1
fi

# Compile the Go script
echo "Compiling the Go script..."
go build -o "$EXECUTABLE" "$GO_SCRIPT"
if [ $? -ne 0 ]; then
    echo "Failed to compile the Go script."
    exit 1
fi
echo "Compilation successful."

# Move the executable to /usr/local/bin
echo "Moving the executable to $INSTALL_DIR..."
sudo mv "$EXECUTABLE" "$INSTALL_DIR"
if [ $? -ne 0 ]; then
    echo "Failed to move the executable to $INSTALL_DIR."
    exit 1
fi
echo "Executable moved successfully."

# Ensure /usr/local/bin is in the PATH
echo "Checking if $INSTALL_DIR is in the PATH..."
SHELL_PROFILE=""
if [ -n "$ZSH_VERSION" ]; then
    SHELL_PROFILE="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_PROFILE="$HOME/.bash_profile"
else
    echo "Unsupported shell. Please add $INSTALL_DIR to your PATH manually."
    exit 1
fi

if ! grep -q "$INSTALL_DIR" "$SHELL_PROFILE"; then
    echo "Adding $INSTALL_DIR to the PATH in $SHELL_PROFILE..."
    echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$SHELL_PROFILE"
    source "$SHELL_PROFILE"
    echo "$INSTALL_DIR added to the PATH."
else
    echo "$INSTALL_DIR is already in the PATH."
fi

echo "Setup completed successfully. You can now run the CLI tool from any directory using '$EXECUTABLE'."
