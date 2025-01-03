//storage/sync.go
package storage

// StoreSyncer 存储同步器
type StoreSyncer struct {
    source     Store
    target     Store
    differ     *StoreDiffer
    batcher    *SyncBatcher
    metrics    *SyncMetrics
}

// SyncOptions 同步选项
type SyncOptions struct {
    BatchSize      int
    RetryAttempts  int
    RetryDelay     time.Duration
    SkipErrors     bool
    DryRun         bool
}

// Sync 执行同步
func (s *StoreSyncer) Sync(ctx context.Context, opts *SyncOptions) error {
    // 获取差异
    diffs, err := s.differ.Diff(ctx)
    if err != nil {
        return err
    }
    
    // 创建批次
    batches := s.batcher.CreateBatches(diffs, opts.BatchSize)
    
    // 同步每个批次
    for _, batch := range batches {
        if err := s.syncBatch(ctx, batch, opts); err != nil {
            if !opts.SkipErrors {
                return err
            }
            s.metrics.Errors.Inc()
        }
    }
    
    return nil
}
