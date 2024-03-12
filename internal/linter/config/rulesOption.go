package config

// RulesOption represents the option for some rules.
type RulesOption struct {
	FileNamesLowerSnakeCase         FileNamesLowerSnakeCaseOption         `yaml:"file_names_lower_snake_case" json:"file_names_lower_snake_case" toml:"file_names_lower_snake_case"`
	QuoteConsistentOption           QuoteConsistentOption                 `yaml:"quote_consistent" json:"quote_consistent" toml:"quote_consistent"`
	ImportsSorted                   ImportsSortedOption                   `yaml:"imports_sorted" json:"imports_sorted" toml:"imports_sorted"`
	MaxLineLength                   MaxLineLengthOption                   `yaml:"max_line_length" json:"max_line_length" toml:"max_line_length"`
	Indent                          IndentOption                          `yaml:"indent" json:"indent" toml:"indent"`
	EnumFieldNamesZeroValueEndWith  EnumFieldNamesZeroValueEndWithOption  `yaml:"enum_field_names_zero_value_end_with" json:"enum_field_names_zero_value_end_with" toml:"enum_field_names_zero_value_end_with"`
	ServiceNamesEndWith             ServiceNamesEndWithOption             `yaml:"service_names_end_with" json:"service_names_end_with" toml:"service_names_end_with"`
	FieldNamesExcludePrepositions   FieldNamesExcludePrepositionsOption   `yaml:"field_names_exclude_prepositions" json:"field_names_exclude_prepositions" toml:"field_names_exclude_prepositions"`
	MessageNamesExcludePrepositions MessageNamesExcludePrepositionsOption `yaml:"message_names_exclude_prepositions" json:"message_names_exclude_prepositions" toml:"message_names_exclude_prepositions"`
	RPCNamesCaseOption              RPCNamesCaseOption                    `yaml:"rpc_names_case" json:"rpc_names_case" toml:"rpc_names_case"`
	MessagesHaveComment             MessagesHaveCommentOption             `yaml:"messages_have_comment" json:"messages_have_comment" toml:"messages_have_comment"`
	ServicesHaveComment             ServicesHaveCommentOption             `yaml:"services_have_comment" json:"services_have_comment" toml:"services_have_comment"`
	RPCsHaveComment                 RPCsHaveCommentOption                 `yaml:"rpcs_have_comment" json:"rpcs_have_comment" toml:"rpcs_have_comment"`
	FieldsHaveComment               FieldsHaveCommentOption               `yaml:"fields_have_comment" json:"fields_have_comment" toml:"fields_have_comment"`
	EnumsHaveComment                EnumsHaveCommentOption                `yaml:"enums_have_comment" json:"enums_have_comment" toml:"enums_have_comment"`
	EnumFieldsHaveComment           EnumFieldsHaveCommentOption           `yaml:"enum_fields_have_comment" json:"enum_fields_have_comment" toml:"enum_fields_have_comment"`
	SyntaxConsistent                SyntaxConsistentOption                `yaml:"syntax_consistent" json:"syntax_consistent" toml:"syntax_consistent"`
	RepeatedFieldNamesPluralized    RepeatedFieldNamesPluralizedOption    `yaml:"repeated_field_names_pluralized" json:"repeated_field_names_pluralized" toml:"repeated_field_names_pluralized"`
	EnumFieldNamesPrefix            CustomizableSeverityOption            `yaml:"enum_field_names_prefix" json:"enum_field_names_prefix" toml:"enum_field_names_prefix"`
	EnumFieldNamesUpperSnakeCase    CustomizableSeverityOption            `yaml:"enum_field_names_upper_snake_case" json:"enum_field_names_upper_snake_case" toml:"enum_field_names_upper_snake_case"`
	EnumNamesUpperCamelCase         CustomizableSeverityOption            `yaml:"enum_names_upper_camel_case" json:"enum_names_upper_camel_case" toml:"enum_names_upper_camel_case"`
	FieldNamesLowerSnakeCase        CustomizableSeverityOption            `yaml:"field_names_lower_snake_case" json:"field_names_lower_snake_case" toml:"field_names_lower_snake_case"`
	FileHasComment                  CustomizableSeverityOption            `yaml:"file_has_comment" json:"file_has_comment" toml:"file_has_comment"`
	MessageNamesUpperCamelCase      CustomizableSeverityOption            `yaml:"message_names_upper_camel_case" json:"message_names_upper_camel_case" toml:"message_names_upper_camel_case"`
	Order                           CustomizableSeverityOption            `yaml:"order" json:"order" toml:"order"`
	PackageNameLowerCase            CustomizableSeverityOption            `yaml:"package_name_lower_case" json:"package_name_lower_case" toml:"package_name_lower_case"`
	Proto3FieldsAvoidRequired       CustomizableSeverityOption            `yaml:"proto3_fields_avoid_required" json:"proto3_fields_avoid_required" toml:"proto3_fields_avoid_required"`
	Proto3GroupsAvoid               CustomizableSeverityOption            `yaml:"proto3_groups_avoid" json:"proto3_groups_avoid" toml:"proto3_groups_avoid"`
	RPCNamesUpperCamelCase          CustomizableSeverityOption            `yaml:"rpc_names_upper_camel_case" json:"rpc_names_upper_camel_case" toml:"rpc_names_upper_camel_case"`
	ServiceNamesUpperCamelCase      CustomizableSeverityOption            `yaml:"service_names_upper_caml_case" json:"service_names_upper_caml_case" toml:"service_names_upper_caml_case"`
}
