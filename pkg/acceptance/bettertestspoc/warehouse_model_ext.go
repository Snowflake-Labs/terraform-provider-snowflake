package bettertestspoc

func BasicWarehouseModel(
	name string,
	comment string,
) *WarehouseModel {
	return NewDefaultWarehouseModel(name).WithComment(comment)
}
