package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateManagedAccountOptions] = new(CreateManagedAccountRequest)
	_ optionsProvider[DropManagedAccountOptions]   = new(DropManagedAccountRequest)
	_ optionsProvider[ShowManagedAccountOptions]   = new(ShowManagedAccountRequest)
)

type CreateManagedAccountRequest struct {
	name                       AccountObjectIdentifier           // required
	CreateManagedAccountParams CreateManagedAccountParamsRequest // required
}

func (r *CreateManagedAccountRequest) GetName() AccountObjectIdentifier {
	return r.name
}

type CreateManagedAccountParamsRequest struct {
	AdminName     string // required
	AdminPassword string // required
	Comment       *string
}

type DropManagedAccountRequest struct {
	name AccountObjectIdentifier // required
}

type ShowManagedAccountRequest struct {
	Like *Like
}
