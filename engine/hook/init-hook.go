package hook

type InitHook interface {
  Initialize() bool
  Finalize()
}
