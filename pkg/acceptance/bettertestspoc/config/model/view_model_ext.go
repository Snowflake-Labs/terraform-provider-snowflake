package model

func (v *ViewModel) WithDependsOn(values ...string) *ViewModel {
	v.SetDependsOn(values...)
	return v
}
