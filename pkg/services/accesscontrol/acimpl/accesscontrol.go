package acimpl

import (
	"context"
	"errors"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/metrics"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
)

var _ accesscontrol.AccessControl = new(AccessControl)

func ProvideAccessControl(cfg *setting.Cfg) *AccessControl {
	logger := log.New("accesscontrol")
	return &AccessControl{
		cfg, logger, accesscontrol.NewResolvers(logger),
	}
}

type AccessControl struct {
	cfg       *setting.Cfg
	log       log.Logger
	resolvers accesscontrol.Resolvers
}

func (a *AccessControl) Evaluate(ctx context.Context, user *user.SignedInUser, evaluator accesscontrol.Evaluator) (bool, error) {
	timer := prometheus.NewTimer(metrics.MAccessEvaluationsSummary)
	defer timer.ObserveDuration()
	metrics.MAccessEvaluationCount.Inc()

	if !verifyPermissions(user) {
		a.log.Warn("no permissions set for user", "userID", user.UserID, "orgID", user.OrgID, "login", user.Login)
		return false, nil
	}
	// Test evaluation without scope resolver first, this will prevent 403 for wildcard scopes when resource does not exist
	if evaluator.Evaluate(user.Permissions[user.OrgID]) {
		return true, nil
	}

	resolvedEvaluator, err := evaluator.MutateScopes(ctx, a.resolvers.GetScopeAttributeMutator(user.OrgID))
	if err != nil {
		if errors.Is(err, accesscontrol.ErrResolverNotFound) {
			return false, nil
		}
		return false, err
	}

	return resolvedEvaluator.Evaluate(user.Permissions[user.OrgID]), nil
}

type Checker func(scopes ...string) bool

func (a *AccessControl) Checker(ctx context.Context, user *user.SignedInUser, action string, prefixes ...string) Checker {
	if !verifyPermissions(user) {
		return func(scope ...string) bool { return false }
	}

	permissions, ok := user.Permissions[user.OrgID]
	if !ok {
		return func(scope ...string) bool { return false }
	}

	checkers := map[string]Checker{}
	scopes, ok := permissions[action]
	if !ok {
		checkers[action] = func(scope ...string) bool { return false }
	}

	wildcards := accesscontrol.WildcardsFromPrefix()
	lookup := make(map[string]bool, len(scopes)-1)
	for _, s := range scopes {
		if wildcards.Contains(s) {
			return func(scope ...string) bool { return true }
		}
		lookup[s] = true
	}

	return func(scopes ...string) bool {
		for _, s := range scopes {
			if lookup[s] {
				return true
			}
		}
		return false
	}
}

func (a *AccessControl) RegisterScopeAttributeResolver(prefix string, resolver accesscontrol.ScopeAttributeResolver) {
	a.resolvers.AddScopeAttributeResolver(prefix, resolver)
}

func (a *AccessControl) IsDisabled() bool {
	return accesscontrol.IsDisabled(a.cfg)
}

func verifyPermissions(u *user.SignedInUser) bool {
	return u.Permissions != nil || u.Permissions[u.OrgID] != nil
}
