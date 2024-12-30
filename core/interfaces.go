// core/interfaces.go

// 添加核心接口定义，使模块间依赖更清晰
type (
    // DaoCore 定义核心功能接口
    DaoCore interface {
        Initialize(ctx context.Context) error
        Terminate(ctx context.Context) error
        Health() Health
    }

    // LifeCycleManager 生命周期管理接口
    LifeCycleManager interface {
        DaoCore
        CreateEntity(id string, opts ...EntityOption) error
        DestroyEntity(id string) error
        GetEntityState(id string) (State, error)
    }

    // StateCoordinator 状态协调接口
    StateCoordinator interface {
        DaoCore
        RegisterState(state State) error
        TransitState(from, to State) error
        ValidateTransition(from, to State) bool
    }
)
