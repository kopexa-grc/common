// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package summarizer

import (
	"context"
	"os"
	"testing"
)

func TestAzureOpenAIIntegration(t *testing.T) {
	// Skip if Azure OpenAI credentials are not available
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	deployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT")

	if apiKey == "" || endpoint == "" || deployment == "" {
		t.Skip("Azure OpenAI credentials not set, skipping integration test")
	}

	t.Run("Azure OpenAI with real credentials", func(t *testing.T) {
		config := NewConfig(
			WithType(TypeLlm),
			WithOpenAI("gpt-4", apiKey,
				WithURL(endpoint),
				WithOption("deployment", deployment),
				WithOption("api_type", "azure"),
				WithOption("api_version", "2023-05-15"),
			),
		)

		client, err := New(config)
		if err != nil {
			t.Fatalf("Failed to create Azure OpenAI summarizer client: %v", err)
		}

		// Test 1: Incident Report Email
		incidentReport := `Subject: Security Incident Report - Unauthorized Access Attempt
From: security@company.com
To: management@company.com
Date: 2024-01-15 14:30:00

Dear Management Team,

I am writing to report a security incident that occurred on January 15, 2024, at approximately 13:45 CET. Our intrusion detection system (IDS) detected multiple failed login attempts to our customer database server from an external IP address (192.168.1.100).

The incident details are as follows:
- Time of detection: 13:45 CET
- Source IP: 192.168.1.100 (external)
- Target system: Customer database server (prod-db-01)
- Attack vector: Brute force login attempts
- Number of failed attempts: 47
- Duration: 15 minutes
- Status: Blocked by firewall

Our security team immediately responded by:
1. Blocking the source IP address
2. Reviewing system logs for any successful access
3. Checking for data exfiltration
4. Notifying the IT operations team

Initial investigation shows no successful unauthorized access or data compromise. However, we recommend:
- Implementing additional rate limiting
- Enabling two-factor authentication for all database access
- Conducting a security audit of our authentication systems
- Reviewing our incident response procedures

The incident has been resolved, and we will provide a detailed post-incident report within 48 hours.

Best regards,
Security Team
Company GmbH`

		// Test 2: Data Subject Request (DSR)
		dsrRequest := `Data Subject Request - Right to Access

Request ID: DSR-2024-001
Date: 2024-01-15
Data Subject: Max Mustermann
Email: max.mustermann@email.com
Phone: +49 123 456789

Dear Data Protection Officer,

I am exercising my right to access under Article 15 of the General Data Protection Regulation (GDPR). I hereby request access to all personal data that your organization processes about me.

Specifically, I request information about:

1. The categories of personal data being processed
2. The purposes of the processing
3. The recipients or categories of recipients to whom the personal data have been or will be disclosed
4. The envisaged period for which the personal data will be stored
5. The existence of the right to request rectification or erasure of personal data
6. The right to lodge a complaint with a supervisory authority
7. Information about the source of the data if not collected from the data subject
8. The existence of automated decision-making, including profiling

I have been a customer of your services since March 2022 and have used the following services:
- Online shopping platform (customer account)
- Newsletter subscription
- Customer support interactions
- Mobile application usage

Please provide this information in a commonly used electronic format. I understand that you have one month to respond to this request, and I may be required to provide additional identification if necessary.

I look forward to your prompt response.

Sincerely,
Max Mustermann
ID Verification: German National ID - 123456789`

		// Test 3: Technical Documentation
		technicalDoc := `API Documentation: User Management Service

Overview:
The User Management Service is a microservice responsible for handling user authentication, authorization, and profile management across our platform. It provides RESTful APIs for user registration, login, profile updates, and account deletion.

Architecture:
The service is built using Go 1.21 and follows a clean architecture pattern with the following layers:
- HTTP handlers for request/response handling
- Business logic layer for core functionality
- Data access layer for database operations
- External service integration layer

Key Features:
1. User Registration: Supports email/password and OAuth2 authentication
2. Profile Management: CRUD operations for user profiles
3. Role-based Access Control: Hierarchical permission system
4. Session Management: JWT-based token authentication
5. Audit Logging: Comprehensive activity tracking
6. Rate Limiting: Protection against abuse

Database Schema:
Users table:
- id (UUID, primary key)
- email (VARCHAR(255), unique)
- password_hash (VARCHAR(255))
- first_name (VARCHAR(100))
- last_name (VARCHAR(100))
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
- is_active (BOOLEAN)

Security Considerations:
- Passwords are hashed using bcrypt with cost factor 12
- JWT tokens have 24-hour expiration
- All endpoints require HTTPS
- Input validation on all parameters
- SQL injection protection through prepared statements

Performance Metrics:
- Average response time: 45ms
- 99th percentile: 120ms
- Throughput: 1000 requests/second
- Uptime: 99.9%

Deployment:
The service is deployed using Docker containers orchestrated by Kubernetes. We use:
- Horizontal Pod Autoscaler for load management
- ConfigMaps for environment configuration
- Secrets for sensitive data
- Ingress controllers for external access

Monitoring and Logging:
- Prometheus metrics for performance monitoring
- ELK stack for log aggregation
- Grafana dashboards for visualization
- AlertManager for incident notification

Future Enhancements:
- Multi-factor authentication support
- Social login integration
- Advanced analytics dashboard
- API versioning strategy`

		deGerman := `Betreff: Zusammenfassung Datenschutz-Folgenabschätzung

Sehr geehrte Damen und Herren,

im Rahmen unserer neuen Cloud-basierten Plattform wurde eine Datenschutz-Folgenabschätzung (DSFA) durchgeführt. Ziel war es, die Risiken für die Rechte und Freiheiten der betroffenen Personen zu identifizieren und geeignete technische sowie organisatorische Maßnahmen zu definieren.

Wesentliche Ergebnisse:
- Die Verarbeitung personenbezogener Daten erfolgt ausschließlich verschlüsselt.
- Zugriff erhalten nur autorisierte Mitarbeitende nach dem Need-to-know-Prinzip.
- Ein Löschkonzept ist implementiert, das die DSGVO-Anforderungen erfüllt.
- Die Datenübermittlung an Drittländer findet nicht statt.

Empfohlene Maßnahmen:
1. Regelmäßige Überprüfung der Zugriffskontrollen
2. Durchführung von Datenschutzschulungen
3. Jährliche Aktualisierung der DSFA

Für Rückfragen stehe ich Ihnen gerne zur Verfügung.

Mit freundlichen Grüßen
Datenschutzbeauftragte`

		tests := []struct {
			name string
			text string
		}{
			{
				name: "Incident Report Email",
				text: incidentReport,
			},
			{
				name: "Data Subject Request",
				text: dsrRequest,
			},
			{
				name: "Technical Documentation",
				text: technicalDoc,
			},
			{
				name: "DSFA Deutsch",
				text: deGerman,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()
				result, err := client.Summarize(ctx, tt.text)

				if err != nil {
					t.Fatalf("Failed to summarize %s: %v", tt.name, err)
				}

				if result == "" {
					t.Fatalf("Expected non-empty summary for %s", tt.name)
				}

				t.Logf("=== %s ===", tt.name)
				t.Logf("Original length: %d characters", len(tt.text))
				t.Logf("Summary length: %d characters", len(result))
				t.Logf("Compression ratio: %.1f%%", float64(len(result))/float64(len(tt.text))*100)
				t.Logf("Summary: %s", result)
				t.Logf("==================")

				// Verify the summary is shorter than the original
				if len(result) >= len(tt.text) {
					t.Errorf("Summary should be shorter than original text for %s", tt.name)
				}

				// Verify the summary is not too short (should capture key information)
				minLength := int(float64(len(tt.text)) * 0.1)
				if len(result) < minLength {
					t.Errorf("Summary seems too short for %s (less than 10%% of original)", tt.name)
				}
			})
		}
	})
}
