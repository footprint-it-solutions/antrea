#!/usr/bin/env bash
# Usage: check-aws-sso.sh <profile_name> [<sso_target>]

PROFILE=$1
SSO_TARGET=$2

# Attempt to get identity; if it fails, trigger login
if ! aws sts get-caller-identity --profile "$PROFILE" >/dev/null 2>&1; then
    echo "⚠️ AWS SSO session expired for profile [$PROFILE]. Logging in..."
    # --use-device-code is helpful for SSH/Remote sessions
    if [ -n "$SSO_TARGET" ]; then
        # Try login using profile first, then fall back to sso-session if it fails
        aws sso login --profile "$SSO_TARGET" --no-browser --use-device-code 2>/dev/null || aws sso login --sso-session "$SSO_TARGET" --no-browser --use-device-code
    else
        aws sso login --profile "$PROFILE" --no-browser --use-device-code
    fi
else
    echo "✅ AWS SSO session is valid for [$PROFILE]"
fi
