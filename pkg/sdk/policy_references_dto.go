package sdk

var _ optionsProvider[getForEntityPolicyReferenceOptions] = new(GetForEntityPolicyReferenceRequest)

type GetForEntityPolicyReferenceRequest struct {
	RefEntityName   string
	RefEntityDomain string
}
