package jio

import (
	"fmt"
	"math"
)

func Number() *NumberSchema {
	return &NumberSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(NumberSchema)

type NumberSchema struct {
	rules []func(*Context)
}

func (n *NumberSchema) Required() *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
	return n
}

func (n *NumberSchema) Optional() *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
	return n
}

func (n *NumberSchema) Default(value float64) *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
	return n
}

func (n *NumberSchema) Valid(values ...float64) *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) {
		var isValid bool
		for _, v := range values {
			if v == ctx.Value {
				isValid = true
				break
			}
		}
		if !isValid {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not in %v", ctx.FieldPath(), ctx.Value, values))
		}
	})
	return n
}

func (n *NumberSchema) Min(min float64) *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) {
		if ctx.Value.(float64) < min {
			ctx.Abort(fmt.Errorf("field `%s` value %v less than %v", ctx.FieldPath(), ctx.Value, min))
		}
	})
	return n
}

func (n *NumberSchema) Max(max float64) *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) {
		if ctx.Value.(float64) > max {
			ctx.Abort(fmt.Errorf("field `%s` value %v exceeded %v", ctx.FieldPath(), ctx.Value, max))
		}
	})
	return n
}

func (n *NumberSchema) Ceil() *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) {
		ctx.Value = math.Ceil(ctx.Value.(float64))
	})
	return n
}

func (n *NumberSchema) Floor() *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) {
		ctx.Value = math.Floor(ctx.Value.(float64))
	})
	return n
}

func (n *NumberSchema) Round() *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) {
		ctx.Value = math.Floor(ctx.Value.(float64) + 0.5)
	})
	return n
}

func (n *NumberSchema) Validate(ctx *Context) {
	for _, rule := range n.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
}
