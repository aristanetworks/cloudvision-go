// Copyright (c) 2022 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.
package fmp

// PassesPartialEqFilter returns whether the MACAddress matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (x *MACAddress) PassesPartialEqFilter(cmp *MACAddress) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}
	if x == nil {
		return false
	}
	return x.Value == cmp.Value
}

// PassesPartialEqFilter returns whether the RepeatedMACAddress matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (i *RepeatedMACAddress) PassesPartialEqFilter(cmp *RepeatedMACAddress) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}
	if i == nil {
		return false
	}
	if len(i.Values) != len(cmp.Values) {
		return false
	}
	for i, f := range i.Values {
		if !f.PassesPartialEqFilter(cmp.Values[i]) {
			return false
		}
	}
	return true
}

// PassesPartialEqFilter returns whether the IPAddress matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (x *IPAddress) PassesPartialEqFilter(cmp *IPAddress) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}
	if x == nil {
		return false
	}
	return x.Value == cmp.Value
}

// PassesPartialEqFilter returns whether the RepeatedIPAddress matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (i *RepeatedIPAddress) PassesPartialEqFilter(cmp *RepeatedIPAddress) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}
	if i == nil {
		return false
	}
	if len(i.Values) != len(cmp.Values) {
		return false
	}
	for i, f := range i.Values {
		if !f.PassesPartialEqFilter(cmp.Values[i]) {
			return false
		}
	}
	return true
}

// PassesPartialEqFilter returns whether the IPv4Address matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (x *IPv4Address) PassesPartialEqFilter(cmp *IPv4Address) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}
	if x == nil {
		return false
	}
	return x.Value == cmp.Value
}

// PassesPartialEqFilter returns whether the RepeatedIPv4Address matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (i *RepeatedIPv4Address) PassesPartialEqFilter(cmp *RepeatedIPv4Address) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}
	if i == nil {
		return false
	}
	if len(i.Values) != len(cmp.Values) {
		return false
	}
	for i, f := range i.Values {
		if !f.PassesPartialEqFilter(cmp.Values[i]) {
			return false
		}
	}
	return true
}

// PassesPartialEqFilter returns whether the IPv6Address matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (x *IPv6Address) PassesPartialEqFilter(cmp *IPv6Address) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}
	if x == nil {
		return false
	}
	return x.Value == cmp.Value
}

// PassesPartialEqFilter returns whether the RepeatedIPv6Address matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (i *RepeatedIPv6Address) PassesPartialEqFilter(cmp *RepeatedIPv6Address) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}
	if i == nil {
		return false
	}
	if len(i.Values) != len(cmp.Values) {
		return false
	}
	for i, f := range i.Values {
		if !f.PassesPartialEqFilter(cmp.Values[i]) {
			return false
		}
	}
	return true
}
