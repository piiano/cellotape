package schema_validator

type SchemaFormat string

// Formats https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-7.3
const (
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid representation according to the "date-time" production.
	dateTimeFormat SchemaFormat = "date-time"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid representation according to the "full-date" production.
	dateFormat SchemaFormat = "date"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid representation according to the "full-time" production.
	timeFormat SchemaFormat = "time"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid representation according to the "duration" production.
	durationFormat SchemaFormat = "duration"
	// use with stringSchemaType
	// As defined by the "Mailbox" ABNF rule in RFC 5321, section 4.1.2 [RFC5321].
	emailFormat SchemaFormat = "email"
	// use with stringSchemaType
	// As defined by the extended "Mailbox" ABNF rule in RFC 6531, section 3.3 [RFC6531].
	idnEmailFormat SchemaFormat = "idn-email"
	// As defined by RFC 1123, section 2.1 [RFC1123], including host names produced using the Punycode algorithm specified in RFC 5891, section 4.4 [RFC5891].
	hostnameFormat SchemaFormat = "hostname"
	// use with stringSchemaType
	// As defined by either RFC 1123 as for hostname, or an internationalized hostname as defined by RFC 5890, section 2.3.2.3 [RFC5890].
	idnHostnameFormat SchemaFormat = "idn-hostname"
	// use with stringSchemaType
	// An IPv4 address according to the "dotted-quad" ABNF syntax as defined in RFC 2673, section 3.2 [RFC2673].
	ipv4Format SchemaFormat = "ipv4"
	// use with stringSchemaType
	// An IPv6 address as defined in RFC 4291, section 2.2 [RFC4291].
	ipv6Format SchemaFormat = "ipv6"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid URI, according to [RFC3986].
	uriFormat SchemaFormat = "uri"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid URI Reference (either a URI or a relative-reference), according to [RFC3986].
	uriReferenceFormat SchemaFormat = "uri-reference"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid IRI, according to [RFC3987].
	iriFormat SchemaFormat = "iri"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid IRI Reference (either an IRI or a relative-reference), according to [RFC3987].
	iriReferenceFormat SchemaFormat = "iri-reference"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid string representation of a UUID, according to [RFC4122].
	uuidFormat SchemaFormat = "uuid"
	// use with stringSchemaType
	// This attribute applies to string instances.
	// A string instance is valid against this attribute if it is a valid URI Template (of any level), according to [RFC6570].
	// Note that URI Templates may be used for IRIs; there is no separate IRI Template specification.
	uriTemplateFormat SchemaFormat = "uri-template"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid JSON string representation of a JSON Pointer, according to RFC 6901, section 5 [RFC6901].
	jsonPointerFormat SchemaFormat = "json-pointer"
	// use with stringSchemaType
	// A string instance is valid against this attribute if it is a valid Relative JSON Pointer [relative-json-pointer].
	relativeJsonPointerFormat SchemaFormat = "relative-json-pointer"
	// use with stringSchemaType
	// This attribute applies to string instances.
	// A regular expression, which SHOULD be valid according to the ECMA-262 [ecma262] regular expression dialect.
	// Implementations that validate formats MUST accept at least the subset of ECMA-262 defined in the Regular Expressions (Section 4.3) section of this specification, and SHOULD accept all valid ECMA-262 expressions.
	regexFormat SchemaFormat = "regex"

	// Additional formats defined by OpenAPI spec https://spec.openapis.org/oas/v3.1.0#data-types
	// use with integerSchemaType
	int32Format SchemaFormat = "int32"
	// use with integerSchemaType
	int64Format SchemaFormat = "int64"
	// use with numberSchemaType
	floatFormat SchemaFormat = "float"
	// use with numberSchemaType
	doubleFormat SchemaFormat = "double"
	// use with stringSchemaType
	passwordFormat SchemaFormat = "password"
)

// TODO add validation for formats
