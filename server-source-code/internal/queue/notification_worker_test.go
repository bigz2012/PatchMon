package queue

import (
	"bufio"
	"bytes"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
	"testing"

	"github.com/PatchMon/PatchMon/server-source-code/internal/notifications"
)

// Issue #851: Mattermost incoming webhooks use a "/hooks/<token>" path on an
// arbitrary self-hosted host and accept Slack-compatible payloads.
func TestIsMattermostWebhookURL(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want bool
	}{
		{"mattermost incoming", "https://mattermost.example.com/hooks/abcd1234efgh5678", true},
		{"mattermost with port", "http://mm.local:8065/hooks/token", true},
		{"slack excluded", "https://hooks.slack.com/services/T000/B000/xxxx", false},
		{"generic endpoint", "https://example.com/api/notify", false},
		{"discord", "https://discord.com/api/webhooks/123/abc", false},
		{"garbage", "://not a url", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isMattermostWebhookURL(c.url); got != c.want {
				t.Fatalf("isMattermostWebhookURL(%q) = %v, want %v", c.url, got, c.want)
			}
		})
	}
}

// Issue #817: with TLS explicitly disabled, PlainAuth must be allowed to send
// credentials over a plaintext connection instead of returning
// "unencrypted connection".
func TestUnencryptedAuthStartOverPlaintext(t *testing.T) {
	base := smtp.PlainAuth("", "user", "pass", "relay.internal")

	// PlainAuth alone refuses a non-TLS, non-localhost server.
	if _, _, err := base.Start(&smtp.ServerInfo{Name: "relay.internal", TLS: false}); err == nil {
		t.Fatal("expected PlainAuth to refuse plaintext connection, got nil error")
	}

	// Wrapped, it proceeds.
	wrapped := unencryptedAuth{base}
	if _, _, err := wrapped.Start(&smtp.ServerInfo{Name: "relay.internal", TLS: false}); err != nil {
		t.Fatalf("unencryptedAuth.Start over plaintext: unexpected error %v", err)
	}
}

// Issue #845: the generated HTML body is a single long line. Quoted-printable
// encoding (as sendEmail now applies) must keep every line within the RFC 5322
// 998-octet limit so strict relays do not reject the message.
func TestEmailBodyQuotedPrintableLineLength(t *testing.T) {
	p := notifications.NotificationDeliverPayload{
		EventType: "host_down",
		Severity:  "critical",
		Title:     "Host is down: " + strings.Repeat("very-long-hostname.example.internal ", 40),
		Message:   strings.Repeat("Something went wrong and this needs attention. ", 60),
		Metadata: map[string]interface{}{
			"host_name":         strings.Repeat("node-", 100) + "01",
			"last_update":       "2026-06-20T00:00:00Z",
			"threshold_minutes": "5",
			"app_link":          "https://patchmon.example.com/hosts/" + strings.Repeat("a", 300),
		},
	}

	html := buildEmailHTML(p)
	if !strings.Contains(html, strings.Repeat("a", 50)) {
		t.Fatal("expected long content in HTML body for the test to be meaningful")
	}

	var buf bytes.Buffer
	w := quotedprintable.NewWriter(&buf)
	if _, err := w.Write([]byte(html)); err != nil {
		t.Fatalf("quotedprintable write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("quotedprintable close: %v", err)
	}

	sc := bufio.NewScanner(bytes.NewReader(buf.Bytes()))
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for sc.Scan() {
		if l := len(sc.Bytes()); l > 998 {
			t.Fatalf("encoded line length %d exceeds RFC 5322 limit of 998 octets", l)
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan: %v", err)
	}
}
