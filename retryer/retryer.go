// 重试通用接口

package retryer

import (
	"sync"
)

var (
	retryPool = sync.Pool{
		New: func() interface{} {
			return &Retryer{}
		},
	}
)

func New() *Retryer {
	r := retryPool.Get().(*Retryer)
	r.policy = nil
	return r
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
