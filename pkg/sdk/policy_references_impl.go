package sdk

import "context"

type policyReference struct {
	client *Client
}

func (v *policyReference) GetForEntity(ctx context.Context, request *GetForEntityPolicyReferenceRequest) ([]PolicyReference, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[policyReferenceDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	resultList := convertRows[policyReferenceDBRow, PolicyReference](dbRows)

	return resultList, nil
}

func (request *GetForEntityPolicyReferenceRequest) toOpts() *getForEntityPolicyReferenceOptions {
	return &getForEntityPolicyReferenceOptions{
		tableFunction: &tableFunction{
			table: Bool(true),
			policyReferenceFunction: &policyReferenceFunction{
				functionFullyQualifiedName: Bool(true),
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewObjectIdentifierFromFullyQualifiedName(request.RefEntityName)},
					refEntityDomain: String(request.RefEntityDomain),
				},
			},
		},
	}
}
