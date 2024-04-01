package opa

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/open-policy-agent/opa/rego"
	"go.uber.org/zap"
	"golang-api-hexagonal/core/domain"
)

// PolicyService policy service
type PolicyService struct {
	log    *zap.SugaredLogger
	policy rego.PreparedEvalQuery
}

// NewPolicyService new policy service
func NewPolicyService(regoFilePath string, log *zap.SugaredLogger) *PolicyService {
	r := rego.New(rego.Query("x = data.golangapitemplate.authz.allow"), rego.Load([]string{regoFilePath}, nil))

	eval, err := r.PrepareForEval(context.Background())
	if err != nil {
		log.Fatalf("Failed to load the policy file: %v", err)
		return nil
	}

	log.Infof("OPA policy loaded")
	return &PolicyService{
		log:    log,
		policy: eval,
	}
}

// EvaluateApiPolicy evaluate the api policy roles
func (p *PolicyService) EvaluateApiPolicy(ctx context.Context, claims domain.AuthClaims, operation, ownerUserName string) bool {
	traceID := ctx.Value(middleware.RequestIDKey)

	token := Token{
		Username: claims.Username,
		Roles:    []string{claims.Role},
	}
	data := EntityData{
		Type:  operation,
		Owner: ownerUserName,
	}
	input := PolicyInput{
		Token:      token,
		EntityData: data,
	}

	result, err := p.policy.Eval(ctx, rego.EvalInput(input))
	if err != nil || len(result) == 0 {
		p.log.With("traceId", traceID).Warnf("Policy evaluation failed: %s", err)
		return false
	} else if _, ok := result[0].Bindings["x"].(bool); !ok {
		p.log.With("traceId", traceID).Warn("Policy evaluation failed")
		return false
	} else {
		p.log.With("traceId", traceID).Infof("Policy result: %v", result)
		return result[0].Bindings["x"].(bool)
	}
}
