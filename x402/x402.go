// Package x402 is the opt-in adapter that plugs the official x402 v2 client
// (github.com/x402-foundation/x402/go/v2) into the Nebula SDK.
//
// The x402 surface serves the 37 fixed-price GET operations listed in
// docs/API-COVERAGE.md on https://nebula-api.hiddensystems.ai/api/v1/x402.
// Each request pays real USDC on Base mainnet (eip155:8453) with the EVM
// "exact" scheme. Payments are irreversible: always configure a SpendPolicy
// and inspect PreviewOffer output before authorizing payments.
//
//	signer, err := x402.NewSignerFromPrivateKeyHex(os.Getenv("NEBULA_X402_KEY_HEX"))
//	policy := x402.SpendPolicy{MaxUSDCPerPayment: "0.05"}
//	httpClient, err := x402.NewHTTPClient(nil, signer, policy)
//	client, err := x402.NewClient(signer, policy)
//	assets, err := client.Lookup.ListAssets(ctx, &nebula.ListAssetsRequest{...})
//
// The SDK never logs private keys or signed payment payloads.
package x402

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	nebula "github.com/tigusigalpa/nebula-go"
	x402v2 "github.com/x402-foundation/x402/go/v2"
	x402http "github.com/x402-foundation/x402/go/v2/http"
	evmmech "github.com/x402-foundation/x402/go/v2/mechanisms/evm"
	evmexact "github.com/x402-foundation/x402/go/v2/mechanisms/evm/exact/client"
	evmsigner "github.com/x402-foundation/x402/go/v2/signers/evm"
)

const (
	// Network is the CAIP-2 network Nebula settles on: Base mainnet.
	Network = "eip155:8453"
	// USDCAsset is the USDC contract on Base mainnet.
	USDCAsset = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
	// USDCDecimals is the atomic precision of USDC.
	USDCDecimals = 6
	// Host is the production x402 host.
	Host = "nebula-api.hiddensystems.ai"
	// Scheme is the only scheme Nebula accepts.
	Scheme = "exact"
)

// Signer is the official client-side EVM signer interface. Any implementation
// works; NewSignerFromPrivateKeyHex wraps the standard private-key signer.
type Signer = evmmech.ClientEvmSigner

// NewSignerFromPrivateKeyHex creates an EVM signer from a hex-encoded
// private key using the official x402 signer. The key authorizes USDC spends;
// keep it in a secret manager, never in source or logs.
func NewSignerFromPrivateKeyHex(privateKeyHex string) (Signer, error) {
	return evmsigner.NewClientSignerFromPrivateKey(privateKeyHex)
}

// SpendPolicy is the explicit pre-signing policy every payment offer is
// checked against. The offer's resource URL, network, USDC asset, recipient
// and amount must all match before the official client signs anything.
type SpendPolicy struct {
	// MaxUSDCPerPayment caps one payment in USD units, e.g. "0.05" for five
	// cents. Required: empty is rejected so a missing cap can never mean
	// "unlimited". The same cap is also enforced by the official
	// SpendControls mechanism.
	MaxUSDCPerPayment string
	// AllowedRecipients restricts the payTo addresses that may be paid.
	// Empty means the recipient is not restricted beyond the other checks;
	// list the addresses Nebula publishes for defense in depth.
	AllowedRecipients []string
	// AllowedHosts restricts which resource hosts may trigger payment.
	// Defaults to Host when empty.
	AllowedHosts []string
	// RequireExactScheme rejects offers whose scheme is not "exact" (the only
	// scheme Nebula accepts). Always true in practice.
	RequireExactScheme bool
}

func (p SpendPolicy) hosts() []string {
	if len(p.AllowedHosts) == 0 {
		return []string{Host}
	}
	return p.AllowedHosts
}

// maxAtomic converts MaxUSDCPerPayment to atomic USDC units (6 decimals).
func (p SpendPolicy) maxAtomic() (*big.Int, error) {
	if p.MaxUSDCPerPayment == "" {
		return nil, fmt.Errorf("nebula/x402: MaxUSDCPerPayment is required")
	}
	rat, ok := new(big.Rat).SetString(p.MaxUSDCPerPayment)
	if !ok || rat.Sign() < 0 {
		return nil, fmt.Errorf("nebula/x402: invalid MaxUSDCPerPayment %q", p.MaxUSDCPerPayment)
	}
	micros := new(big.Rat).Mul(rat, big.NewRat(1_000_000, 1))
	if !micros.IsInt() {
		return nil, fmt.Errorf("nebula/x402: MaxUSDCPerPayment %q exceeds 6 decimals", p.MaxUSDCPerPayment)
	}
	return micros.Num(), nil
}

