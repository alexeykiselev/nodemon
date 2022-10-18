package criteria

import (
	"github.com/pkg/errors"
	"nodemon/pkg/entities"
)

type BaseTargetCriterionOptions struct {
	Threshold int64
}

type BaseTargetCriterion struct {
	opts *BaseTargetCriterionOptions
}

func NewBaseTargetCriterion(opts *BaseTargetCriterionOptions) (*BaseTargetCriterion, error) {
	if opts == nil {
		return nil, errors.New("provided <nil> base target criterion options")
	}
	return &BaseTargetCriterion{opts: opts}, nil
}

func (c *BaseTargetCriterion) Analyze(alerts chan<- entities.Alert, timestamp int64, statements entities.NodeStatements) {
	var baseTargetValues []entities.BaseTargetValue
	baseTarget := mostFrequentBaseTarget(statements)

	for _, nodeStatement := range statements {
		baseTargetValues = append(baseTargetValues, entities.BaseTargetValue{Node: nodeStatement.Node, BaseTarget: nodeStatement.BaseTarget})
	}

	if baseTarget > c.opts.Threshold || baseTarget == 0 {
		alerts <- &entities.BaseTargetAlert{Timestamp: timestamp, BaseTargetValues: baseTargetValues, Threshold: c.opts.Threshold}
	}
}

func mostFrequentBaseTarget(statements entities.NodeStatements) int64 {
	m := make(map[int64]int)
	var max int
	var maxKey int64
	for _, s := range statements {
		v := s.BaseTarget
		m[v]++
		if m[v] > max {
			max = m[v]
			maxKey = v
		}
	}
	return maxKey
}
