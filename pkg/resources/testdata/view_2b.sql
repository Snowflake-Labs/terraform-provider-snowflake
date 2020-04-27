 SELECT
             value:configRuleName::VARCHAR AS CONFIG_RULE_NAME,
             VALUE:complianceType::VARCHAR AS COMPLIANCE_TYPE,
             *
         FROM (
             SELECT
                 parse_json(CONFIGURATION:configRuleList) AS SRC,
                 *
             FROM "SNOWALERT"."DATA"."AWS_CONFIG_DEFAULT_EVENTS_CONNECTION"
             WHERE RESOURCE_TYPE = 'AWS::Config::ResourceCompliance'
         ), lateral flatten(input => SRC)