// PaymentOffer is the decoded PAYMENT-REQUIRED offer presented for review
// before any signature. Amount is in atomic USDC units.
type PaymentOffer struct {
	ResourceURL       string
	Network           string
	Scheme            string
	Asset             string
	Amount            *big.Int
	PayTo             string
	MaxTimeoutSeconds int
	Raw               json.RawMessage
}

// Check evaluates the offer against the policy and returns a descriptive
// error when it must not be signed.
func (p SpendPolicy) Check(offer *PaymentOffer) error {
	if err := p.checkResource(offer.ResourceURL); err != nil {
		return err
	}
	return p.checkRequirements(offer)
}

// checkRequirements validates everything except the resource URL: scheme,
// network, USDC asset, recipient, and the amount cap.
func (p SpendPolicy) checkRequirements(offer *PaymentOffer) error {
	if p.RequireExactScheme && offer.Scheme != Scheme {
		return fmt.Errorf("nebula/x402: offer scheme %q is not %q", offer.Scheme, Scheme)
	}
	if offer.Network != Network {
		return fmt.Errorf("nebula/x402: offer network %q is not %q", offer.Network, Network)
	}
	if !strings.EqualFold(offer.Asset, USDCAsset) {
		return fmt.Errorf("nebula/x402: offer asset %q is not USDC on Base (%s)", offer.Asset, USDCAsset)
	}
	if err := p.checkRecipient(offer.PayTo); err != nil {
		return err
	}
	maxAtomic, err := p.maxAtomic()
	if err != nil {
		return err
	}
	if offer.Amount == nil || offer.Amount.Sign() < 0 {
		return fmt.Errorf("nebula/x402: offer has no amount")
	}
	if offer.Amount.Cmp(maxAtomic) > 0 {
		return fmt.Errorf("nebula/x402: offer amount %s micro-USDC exceeds policy cap %s",
			offer.Amount.String(), p.MaxUSDCPerPayment)
	}
	return nil
}

func (p SpendPolicy) checkResource(resourceURL string) error {
	if resourceURL == "" {
		return fmt.Errorf("nebula/x402: offer carries no resource URL")
	}
	for _, h := range p.hosts() {
		if strings.Contains(resourceURL, "://"+h+"/") || strings.HasSuffix(resourceURL, "://"+h) {
			return nil
		}
	}
	return fmt.Errorf("nebula/x402: offer resource %q is not on an allowed host", resourceURL)
}

func (p SpendPolicy) checkRecipient(payTo string) error {
	if len(p.AllowedRecipients) == 0 {
		return nil
	}
	for _, r := range p.AllowedRecipients {
		if strings.EqualFold(r, payTo) {
			return nil
		}
	}
	return fmt.Errorf("nebula/x402: offer recipient %q is not in AllowedRecipients", payTo)
}

// requirementsPolicy adapts SpendPolicy to the official PaymentPolicy hook:
// requirements that fail the check are filtered out, which makes payment
// creation fail safely instead of signing. The resource URL itself is
// verified separately in the OnPaymentRequired hook, where it is visible.
func (p SpendPolicy) requirementsPolicy() x402v2.PaymentPolicy {
	return func(requirements []x402v2.PaymentRequirementsView) []x402v2.PaymentRequirementsView {
		out := make([]x402v2.PaymentRequirementsView, 0, len(requirements))
		for _, rv := range requirements {
			amount, _ := new(big.Int).SetString(rv.GetAmount(), 10)
			offer := &PaymentOffer{
				Network: rv.GetNetwork(),
				Scheme:  rv.GetScheme(),
				Asset:   rv.GetAsset(),
				Amount:  amount,
				PayTo:   rv.GetPayTo(),
			}
			if p.checkRequirements(offer) == nil {
				out = append(out, rv)
			}
		}
		return out
	}
}

