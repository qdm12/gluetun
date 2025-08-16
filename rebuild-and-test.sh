#!/bin/bash

# Gluetun Development: Rebuild and Test Script
# Usage: ./rebuild-and-test.sh

set -e  # Exit on any error

# Get the directory where this script is located
SCRIPT_DIR="$HOME/GITHUB/gluetun"
cd "$SCRIPT_DIR"

echo "🔧 Building Gluetun Docker image..."
docker build -t gluetun-dev .

echo "🚀 Restarting containers with updated image..."
cd test
docker compose down
docker compose up -d

echo "⏳ Waiting for containers to start..."
sleep 10

echo "🧪 Testing WireGuard routing fix..."

echo "📋 Route table 51820:"
docker exec gluetun-custom ip route show table 51820

echo "📋 WireGuard interface status:"
docker exec gluetun-custom ip addr show wg0

echo "📋 Firewall rules (fwmark 51820):"
docker exec gluetun-custom ip rule show | grep 51820 || echo "No fwmark rules found"

echo "📋 Policy routing setup:"
docker exec gluetun-custom ip rule show

echo "📋 Comparing with wg-quick setup:"
echo "Checking for missing wg-quick rules:"
echo "1. suppress_prefixlength rule:"
docker exec gluetun-custom ip rule show | grep "suppress_prefixlength" || echo "❌ MISSING: suppress_prefixlength 0 rule"

echo "2. src_valid_mark sysctl:"
docker exec gluetun-custom sysctl net.ipv4.conf.all.src_valid_mark || echo "❌ MISSING: src_valid_mark sysctl"

echo "3. All sysctl routing settings:"
docker exec gluetun-custom sysctl -a 2>/dev/null | grep -E "(src_valid_mark|rp_filter)" || echo "No routing sysctls found"

echo "📋 NAT table rules:"
docker exec gluetun-custom iptables -t nat -L -n -v

echo "📋 Detailed NAT POSTROUTING chain:"
docker exec gluetun-custom iptables -t nat -L POSTROUTING -n -v --line-numbers

echo "📋 Packet flow test:"
echo "Testing if traffic is getting marked and routed correctly..."
docker exec gluetun-test-client timeout 5 wget -qO- --timeout=3 1.1.1.1 &
sleep 2
echo "Checking wg0 interface traffic during test:"
docker exec gluetun-custom cat /proc/net/dev | grep wg0
wait

echo "🌐 Testing VPN connectivity..."

echo "📋 Current DNS configuration:"
docker exec gluetun-custom cat /etc/resolv.conf

echo "📋 VPN endpoint connectivity:"
VPN_ENDPOINT=$(docker exec gluetun-custom grep -o '[0-9]\+\.[0-9]\+\.[0-9]\+\.[0-9]\+:[0-9]\+' /gluetun/custom.conf | head -1 | cut -d: -f1)
echo "Testing ping to VPN server: $VPN_ENDPOINT"
docker exec gluetun-custom ping -c 2 -W 2 $VPN_ENDPOINT || echo "VPN server unreachable"

echo "📋 Testing with different DNS servers:"
echo "Google DNS (8.8.8.8):"
timeout 15 docker exec gluetun-test-client nslookup ifconfig.me 8.8.8.8 || echo "❌ Google DNS failed"

echo "Cloudflare DNS (1.1.1.1):"
timeout 15 docker exec gluetun-test-client nslookup ifconfig.me 1.1.1.1 || echo "❌ Cloudflare DNS failed"

echo "📋 Direct IP test (bypassing DNS):"
timeout 15 docker exec gluetun-test-client wget -qO- --timeout=2 1.1.1.1 || echo "❌ Direct IP failed"

echo "📋 DNS resolution test (default):"
docker exec gluetun-test-client nslookup ifconfig.me || echo "❌ DNS failed"

echo "📋 Public IP test:"
timeout 30 docker exec gluetun-test-client wget -qO- --timeout=5 ifconfig.me || echo "❌ Connection failed"

echo "✅ Test complete!"
echo ""
echo "Expected results:"
echo "- Route table: 'default dev wg0'"
echo "- WireGuard: 'fwmark: 0xca6c' (51820 in hex)"
echo "- Public IP: Should be ProtonVPN server IP"
