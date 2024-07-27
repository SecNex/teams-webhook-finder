# Programmname
PROGRAM_NAME="webhook-finder"
VERSION="0.1.2"
PROGRAM="$PROGRAM_NAME-$VERSION"

# Windows (64-Bit und 32-Bit)
export GOOS=windows
export GOARCH=amd64
go build -o "build/${PROGRAM}_windows_amd64.exe"
export GOARCH=386
go build -o "build/${PROGRAM}_windows_386.exe"

# Linux (64-Bit, 32-Bit und ARM64)
export GOOS=linux
export GOARCH=amd64
go build -o "build/${PROGRAM}_linux_amd64"
export GOARCH=386
go build -o "build/${PROGRAM}_linux_386"
export GOARCH=arm64
go build -o "build/${PROGRAM}_linux_arm64"

# macOS (64-Bit und ARM64)
export GOOS=darwin
export GOARCH=amd64
go build -o "build/${PROGRAM}_darwin_amd64"
export GOARCH=arm64
go build -o "build/${PROGRAM}_darwin_arm64"

# Zurücksetzen der Umgebungsvariablen
unset GOOS
unset GOARCH

echo "Kompilierung abgeschlossen."