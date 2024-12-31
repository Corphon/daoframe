package force

type Force uint8

const (
    Create  Force = iota // 生之力
    Destroy             // 灭之力
    Transform           // 变之力
    Balance            // 衡之力
)

// ForceInteraction 定义力之间的相互作用规则
var ForceInteraction = map[Force][]Force{
    Create:    {Transform, Balance},
    Destroy:   {Transform, Balance},
    Transform: {Create, Destroy, Balance},
    Balance:   {Create, Destroy, Transform},
}
