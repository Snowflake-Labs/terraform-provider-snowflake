package bettertestspoc

func BasicWarehouseModel(
	name string,
	comment string,
) *WarehouseModel {
	return NewWarehouseModel(name).WithComment(comment)
}
