package resourceshowoutputassert

func (o *OpenflowConnectorShowOutputAssert) HasCreatedOnNotEmpty() *OpenflowConnectorShowOutputAssert {
	o.ValuePresent("created_on")
	return o
}

func (o *OpenflowConnectorShowOutputAssert) HasUpdatedOnNotEmpty() *OpenflowConnectorShowOutputAssert {
	o.ValuePresent("updated_on")
	return o
}

func (o *OpenflowConnectorShowOutputAssert) HasConnectorUrlNotEmpty() *OpenflowConnectorShowOutputAssert {
	o.ValuePresent("connector_url")
	return o
}

func (o *OpenflowConnectorShowOutputAssert) HasLiveVersionLocationUriNotEmpty() *OpenflowConnectorShowOutputAssert {
	o.ValuePresent("live_version_location_uri")
	return o
}
