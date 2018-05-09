// 重试通用接口

package retryer

func New() *Retryer {
	return &Retryer{}
}

// 重试器，如果policy为空，则不重试
// 不保证线程安全
type Retryer struct {
	policy Policy
}

func (r *Retryer) SetPolicy(policy Policy) *Retryer {
	r.policy = policy
	return r
}

func (r *Retryer) ShouldRetry() bool {
	if r.policy == nil {
		return false
	} else {
		return r.policy.Continue()
	}
}
