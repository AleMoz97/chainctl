#!/bin/sh

set -eu

OWNER="AleMoz97"
REPO="chainctl"
BIN_NAME="chainctl"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="latest"

usage() {
  cat <<EOF
Usage: install.sh [--version <tag>] [--install-dir <dir>]

Examples:
  sh install.sh
  sh install.sh --version v0.0.1
  INSTALL_DIR="\$HOME/.local/bin" sh install.sh
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      VERSION="$2"
      shift 2
      ;;
    --install-dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

need_cmd uname
need_cmd mktemp
need_cmd tar
need_cmd curl

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

checksum_cmd=""
if command -v sha256sum >/dev/null 2>&1; then
  checksum_cmd="sha256sum"
elif command -v shasum >/dev/null 2>&1; then
  checksum_cmd="shasum -a 256"
else
  echo "Missing required command: sha256sum or shasum" >&2
  exit 1
fi

case "$os" in
  linux) os="linux" ;;
  darwin) os="darwin" ;;
  mingw*|msys*|cygwin*) os="windows" ;;
  *)
    echo "Unsupported OS: $os" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "Unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

if [ "$VERSION" = "latest" ]; then
  release_url="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
  VERSION="$(curl -fsSL "$release_url" | sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)"
  if [ -z "$VERSION" ]; then
    echo "Unable to resolve latest release version" >&2
    exit 1
  fi
fi

archive_ext="tar.gz"
if [ "$os" = "windows" ]; then
  archive_ext="zip"
  need_cmd unzip
fi

archive_name="${BIN_NAME}_${VERSION}_${os}_${arch}.${archive_ext}"
checksums_name="checksums.txt"
base_url="https://github.com/$OWNER/$REPO/releases/download/$VERSION"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT INT TERM

curl -fsSL "$base_url/$archive_name" -o "$tmp_dir/$archive_name"
curl -fsSL "$base_url/$checksums_name" -o "$tmp_dir/$checksums_name"

expected_checksum="$(grep "  $archive_name\$" "$tmp_dir/$checksums_name" | awk '{print $1}' | head -n 1)"
if [ -z "$expected_checksum" ]; then
  echo "Checksum for $archive_name not found" >&2
  exit 1
fi

if [ "$checksum_cmd" = "sha256sum" ]; then
  actual_checksum="$(sha256sum "$tmp_dir/$archive_name" | awk '{print $1}')"
else
  actual_checksum="$(shasum -a 256 "$tmp_dir/$archive_name" | awk '{print $1}')"
fi

if [ "$expected_checksum" != "$actual_checksum" ]; then
  echo "Checksum verification failed for $archive_name" >&2
  exit 1
fi

if [ "$os" = "windows" ]; then
  unzip -q "$tmp_dir/$archive_name" -d "$tmp_dir/unpack"
else
  tar -xzf "$tmp_dir/$archive_name" -C "$tmp_dir"
fi

mkdir -p "$INSTALL_DIR"

if [ "$os" = "windows" ]; then
  install_target="$INSTALL_DIR/${BIN_NAME}.exe"
  src_path="$(find "$tmp_dir/unpack" -type f -name "${BIN_NAME}.exe" | head -n 1)"
else
  install_target="$INSTALL_DIR/$BIN_NAME"
  src_path="$(find "$tmp_dir" -type f -name "$BIN_NAME" | head -n 1)"
fi

if [ -z "${src_path:-}" ]; then
  echo "Downloaded archive does not contain $BIN_NAME" >&2
  exit 1
fi

install -m 0755 "$src_path" "$install_target"

echo "$BIN_NAME installed in $install_target"
echo "Run: $BIN_NAME --help"
