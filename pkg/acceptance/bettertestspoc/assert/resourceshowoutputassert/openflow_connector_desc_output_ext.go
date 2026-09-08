package resourceshowoutputassert

func (o *OpenflowConnectorDescribeOutputAssert) HasConnectorUrlNotEmpty() *OpenflowConnectorDescribeOutputAssert {
	o.ValuePresent("connector_url")
	return o
}

func (o *OpenflowConnectorDescribeOutputAssert) HasLiveVersionLocationUriNotEmpty() *OpenflowConnectorDescribeOutputAssert {
	o.ValuePresent("live_version_location_uri")
	return o
}
