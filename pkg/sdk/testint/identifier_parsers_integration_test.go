package testint

//func TestInt_IdentifierParsers(t *testing.T) {
//	testCases := []struct {
//		IdentifierType string
//		Input          string
//		Expected       sdk.ObjectIdentifier
//	}{
//		// Note: Won't work because of the internal SDK identifier validation and the fact
//		//{IdentifierType: "AccountObjectIdentifier", Input: `""`, Expected: sdk.NewAccountObjectIdentifier(``)},
//		//{IdentifierType: "AccountObjectIdentifier", Input: `""""`, Expected: sdk.NewAccountObjectIdentifier(`"`)},
//		//{IdentifierType: "AccountObjectIdentifier", Input: sdk.NewAccountObjectIdentifier(`"""`)},
//
//		{IdentifierType: "AccountObjectIdentifier", Input: "", Expected: sdk.NewAccountObjectIdentifier(`abc`)},
//		{IdentifierType: "AccountObjectIdentifier", Input: "", Expected: sdk.NewAccountObjectIdentifier(`abc`)},
//		{IdentifierType: "AccountObjectIdentifier", Input: "", Expected: sdk.NewAccountObjectIdentifier(`ab.c`)},
//
//		// TODO(): Won't work because of Show like 'a\"\"bc' # additional escaping
//		// 	Additional rules for SDK validation should be added (like quotes inside identifiers)
//		//{IdentifierType: "AccountObjectIdentifier", Input: sdk.NewAccountObjectIdentifier(`a""bc`)},
//
//		//{IdentifierType: "DatabaseObjectIdentifier", Input: `"".""`, Expected: sdk.NewDatabaseObjectIdentifier(``, ``)},
//		//{IdentifierType: "DatabaseObjectIdentifier", Input: `"""".""""`, Expected: sdk.NewDatabaseObjectIdentifier(`"`, `"`)},
//		//{IdentifierType: "DatabaseObjectIdentifier", Input: `abc.cde`, Expected: sdk.NewDatabaseObjectIdentifier(`abc`, `cde`)},
//		//{IdentifierType: "DatabaseObjectIdentifier", Input: `"abc"."cde"`, Expected: sdk.NewDatabaseObjectIdentifier(`abc`, `cde`)},
//		//{IdentifierType: "DatabaseObjectIdentifier", Input: `"ab.c"."cd.e"`, Expected: sdk.NewDatabaseObjectIdentifier(`ab.c`, `cd.e`)},
//		//{IdentifierType: "DatabaseObjectIdentifier", Input: `"a""bc"."cd""e"`, Expected: sdk.NewDatabaseObjectIdentifier(`a"bc`, `cd"e`)},
//		//
//		//{IdentifierType: "SchemaObjectIdentifier", Input: ``, Error: "incompatible identifier: "},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: "a\nb.cde.efg", Error: "unable to read identifier: a\nb.cde.efg, err = record on line 2: wrong number of fields"},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: `a"b.cde.efg`, Error: "unable to read identifier: a\"b.cde.efg, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: `abc.cde.efg.ghi`, Error: `unexpected number of parts 4 in identifier abc.cde.efg.ghi, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: `""."".""`, Expected: sdk.NewSchemaObjectIdentifier(``, ``, ``)},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: `""""."""".""""`, Expected: sdk.NewSchemaObjectIdentifier(`"`, `"`, `"`)},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: `abc.cde.efg`, Expected: sdk.NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: `"abc"."cde"."efg"`, Expected: sdk.NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: `"ab.c"."cd.e"."ef.g"`, Expected: sdk.NewSchemaObjectIdentifier(`ab.c`, `cd.e`, `ef.g`)},
//		//{IdentifierType: "SchemaObjectIdentifier", Input: `"a""bc"."cd""e"."ef""g"`, Expected: sdk.NewSchemaObjectIdentifier(`a"bc`, `cd"e`, `ef"g`)},
//		//
//		//{IdentifierType: "TableColumnIdentifier", Input: ``, Error: "incompatible identifier: "},
//		//{IdentifierType: "TableColumnIdentifier", Input: "a\nb.cde.efg.ghi", Error: "unable to read identifier: a\nb.cde.efg.ghi, err = record on line 2: wrong number of fields"},
//		//{IdentifierType: "TableColumnIdentifier", Input: `a"b.cde.efg.ghi`, Error: "unable to read identifier: a\"b.cde.efg.ghi, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
//		//{IdentifierType: "TableColumnIdentifier", Input: `abc.cde.efg.ghi.ijk`, Error: `unexpected number of parts 5 in identifier abc.cde.efg.ghi.ijk, expected 4 in a form of "<database_name>.<schema_name>.<table_name>.<table_column_name>"`},
//		//{IdentifierType: "TableColumnIdentifier", Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 4 in a form of "<database_name>.<schema_name>.<table_name>.<table_column_name>"`},
//		//{IdentifierType: "TableColumnIdentifier", Input: `"".""."".""`, Expected: sdk.NewTableColumnIdentifier(``, ``, ``, ``)},
//		//{IdentifierType: "TableColumnIdentifier", Input: `"""".""""."""".""""`, Expected: sdk.NewTableColumnIdentifier(`"`, `"`, `"`, `"`)},
//		//{IdentifierType: "TableColumnIdentifier", Input: `abc.cde.efg.ghi`, Expected: sdk.NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
//		//{IdentifierType: "TableColumnIdentifier", Input: `"abc"."cde"."efg"."ghi"`, Expected: sdk.NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
//		//{IdentifierType: "TableColumnIdentifier", Input: `"ab.c"."cd.e"."ef.g"."gh.i"`, Expected: sdk.NewTableColumnIdentifier(`ab.c`, `cd.e`, `ef.g`, `gh.i`)},
//		//{IdentifierType: "TableColumnIdentifier", Input: `"a""bc"."cd""e"."ef""g"."gh""i"`, Expected: sdk.NewTableColumnIdentifier(`a"bc`, `cd"e`, `ef"g`, `gh"i`)},
//		//
//		//{IdentifierType: "AccountIdentifier", Input: ``, Error: "incompatible identifier: "},
//		//{IdentifierType: "AccountIdentifier", Input: "a\nb.cde", Error: "unable to read identifier: a\nb.cde, err = record on line 2: wrong number of fields"},
//		//{IdentifierType: "AccountIdentifier", Input: `a"b.cde`, Error: "unable to read identifier: a\"b.cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
//		//{IdentifierType: "AccountIdentifier", Input: `abc.cde.efg`, Error: `unexpected number of parts 3 in identifier abc.cde.efg, expected 2 in a form of "<organization_name>.<account_name>"`},
//		//{IdentifierType: "AccountIdentifier", Input: `abc`, Error: `unexpected number of parts 1 in identifier abc, expected 2 in a form of "<organization_name>.<account_name>"`},
//		//{IdentifierType: "AccountIdentifier", Input: `"".""`, Expected: sdk.NewAccountIdentifier(``, ``)},
//		//{IdentifierType: "AccountIdentifier", Input: `"""".""""`, Expected: sdk.NewAccountIdentifier(`"`, `"`)},
//		//{IdentifierType: "AccountIdentifier", Input: `abc.cde`, Expected: sdk.NewAccountIdentifier(`abc`, `cde`)},
//		//{IdentifierType: "AccountIdentifier", Input: `"abc"."cde"`, Expected: sdk.NewAccountIdentifier(`abc`, `cde`)},
//		//{IdentifierType: "AccountIdentifier", Input: `"ab.c"."cd.e"`, Expected: sdk.NewAccountIdentifier(`ab.c`, `cd.e`)},
//		//{IdentifierType: "AccountIdentifier", Input: `"a""bc"."cd""e"`, Expected: sdk.NewAccountIdentifier(`a"bc`, `cd"e`)},
//		//
//		//{IdentifierType: "ExternalObjectIdentifier", Input: ``, Error: "incompatible identifier: "},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: "a\nb.cde.efg", Error: "unable to read identifier: a\nb.cde.efg, err = record on line 2: wrong number of fields"},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: `a"b.cde.efg`, Error: "unable to read identifier: a\"b.cde.efg, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: `abc.cde.efg.ghi`, Error: `unexpected number of parts 4 in identifier abc.cde.efg.ghi, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: `""."".""`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(``, ``), sdk.NewAccountObjectIdentifier(``))},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: `""""."""".""""`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`"`, `"`), sdk.NewAccountObjectIdentifier(`"`))},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: `abc.cde.efg`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`abc`, `cde`), sdk.NewAccountObjectIdentifier(`efg`))},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: `"abc"."cde"."efg"`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`abc`, `cde`), sdk.NewAccountObjectIdentifier(`efg`))},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: `"ab.c"."cd.e"."ef.g"`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`ab.c`, `cd.e`), sdk.NewAccountObjectIdentifier(`ef.g`))},
//		//{IdentifierType: "ExternalObjectIdentifier", Input: `"a""bc"."cd""e"."ef""g"`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`a"bc`, `cd"e`), sdk.NewAccountObjectIdentifier(`ef"g`))},
//	}
//
//	role, roleCleanup := testClientHelper().Role.CreateRole(t)
//	t.Cleanup(roleCleanup)
//
//	for _, testCase := range testCases {
//		t.Run(fmt.Sprintf(`Parsing %s with input: "%s"`, testCase.IdentifierType, testCase.Input), func(t *testing.T) {
//			switch testCase.IdentifierType {
//			case "AccountObjectIdentifier":
//				//id := testCase.Input.(sdk.AccountObjectIdentifier)
//				//
//				//database, databaseCleanup := testClientHelper().Database.CreateDatabaseWithIdentifier(t, id)
//				//t.Cleanup(databaseCleanup)
//				//
//				//err := testClient(t).Grants.GrantPrivilegesToAccountRole(
//				//	context.Background(),
//				//	&sdk.AccountRoleGrantPrivileges{AccountObjectPrivileges: []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}},
//				//	&sdk.AccountRoleGrantOn{
//				//		AccountObject: &sdk.GrantOnAccountObject{
//				//			Database: sdk.Pointer(database.ID()),
//				//		},
//				//	},
//				//	role.ID(),
//				//	&sdk.GrantPrivilegesToAccountRoleOptions{},
//				//)
//				//require.NoError(t, err)
//				//
//				//grants, err := testClient(t).Grants.Show(context.Background(), &sdk.ShowGrantOptions{
//				//	To: &sdk.ShowGrantsTo{
//				//		Role: role.ID(),
//				//	},
//				//})
//				//require.NoError(t, err)
//				//
//				//assert.Equal(t, 1, len(grants))
//				//
//				//databaseId, err := sdk.ParseAccountObjectIdentifier(grants[0].Name.FullyQualifiedName())
//				//require.NoError(t, err)
//				//assert.Equal(t, database.ID(), databaseId)
//
//				//assert.Equal(t, id, database.ID())
//			default:
//				t.SkipNow()
//
//				//case "DatabaseObjectIdentifier":
//				//	id, err = sdk.ParseDatabaseObjectIdentifier(testCase.Input)
//				//case "SchemaObjectIdentifier":
//				//	id, err = sdk.ParseSchemaObjectIdentifier(testCase.Input)
//				//case "TableColumnIdentifier":
//				//	id, err = sdk.ParseTableColumnIdentifier(testCase.Input)
//				//case "AccountIdentifier":
//				//	id, err = sdk.ParseAccountIdentifier(testCase.Input)
//				//case "ExternalObjectIdentifier":
//				//	id, err = sdk.ParseExternalObjectIdentifier(testCase.Input)
//				//default:
//				//	t.Fatalf("unknown identifier type: %s", testCase.IdentifierType)
//			}
//		})
//	}
//}
