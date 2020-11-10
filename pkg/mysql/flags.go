/*
   --------------------------------------------------------------------------
   Copyright (c) Arroyo Networks - 2020 - All Rights Reserved
   Proprietary and Confidential

   Unauthorized copying of this file, via any medium, is strictly prohibited.
   --------------------------------------------------------------------------
*/

package mysql

import (
	"encoding/json"
)

// Capability is a capability composite flag field.
// Each bit represents an optional feature of the protocol.
// Both the client and server send these.
type Capability uint32

// Capability Flags
const (
	CapabilityLongPassword               Capability = 1 << 0
	CapabilityFoundRows                  Capability = 1 << 1
	CapabilityLongFlag                   Capability = 1 << 2
	CapabilityConnectWithDB              Capability = 1 << 3
	CapabilityNoSchema                   Capability = 1 << 4
	CapabilityCompress                   Capability = 1 << 5
	CapabilityODBC                       Capability = 1 << 6
	CapabilityLocalFiles                 Capability = 1 << 7
	CapabilityIgnoreSpace                Capability = 1 << 8
	CapabilityProtocol41                 Capability = 1 << 9
	CapabilityInteractive                Capability = 1 << 10
	CapabilitySSL                        Capability = 1 << 11
	CapabilityIgnoreSigPipe              Capability = 1 << 12
	CapabilityTransactions               Capability = 1 << 13
	CapabilityReserved                   Capability = 1 << 14
	CapabilityReserved2                  Capability = 1 << 15
	CapabilityMultiStatements            Capability = 1 << 16
	CapabilityMultiResults               Capability = 1 << 17
	CapabilityPSMultiResults             Capability = 1 << 18
	CapabilityPluginAuth                 Capability = 1 << 19
	CapabilityConnectAttrs               Capability = 1 << 20
	CapabilityPluginAuthLenEncClientData Capability = 1 << 21
	CapabilityCanHandleExpiredPasswords  Capability = 1 << 22
	CapabilitySessionTrack               Capability = 1 << 23
	CapabilityDeprecateEOF               Capability = 1 << 24
	CapabilityOptionalResultSetMetadata  Capability = 1 << 25
	CapabilitySSLVerifyServerCert        Capability = 1 << 30
	CapabilityRememberOptions            Capability = 1 << 31
)

func (c Capability) Has(cap Capability) bool {
	return (c & cap) == cap
}

func (c Capability) MarshalJSON() ([]byte, error) {
	names := []string{}

	if c.Has(CapabilityLongPassword) {
		names = append(names, "CLIENT_LONG_PASSWORD")
	}
	if c.Has(CapabilityFoundRows) {
		names = append(names, "CLIENT_FOUND_ROWS")
	}
	if c.Has(CapabilityLongFlag) {
		names = append(names, "CLIENT_LONG_FLAG")
	}
	if c.Has(CapabilityConnectWithDB) {
		names = append(names, "CLIENT_CONNECT_WITH_DB")
	}
	if c.Has(CapabilityNoSchema) {
		names = append(names, "CLIENT_NO_SCHEMA")
	}
	if c.Has(CapabilityCompress) {
		names = append(names, "CLIENT_COMPRESS")
	}
	if c.Has(CapabilityODBC) {
		names = append(names, "CLIENT_ODBC")
	}
	if c.Has(CapabilityLocalFiles) {
		names = append(names, "CLIENT_LOCAL_FILES")
	}
	if c.Has(CapabilityIgnoreSpace) {
		names = append(names, "CLIENT_IGNORE_SPACE")
	}
	if c.Has(CapabilityProtocol41) {
		names = append(names, "CLIENT_PROTOCOL_41")
	}
	if c.Has(CapabilityInteractive) {
		names = append(names, "CLIENT_INTERACTIVE")
	}
	if c.Has(CapabilitySSL) {
		names = append(names, "CLIENT_SSL")
	}
	if c.Has(CapabilityIgnoreSigPipe) {
		names = append(names, "CLIENT_IGNORE_SIGPIPE")
	}
	if c.Has(CapabilityTransactions) {
		names = append(names, "CLIENT_TRANSACTIONS")
	}
	if c.Has(CapabilityReserved) {
		names = append(names, "CLIENT_RESERVED")
	}
	if c.Has(CapabilityReserved2) {
		names = append(names, "CLIENT_RESERVED2")
	}
	if c.Has(CapabilityMultiStatements) {
		names = append(names, "CLIENT_MULTI_STATEMENTS")
	}
	if c.Has(CapabilityMultiResults) {
		names = append(names, "CLIENT_MULTI_RESULTS")
	}
	if c.Has(CapabilityPSMultiResults) {
		names = append(names, "CLIENT_PS_MULTI_RESULTS")
	}
	if c.Has(CapabilityPluginAuth) {
		names = append(names, "CLIENT_PLUGIN_AUTH")
	}
	if c.Has(CapabilityConnectAttrs) {
		names = append(names, "CLIENT_CONNECT_ATTRS")
	}
	if c.Has(CapabilityPluginAuthLenEncClientData) {
		names = append(names, "CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA")
	}
	if c.Has(CapabilityCanHandleExpiredPasswords) {
		names = append(names, "CLIENT_CAN_HANDLED_EXPIRED_PASSWORDS")
	}
	if c.Has(CapabilitySessionTrack) {
		names = append(names, "CLIENT_SESSION_TRACK")
	}
	if c.Has(CapabilityDeprecateEOF) {
		names = append(names, "CLIENT_DEPRECATE_EOF")
	}
	if c.Has(CapabilityOptionalResultSetMetadata) {
		names = append(names, "CLIENT_OPTIONAL_RESULTSET_METADATA")
	}
	if c.Has(CapabilitySSLVerifyServerCert) {
		names = append(names, "CLIENT_SSL_VERIFY_SERVER_CERT")
	}
	if c.Has(CapabilityRememberOptions) {
		names = append(names, "CLIENT_REMEMBER_OPTIONS")
	}

	return json.Marshal(names)
}

