// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/orestonce/ChessGame/ent/dchat"
)

// DChatCreate is the builder for creating a DChat entity.
type DChatCreate struct {
	config
	mutation *DChatMutation
	hooks    []Hook
}

// SetSessionID sets the "session_id" field.
func (dc *DChatCreate) SetSessionID(s string) *DChatCreate {
	dc.mutation.SetSessionID(s)
	return dc
}

// SetUserID sets the "user_id" field.
func (dc *DChatCreate) SetUserID(s string) *DChatCreate {
	dc.mutation.SetUserID(s)
	return dc
}

// SetRoomID sets the "room_id" field.
func (dc *DChatCreate) SetRoomID(s string) *DChatCreate {
	dc.mutation.SetRoomID(s)
	return dc
}

// SetText sets the "text" field.
func (dc *DChatCreate) SetText(s string) *DChatCreate {
	dc.mutation.SetText(s)
	return dc
}

// SetCreateTime sets the "create_time" field.
func (dc *DChatCreate) SetCreateTime(t time.Time) *DChatCreate {
	dc.mutation.SetCreateTime(t)
	return dc
}

// SetID sets the "id" field.
func (dc *DChatCreate) SetID(s string) *DChatCreate {
	dc.mutation.SetID(s)
	return dc
}

// Mutation returns the DChatMutation object of the builder.
func (dc *DChatCreate) Mutation() *DChatMutation {
	return dc.mutation
}

// Save creates the DChat in the database.
func (dc *DChatCreate) Save(ctx context.Context) (*DChat, error) {
	var (
		err  error
		node *DChat
	)
	if len(dc.hooks) == 0 {
		if err = dc.check(); err != nil {
			return nil, err
		}
		node, err = dc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*DChatMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = dc.check(); err != nil {
				return nil, err
			}
			dc.mutation = mutation
			if node, err = dc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(dc.hooks) - 1; i >= 0; i-- {
			if dc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = dc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, dc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (dc *DChatCreate) SaveX(ctx context.Context) *DChat {
	v, err := dc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dc *DChatCreate) Exec(ctx context.Context) error {
	_, err := dc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dc *DChatCreate) ExecX(ctx context.Context) {
	if err := dc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (dc *DChatCreate) check() error {
	if _, ok := dc.mutation.SessionID(); !ok {
		return &ValidationError{Name: "session_id", err: errors.New(`ent: missing required field "session_id"`)}
	}
	if _, ok := dc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "user_id"`)}
	}
	if _, ok := dc.mutation.RoomID(); !ok {
		return &ValidationError{Name: "room_id", err: errors.New(`ent: missing required field "room_id"`)}
	}
	if _, ok := dc.mutation.Text(); !ok {
		return &ValidationError{Name: "text", err: errors.New(`ent: missing required field "text"`)}
	}
	if _, ok := dc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "create_time"`)}
	}
	return nil
}

func (dc *DChatCreate) sqlSave(ctx context.Context) (*DChat, error) {
	_node, _spec := dc.createSpec()
	if err := sqlgraph.CreateNode(ctx, dc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}

func (dc *DChatCreate) createSpec() (*DChat, *sqlgraph.CreateSpec) {
	var (
		_node = &DChat{config: dc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: dchat.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: dchat.FieldID,
			},
		}
	)
	if id, ok := dc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := dc.mutation.SessionID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: dchat.FieldSessionID,
		})
		_node.SessionID = value
	}
	if value, ok := dc.mutation.UserID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: dchat.FieldUserID,
		})
		_node.UserID = value
	}
	if value, ok := dc.mutation.RoomID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: dchat.FieldRoomID,
		})
		_node.RoomID = value
	}
	if value, ok := dc.mutation.Text(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: dchat.FieldText,
		})
		_node.Text = value
	}
	if value, ok := dc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: dchat.FieldCreateTime,
		})
		_node.CreateTime = value
	}
	return _node, _spec
}

// DChatCreateBulk is the builder for creating many DChat entities in bulk.
type DChatCreateBulk struct {
	config
	builders []*DChatCreate
}

// Save creates the DChat entities in the database.
func (dcb *DChatCreateBulk) Save(ctx context.Context) ([]*DChat, error) {
	specs := make([]*sqlgraph.CreateSpec, len(dcb.builders))
	nodes := make([]*DChat, len(dcb.builders))
	mutators := make([]Mutator, len(dcb.builders))
	for i := range dcb.builders {
		func(i int, root context.Context) {
			builder := dcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*DChatMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, dcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, dcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, dcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (dcb *DChatCreateBulk) SaveX(ctx context.Context) []*DChat {
	v, err := dcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dcb *DChatCreateBulk) Exec(ctx context.Context) error {
	_, err := dcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dcb *DChatCreateBulk) ExecX(ctx context.Context) {
	if err := dcb.Exec(ctx); err != nil {
		panic(err)
	}
}
