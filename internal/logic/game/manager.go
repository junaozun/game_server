package game

type IManager interface {
	RegisterRouter() ExecCommand
}

type MgrManager struct {
	Modules []IManager
}

func NewMgrManager() *MgrManager {
	mgr := &MgrManager{
		Modules: make([]IManager, 0),
	}
	return mgr
}

func (mgr *MgrManager) Register(module IManager) {
	mgr.Modules = append(mgr.Modules, module)
}

func (mgr *MgrManager) UnRegister(module IManager) {
	for k, v := range mgr.Modules {
		if v == module {
			mgr.Modules = append(mgr.Modules[:k], mgr.Modules[k+1:]...)
			break
		}
	}
}