// ServerStatus is a server status composite flag field.
type ServerStatus uint16

// ServerStatus Flags
const (
	ServerStatusInTrans           ServerStatus = 1 << 0
	ServerStatusAutoCommit        ServerStatus = 1 << 1
	ServerStatusMoreResultsExist  ServerStatus = 1 << 3
	ServerStatusNoGoodIndexUsed   ServerStatus = 1 << 4
	ServerStatusNoIndexUsed       ServerStatus = 1 << 5
	ServerStatusCursorExists      ServerStatus = 1 << 6
	ServerStatusLastRowSent       ServerStatus = 1 << 7
	ServerStatusDBDropped         ServerStatus = 1 << 8
	ServerStatusNoBacklashEscapes ServerStatus = 1 << 9
	ServerStatusMetadataChanged   ServerStatus = 1 << 10
	ServerStatusQuerySlow         ServerStatus = 1 << 11
	ServerStatusPSOutParams       ServerStatus = 1 << 12
	ServerStatusInTransReadOnly   ServerStatus = 1 << 13
	ServerStatusStateChanged      ServerStatus = 1 << 14
)

// Has determines if the ServerStatus contains the given ServerStatus flag.
func (s ServerStatus) Has(status ServerStatus) bool {
	return (s & status) == status
}

func (s ServerStatus) MarshalJSON() ([]byte, error) {
	names := []string{}

	if s.Has(ServerStatusInTrans) {
		names = append(names, "SERVER_STATUS_IN_TRANS")
	}
	if s.Has(ServerStatusAutoCommit) {
		names = append(names, "SERVER_STATUS_AUTOCOMMIT")
	}
	if s.Has(ServerStatusMoreResultsExist) {
		names = append(names, "SERVER_MORE_RESULT_EXISTS")
	}
	if s.Has(ServerStatusNoGoodIndexUsed) {
		names = append(names, "SERVER_QUERY_NO_GOOD_INDEX_USED")
	}
	if s.Has(ServerStatusNoIndexUsed) {
		names = append(names, "SERVER_QUERY_NO_INDEX_USED")
	}
	if s.Has(ServerStatusCursorExists) {
		names = append(names, "SERVER_STATUS_CURSOR_EXISTS")
	}
	if s.Has(ServerStatusLastRowSent) {
		names = append(names, "SERVER_STATUS_LAST_ROW_SENT")
	}
	if s.Has(ServerStatusDBDropped) {
		names = append(names, "SERVER_STATUS_DB_DROPPED")
	}
	if s.Has(ServerStatusNoBacklashEscapes) {
		names = append(names, "SERVER_STATUS_NO_BACKSLASH_ESCAPES")
	}
	if s.Has(ServerStatusMetadataChanged) {
		names = append(names, "SERVER_STATUS_METADATA_CHANGED")
	}
	if s.Has(ServerStatusQuerySlow) {
		names = append(names, "SERVER_QUERY_WAS_SLOW")
	}
	if s.Has(ServerStatusPSOutParams) {
		names = append(names, "SERVER_PS_OUT_PARAMS")
	}
	if s.Has(ServerStatusInTransReadOnly) {
		names = append(names, "SERVER_STATUS_IN_TRANS_READONLY")
	}
	if s.Has(ServerStatusStateChanged) {
		names = append(names, "SERVER_SESSION_STATE_CHANGED")
	}

	return json.Marshal(names)
}
