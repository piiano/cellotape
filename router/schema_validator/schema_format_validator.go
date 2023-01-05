package schema_validator

// Formats https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-7.3
const (
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid representation according to the "date-time" production.
	dateTimeFormat = "date-time"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid representation according to the "full-date" production.
	dateFormat = "date"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid representation according to the "full-time" production.
	timeFormat = "time"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid representation according to the "duration" production.
	durationFormat = "duration"
	// use with openapi3.TypeString
	// As defined by the "Mailbox" ABNF rule in RFC 5321, section 4.1.2 [RFC5321].
	emailFormat = "email"
	// use with openapi3.TypeString
	// As defined by the extended "Mailbox" ABNF rule in RFC 6531, section 3.3 [RFC6531].
	idnEmailFormat = "idn-email"
	// As defined by RFC 1123, section 2.1 [RFC1123], including host names produced using the Punycode algorithm specified in RFC 5891, section 4.4 [RFC5891].
	hostnameFormat = "hostname"
	// use with openapi3.TypeString
	// As defined by either RFC 1123 as for hostname, or an internationalized hostname as defined by RFC 5890, section 2.3.2.3 [RFC5890].
	idnHostnameFormat = "idn-hostname"
	// use with openapi3.TypeString
	// An IPv4 address according to the "dotted-quad" ABNF syntax as defined in RFC 2673, section 3.2 [RFC2673].
	ipv4Format = "ipv4"
	// use with openapi3.TypeString
	// An IPv6 address as defined in RFC 4291, section 2.2 [RFC4291].
	ipv6Format = "ipv6"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid URI, according to [RFC3986].
	uriFormat = "uri"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid URI Reference (either a URI or a relative-reference), according to [RFC3986].
	uriReferenceFormat = "uri-reference"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid IRI, according to [RFC3987].
	iriFormat = "iri"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid IRI Reference (either an IRI or a relative-reference), according to [RFC3987].
	iriReferenceFormat = "iri-reference"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid string representation of a UUID, according to [RFC4122].
	uuidFormat = "uuid"
	// use with openapi3.TypeString
	// This attribute applies to string instances.
	// A string instance is valid against this attribute if it is a valid URI Template (of any level), according to [RFC6570].
	// Note that URI Templates may be used for IRIs; there is no separate IRI Template specification.
	uriTemplateFormat = "uri-template"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid JSON string representation of a JSON Pointer, according to RFC 6901, section 5 [RFC6901].
	jsonPointerFormat = "json-pointer"
	// use with openapi3.TypeString
	// A string instance is valid against this attribute if it is a valid Relative JSON Pointer [relative-json-pointer].
	relativeJsonPointerFormat = "relative-json-pointer"
	// use with openapi3.TypeString
	// This attribute applies to string instances.
	// A regular expression, which SHOULD be valid according to the ECMA-262 [ecma262] regular expression dialect.
	// Implementations that validate formats MUST accept at least the subset of ECMA-262 defined in the Regular Expressions (Section 4.3) section of this specification, and SHOULD accept all valid ECMA-262 expressions.
	regexFormat = "regex"

	// Additional formats defined by OpenAPI spec https://spec.openapis.org/oas/v3.1.0#data-types
	// use with openapi3.TypeInteger
	int32Format = "int32"
	// use with openapi3.TypeInteger
	int64Format = "int64"
	// use with openapi3.TypeNumber
	floatFormat = "float"
	// use with openapi3.TypeNumber
	doubleFormat = "double"
	// use with openapi3.TypeString
	passwordFormat = "password"
)

// TODO add validation for formats
