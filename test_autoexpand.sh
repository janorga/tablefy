#!/bin/bash

# Quick test script for auto-expand feature

echo "=== Tablefy Auto-Expand Feature Test ==="
echo ""
echo "Creating test data with helm-like output..."
echo ""

# Create test data
TEST_DATA="NAME	NAMESPACE	REVISION	UPDATED	STATUS	CHART	APP VERSION
argocd	argocd	2	2025-07-29 20:02:02.89383204 +0200 CEST	deployed	argo-cd-8.2.3	v3.0.12
longhorn	longhorn-system	2	2025-08-01 14:13:13.960555922 +0200 CEST	deployed	longhorn-1.9.1	v1.9.1
rke2-canal	kube-system	8	2025-07-29 12:37:00.299221083 +0000 UTC	deployed	rke2-canal-v3.30.1-build2025061101	v3.30.1
sealed-secrets	kube-system	1	2025-11-13 15:26:36.16453789 +0100 CET	deployed	sealed-secrets-2.17.9	0.33.1"

echo "Test data prepared (4 rows)"
echo ""
echo "To test auto-expand:"
echo ""
echo "1. Open a narrow terminal window (around 80-100 columns wide)"
echo "2. Run: cat << 'DATA' | ./bin/tablefy --auto-expand"
echo "3. Navigate with arrow keys to column 4 (UPDATED)"
echo "4. Observe the column expanding to show full timestamps"
echo "5. Move away with arrow keys and watch it shrink back"
echo ""
echo "✅ Ready to test!"
