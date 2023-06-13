package config

// RulesOption represents the option for some rules.
type RulesOption struct {
	FileNamesLowerSnakeCase         FileNamesLowerSnakeCaseOption         `yaml:"file_names_lower_snake_case"`
	QuoteConsistentOption           QuoteConsistentOption                 `yaml:"quote_consistent"`
	ImportsSorted                   ImportsSortedOption                   `yaml:"imports_sorted"`
	MaxLineLength                   MaxLineLengthOption                   `yaml:"max_line_length"`
	Indent                          IndentOption                          `yaml:"indent"`
	EnumFieldNamesZeroValueEndWith  EnumFieldNamesZeroValueEndWithOption  `yaml:"enum_field_names_zero_value_end_with"`
	ServiceNamesEndWith             ServiceNamesEndWithOption             `yaml:"service_names_end_with"`
	FieldNamesExcludePrepositions   FieldNamesExcludePrepositionsOption   `yaml:"field_names_exclude_prepositions"`
	MessageNamesExcludePrepositions MessageNamesExcludePrepositionsOption `yaml:"message_names_exclude_prepositions"`
	RPCNamesCaseOption              RPCNamesCaseOption                    `yaml:"rpc_names_case"`
	MessagesHaveComment             MessagesHaveCommentOption             `yaml:"messages_have_comment"`
	ServicesHaveComment             ServicesHaveCommentOption             `yaml:"services_have_comment"`
	RPCsHaveComment                 RPCsHaveCommentOption                 `yaml:"rpcs_have_comment"`
	FieldsHaveComment               FieldsHaveCommentOption               `yaml:"fields_have_comment"`
	EnumsHaveComment                EnumsHaveCommentOption                `yaml:"enums_have_comment"`
	EnumFieldsHaveComment           EnumFieldsHaveCommentOption           `yaml:"enum_fields_have_comment"`
	SyntaxConsistent                SyntaxConsistentOption                `yaml:"syntax_consistent"`
	RepeatedFieldNamesPluralized    RepeatedFieldNamesPluralizedOption    `yaml:"repeated_field_names_pluralized"`
	EnumFieldNamesPrefix            CustomizableSeverityOption            `yaml:"enum_field_names_prefix"`
	EnumFieldNamesUpperSnakeCase    CustomizableSeverityOption            `yaml:"enum_field_names_upper_snake_case"`
	EnumNamesUpperCamelCase         CustomizableSeverityOption            `yaml:"enum_names_upper_camel_case"`
	FieldNamesLowerSnakeCase        CustomizableSeverityOption            `yaml:"field_names_lower_snake_case"`
	FileHasComment                  CustomizableSeverityOption            `yaml:"file_has_comment"`
	MessageNamesUpperCamelCase      CustomizableSeverityOption            `yaml:"message_names_upper_camel_case"`
	Order                           CustomizableSeverityOption            `yaml:"order"`
	PackageNameLowerCase            CustomizableSeverityOption            `yaml:"package_name_lower_case"`
	Proto3FieldsAvoidRequired       CustomizableSeverityOption            `yaml:"proto3_fields_avoid_required"`
	Proto3GroupsAvoid               CustomizableSeverityOption            `yaml:"proto3_groups_avoid"`
	RPCNamesUpperCamelCase          CustomizableSeverityOption            `yaml:"rpc_names_upper_camel_case"`
	ServiceNamesUpperCamelCase      CustomizableSeverityOption            `yaml:"service_names_upper_caml_case"`
}
