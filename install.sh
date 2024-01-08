#!/bin/bash

# Check if the script is run as root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

# Check if the script is run on Linux
os=$(uname -s)
if [ "$os" != "Linux" ]; then
    echo "This script is only supported on Linux."
    echo "Please build the binary from source, or use go install"
    exit 1
fi

architecture=$(uname -m)

# Determine the download URL for the latest release
if [ "$architecture" = "x86_64" ] || [ "$architecture" = "amd64" ]; then
    download_url="https://github.com/majermarci/hrt/releases/latest/download/hrt-linux-amd64"
elif [ "$architecture" = "aarch64" ] || [ "$architecture" = "arm64" ]; then
    download_url="https://github.com/majermarci/hrt/releases/latest/download/hrt-linux-arm64"
else
    echo "Unsupported architecture, no downloadable release for '$architecture'"
    echo "Please build the binary from source, or use go install"
    exit 1
fi

# Download the file and check if it's successful
if ! curl -fsSL "$download_url" -o /usr/local/bin/hrt; then
    echo "Failed to download file"
    exit 1
fi

# Make the file executable and check if it's successful
if ! chmod +x /usr/local/bin/hrt; then
    echo "Failed to make file executable"
    exit 1
fi

echo "hrt installed to /usr/local/bin/hrt"
