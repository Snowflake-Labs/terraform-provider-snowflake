package model

func BasicWarehouseModel(
	name string,
	comment string,
) *WarehouseModel {
	return WarehouseWithDefaultMeta(name).WithComment(comment)
}