// NewHTTPClient wraps base with the official x402 payment transport:
// 402 offers are checked against policy, signed by the official EVM exact
// scheme, and receipts land in response metadata via PAYMENT-RESPONSE.
//
// The returned *http.Client plugs straight into nebula.New /
// nebula.NewX402Client. Spend controls stay enabled; DisableSpendControls is
// never called.
func NewHTTPClient(base *http.Client, signer Signer, policy SpendPolicy) (*http.Client, error) {
	if signer == nil {
		return nil, fmt.Errorf("nebula/x402: signer is required")
	}
	if _, err := policy.maxAtomic(); err != nil {
		return nil, err
	}
	if base == nil {
		base = http.DefaultClient
	}
	x402c := x402v2.Newx402Client(
		x402v2.WithSpendControls(x402v2.SpendControls{
			MaxAmountPerPayment: policy.MaxUSDCPerPayment,
		}),
		x402v2.WithPolicy(policy.requirementsPolicy()),
	)
	x402c.Register(x402v2.Network(Network), evmexact.NewExactEvmScheme(signer, nil))
	hc := x402http.Newx402HTTPClient(x402c)
	hc.OnPaymentRequired(func(ctx context.Context, required x402v2.PaymentRequired, requestURL string) (*x402http.PaymentRequiredHookResult, error) {
		resource := requestURL
		if required.Resource != nil && strings.Contains(required.Resource.URL, "://") {
			resource = required.Resource.URL
		}
		if err := policy.checkResource(resource); err != nil {
			return nil, err
		}
		return &x402http.PaymentRequiredHookResult{}, nil
	})
	return x402http.WrapHTTPClientWithPayment(base, hc), nil
}

// NewClient builds a ready Nebula client bound to the x402 surface with the
// payment transport wired in. Eligible operations behave exactly like the
// account API; excluded operations return ErrNotX402Eligible locally.
func NewClient(signer Signer, policy SpendPolicy, opts ...nebula.Option) (*nebula.Client, error) {
	httpClient, err := NewHTTPClient(nil, signer, policy)
	if err != nil {
		return nil, err
	}
	return nebula.NewX402Client(httpClient, opts...)
}

// PreviewOffer performs the request without paying and decodes the
// PAYMENT-REQUIRED offer so callers can inspect price, network, asset and
// recipient before authorizing anything. It never signs.
//
// Pass httpClient nil to use a bare client; the server must respond 402 with
// the offer header for the decode to succeed.
func PreviewOffer(ctx context.Context, httpClient *http.Client, method, fullURL string) (*PaymentOffer, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	offer, err := DecodePaymentRequired(resp.Header)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

// DecodePaymentRequired decodes the base64 PAYMENT-REQUIRED header into a
// reviewable offer. The first accepted requirement is surfaced; Raw keeps
// the full payload.
func DecodePaymentRequired(h http.Header) (*PaymentOffer, error) {
	raw := h.Get("PAYMENT-REQUIRED")
	if raw == "" {
		return nil, fmt.Errorf("nebula/x402: no PAYMENT-REQUIRED header (status must be 402 with an offer)")
	}
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("nebula/x402: invalid PAYMENT-REQUIRED header: %w", err)
	}
	var pr x402v2.PaymentRequired
	if err := json.Unmarshal(decoded, &pr); err != nil {
		return nil, fmt.Errorf("nebula/x402: malformed payment offer: %w", err)
	}
	offer := &PaymentOffer{Raw: json.RawMessage(decoded)}
	if pr.Resource != nil {
		offer.ResourceURL = pr.Resource.URL
	}
	if len(pr.Accepts) > 0 {
		a := pr.Accepts[0]
		offer.Scheme = a.Scheme
		offer.Network = a.Network
		offer.Asset = a.Asset
		offer.PayTo = a.PayTo
		offer.MaxTimeoutSeconds = a.MaxTimeoutSeconds
		offer.Amount, _ = new(big.Int).SetString(a.Amount, 10)
	}
	return offer, nil
}

// DecodePaymentResponse decodes the PAYMENT-RESPONSE settlement receipt
// header returned after a paid request. Callers should retain it for
// reconciliation.
func DecodePaymentResponse(h http.Header) (json.RawMessage, error) {
	raw := h.Get("PAYMENT-RESPONSE")
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("nebula/x402: invalid PAYMENT-RESPONSE header: %w", err)
	}
	return json.RawMessage(decoded), nil
}
