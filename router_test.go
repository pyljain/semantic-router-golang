package semantic_router

import (
	"semantic_router/pkg/encoder"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeRoutingDecision(t *testing.T) {

	techRoutes := []*Route{
		{
			Name: "kubernetes",
			Metadata: RouteMetadata{
				Utterances: []string{
					"how do I deploy a pod in Kubernetes?",
					"what is a Kubernetes namespace?",
					"how can I check the status of my Kubernetes cluster?",
					"how do I create a service?",
					"why is my pod in CrashLoopBackOff state?",
					"how do I create a network policy",
				},
			},
		},
		{
			Name: "terraform",
			Metadata: RouteMetadata{
				Utterances: []string{
					"how do I initialize a Terraform configuration?",
					"what is the purpose of a Terraform provider?",
					"how can I import existing infrastructure into Terraform?",
					"why is my Terraform apply failing?",
					"how do I use Terraform with AWS?",
				},
			},
		},
		{
			Name: "linux",
			Metadata: RouteMetadata{
				Utterances: []string{
					"how do I list all files in a directory in Linux?",
					"what is the difference between a symlink and a hard link?",
					"how do I check disk usage in Linux?",
					"how can I change file permissions in Linux?",
					"why is my shell script not executing?",
				},
			},
		},
		{
			Name: "cicd",
			Metadata: RouteMetadata{
				Utterances: []string{
					"how do I set up a CI/CD pipeline with Jenkins?",
					"what is continuous integration?",
					"how can I deploy code using GitHub Actions?",
					"why is my CI build failing?",
					"how do I implement continuous deployment?",
				},
			},
		},
	}

	hrRoutes := []*Route{
		{
			Name: "hr_general",
			Metadata: RouteMetadata{
				Utterances: []string{
					"what are the company policies on remote work?",
					"how can I apply for leave?",
					"what is the process for performance reviews?",
					"how do I report workplace harassment?",
					"what are the benefits of working here?",
				},
			},
		},
		{
			Name: "hr_payroll",
			Metadata: RouteMetadata{
				Utterances: []string{
					"when is payday?",
					"how can I view my payslip?",
					"what is the process for salary revision?",
					"how do I update my bank details for payroll?",
					"who do I contact for payroll issues?",
				},
			},
		},
	}

	tt := []struct {
		description    string
		routes         []*Route
		userInput      string
		expectedResult string
	}{
		{
			description:    "Identifies the expected route - Kubernetes",
			routes:         techRoutes,
			userInput:      "How do I create a service of type loadbalancer?",
			expectedResult: "kubernetes",
		},
		{
			description:    "Identifies the expected route - Terraform",
			routes:         techRoutes,
			userInput:      "Why is my Terraform apply failing?",
			expectedResult: "terraform",
		},
		{
			description:    "Identifies the expected route - Linux",
			routes:         techRoutes,
			userInput:      "How can I change file permissions in Linux?",
			expectedResult: "linux",
		},
		{
			description:    "Identifies the expected route - CI/CD",
			routes:         techRoutes,
			userInput:      "Why is the pipeline taking so long?",
			expectedResult: "cicd",
		},
		{
			description:    "Identifies the expected route - HR General",
			routes:         hrRoutes,
			userInput:      "What are the company policies on remote work?",
			expectedResult: "hr_general",
		},
		{
			description:    "Identifies the expected route - HR Payroll",
			routes:         hrRoutes,
			userInput:      "When is payday?",
			expectedResult: "hr_payroll",
		},
	}

	enc, err := encoder.NewOpenAIEncoder("text-embedding-3-small")
	require.NoError(t, err)

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			r, err := New(enc, test.routes)
			require.NoError(t, err)

			result, err := r.MakeRoutingDecision(test.userInput)
			require.NoError(t, err)

			require.Equal(t, test.expectedResult, result)
		})
	}
}
