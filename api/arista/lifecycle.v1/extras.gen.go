// Copyright (c) 2023 Arista Networks, Inc.  All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

//
// Code generated by boomtown. DO NOT EDIT.
//

package lifecycle

// HasKey returns whether the given DeviceLifecycleSummary has a key provided in the model.
func (d *DeviceLifecycleSummary) HasKey() bool {
	return d.GetKey() != nil
}

// HasKey returns whether the given DeviceLifecycleSummaryRequest has a key provided in the request.
func (d *DeviceLifecycleSummaryRequest) HasKey() bool {
	return d.GetKey() != nil
}

// PassesPartialEqFilter returns whether the DateAndModels matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (d *DateAndModels) PassesPartialEqFilter(cmp *DateAndModels) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}

	if d == nil {
		return false
	}

	if cmp.Date != nil {
		if cmp.Date.Seconds != 0 {
			if d.Date.Seconds != cmp.Date.Seconds {
				return false
			}
		}
		if cmp.Date.Nanos != 0 {
			if d.Date.Nanos != cmp.Date.Nanos {
				return false
			}
		}
	}
	if !d.Models.PassesPartialEqFilter(cmp.Models) {
		return false
	}

	return true
}

// PassesPartialEqFilter returns whether the HardwareLifecycleSummary matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (h *HardwareLifecycleSummary) PassesPartialEqFilter(cmp *HardwareLifecycleSummary) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}

	if h == nil {
		return false
	}
	if !h.EndOfLife.PassesPartialEqFilter(cmp.EndOfLife) {
		return false
	}
	if !h.EndOfSale.PassesPartialEqFilter(cmp.EndOfSale) {
		return false
	}
	if !h.EndOfTacSupport.PassesPartialEqFilter(cmp.EndOfTacSupport) {
		return false
	}
	if !h.EndOfHardwareRmaRequests.PassesPartialEqFilter(cmp.EndOfHardwareRmaRequests) {
		return false
	}

	return true
}

// PassesPartialEqFilter returns whether the SoftwareEOL matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (s *SoftwareEOL) PassesPartialEqFilter(cmp *SoftwareEOL) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}

	if s == nil {
		return false
	}

	if cmp.Version != nil {
		if s.Version == nil {
			return false
		}
		if s.Version.Value != cmp.Version.Value {
			return false
		}
	}

	if cmp.EndOfSupport != nil {
		if cmp.EndOfSupport.Seconds != 0 {
			if s.EndOfSupport.Seconds != cmp.EndOfSupport.Seconds {
				return false
			}
		}
		if cmp.EndOfSupport.Nanos != 0 {
			if s.EndOfSupport.Nanos != cmp.EndOfSupport.Nanos {
				return false
			}
		}
	}

	return true
}

// PassesPartialEqFilter returns whether the DeviceLifecycleSummaryKey matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (d *DeviceLifecycleSummaryKey) PassesPartialEqFilter(cmp *DeviceLifecycleSummaryKey) bool {
	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}

	if d == nil {
		return false
	}

	if cmp.DeviceId != nil {
		if d.DeviceId == nil {
			return false
		}
		if d.DeviceId.Value != cmp.DeviceId.Value {
			return false
		}
	}

	return true
}

// PassesPartialEqFilter returns whether the DeviceLifecycleSummary matches the passed in filter.
// On a nil comparison, we consider it a pass. Otherwise, all set (non-nil, initialized)
// fields are expected to match their "sibling" field in the comparison. Any non-matching
// value is considered a mismatch of the filter.
func (d *DeviceLifecycleSummary) PassesPartialEqFilter(cmp *DeviceLifecycleSummary) bool {
	// if the resource is nil, there is nothing to send to the client
	if d == nil {
		return false
	}

	// gave nothing to filter on, consider it passing
	if cmp == nil {
		return true
	}

	if !d.Key.PassesPartialEqFilter(cmp.Key) {
		return false
	}

	if !d.SoftwareEol.PassesPartialEqFilter(cmp.SoftwareEol) {
		return false
	}

	if !d.HardwareLifecycleSummary.PassesPartialEqFilter(cmp.HardwareLifecycleSummary) {
		return false
	}

	return true
}

// MatchesAnyPartialEqFilter returns whether the receiver matches any filters in the given set.
func (d *DeviceLifecycleSummary) MatchesAnyPartialEqFilter(filters []*DeviceLifecycleSummary) bool {
	if len(filters) == 0 {
		return true
	}

	for _, filt := range filters {
		if d.PassesPartialEqFilter(filt) {
			return true
		}
	}

	return false
}
