// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1
//
// Package tokens implements creation, signing and verification of short‑lived
// URL tokens used in the IAM subsystem (invite, email verification, password reset).
//
// Design Overview
// A token type (e.g. OrganizationInviteToken, VerificationToken, ResetToken) embeds
// SigningInfo. SigningInfo supplies:
//   - ExpiresAt (UTC timestamp) for short‑lived validity windows
//   - Nonce (random bytes) to ensure each issued token instance has unique
//     input material even if logical data is identical.
//
// Signing Process
//  1. The full token struct (including the freshly generated Nonce & ExpiresAt)
//     is msgpack marshalled.
//  2. A per‑token random HMAC key (keyLength bytes) is generated.
//  3. HMAC‑SHA256(data, key) is computed. The resulting digest is base64 (RawURL) encoded;
//     this becomes the "signature" returned to clients (e.g. embedded in a link).
//  4. The server stores (or otherwise associates) the returned secret which concatenates
//     Nonce||Key (nonceLength + keyLength bytes). No secret -> cannot verify.
//
// Verification Process
//   - Reconstruct the token struct (including setting the Nonce from secret) and marshal it
//     identically, recompute the HMAC and constant‑time compare with the provided signature.
//   - Reject if: expired, malformed secret length, missing required logical fields, or signature mismatch.
//
// Expiration Semantics
// Tokens are considered expired strictly when ExpiresAt.Before(time.Now()). A token expiring
// at the exact call time (== now) is treated as expired (consistent with tests). Negative
// durations to NewSigningInfo intentionally yield immediately expired tokens (used in tests).
//
// Security Notes
//   - Each token uses an independent random HMAC key; compromise does not cascade.
//   - HMAC comparison uses hmac.Equal (constant‑time) through the verifyData helper.
//   - Secrets must be stored securely server‑side and zeroed when no longer needed if
//     long‑term memory disclosure is a concern.
//   - msgpack is chosen for compact, deterministic binary representation.
//
// Migration / Extension
// For new token types: define struct embedding SigningInfo, provide constructor that calls
// NewSigningInfo with domain‑appropriate TTL, a Sign method that marshals & calls signData,
// and a Verify method mirroring existing examples OR implement URLToken and reuse VerifyToken.
//
// Testing Guidance
// Tests should cover: successful Sign/Verify, tampering (field modification), invalid secret
// length, expired instances, and required field validation.
package tokens
