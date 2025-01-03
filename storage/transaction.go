//storage/transaction.go
package storage

import (
    "context"
)

// Transaction 事务接口
type Transaction interface {
    Get(key string) (*Item, error)
    Set(key string, value []byte, opts *Options) error
    Delete(key string) error
    Commit() error
    Rollback() error
}

// TransactionImpl 事务实现
type TransactionImpl struct {
    store      Store
    operations []Operation
    changes    map[string]*Item
    committed  bool
    rolledback bool
}

// Operation 事务操作
type Operation struct {
    Type      OperationType
    Key       string
    Item      *Item
    Options   *Options
}

type OperationType int

const (
    OpGet OperationType = iota
    OpSet
    OpDelete
)

// Set 设置数据
func (tx *TransactionImpl) Set(key string, value []byte, opts *Options) error {
    if tx.committed || tx.rolledback {
        return errors.New("transaction already finished")
    }
    
    item := &Item{
        Key:      key,
        Value:    value,
        Created:  time.Now(),
        Modified: time.Now(),
    }
    
    if opts != nil && !opts.ExpireAt.IsZero() {
        item.ExpireAt = opts.ExpireAt
    }
    
    tx.operations = append(tx.operations, Operation{
        Type:    OpSet,
        Key:     key,
        Item:    item,
        Options: opts,
    })
    
    tx.changes[key] = item
    return nil
}
