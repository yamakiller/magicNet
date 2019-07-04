package hook

type InitializeHook interface {
  Initialize() bool
  Finalize()
}
