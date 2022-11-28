package criteria

import (
	"github.com/pkg/errors"
	"nodemon/pkg/entities"
	"nodemon/pkg/storing/events"
)

type IncompleteCriterionOptions struct {
	Streak                              int
	Depth                               int
	ConsiderPrevUnreachableAsIncomplete bool
}

type IncompleteCriterion struct {
	opts *IncompleteCriterionOptions
	es   *events.Storage
}

func NewIncompleteCriterion(es *events.Storage, opts *IncompleteCriterionOptions) *IncompleteCriterion {
	if opts == nil { // by default
		opts = &IncompleteCriterionOptions{
			Streak:                              3,
			Depth:                               5,
			ConsiderPrevUnreachableAsIncomplete: true,
		}
	}
	return &IncompleteCriterion{opts: opts, es: es}
}

func (c *IncompleteCriterion) Analyze(alerts chan<- entities.Alert, statements entities.NodeStatements) error {
	for _, statement := range statements {
		if err := c.analyzeNode(alerts, statement); err != nil {
			return err
		}
	}
	return nil
}

func (c *IncompleteCriterion) analyzeNode(alerts chan<- entities.Alert, statement entities.NodeStatement) error {
	var (
		streak = 0
		depth  = 0
	)
	err := c.es.ViewStatementsByNodeWithDescendKeys(statement.Node, func(statement *entities.NodeStatement) bool {
		if s := statement.Status; s == entities.Incomplete || (c.opts.ConsiderPrevUnreachableAsIncomplete && s == entities.Unreachable) {
			streak += 1
		} else {
			streak = 0
		}
		depth += 1
		if streak >= c.opts.Streak || depth >= c.opts.Depth {
			return false
		}
		return true
	})
	if err != nil {
		return errors.Wrapf(err, "failed to analyze %q by incomplete criterion", statement.Node)
	}
	if streak >= c.opts.Streak {
		alerts <- &entities.IncompleteAlert{NodeStatement: statement}
	}
	return nil
}