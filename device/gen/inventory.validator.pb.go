// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: inventory.proto

package gen

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	regexp "regexp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var _regex_DeviceConfig_DeviceType = regexp.MustCompile(`.+`)

func (this *DeviceConfig) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	if !_regex_DeviceConfig_DeviceType.MatchString(this.DeviceType) {
		return github_com_mwitkow_go_proto_validators.FieldError("DeviceType", fmt.Errorf(`value '%v' must be a string conforming to regex ".+"`, this.DeviceType))
	}
	return nil
}
func (this *DeviceConfigs) Validate() error {
	for _, item := range this.DeviceConfigs {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DeviceConfigs", err)
			}
		}
	}
	return nil
}

var _regex_DeviceInfo_DeviceID = regexp.MustCompile(`.+`)

func (this *DeviceInfo) Validate() error {
	if this.DeviceConfig != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DeviceConfig); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DeviceConfig", err)
		}
	}
	if !_regex_DeviceInfo_DeviceID.MatchString(this.DeviceID) {
		return github_com_mwitkow_go_proto_validators.FieldError("DeviceID", fmt.Errorf(`value '%v' must be a string conforming to regex ".+"`, this.DeviceID))
	}
	return nil
}
func (this *AddRequest) Validate() error {
	if nil == this.DeviceConfig {
		return github_com_mwitkow_go_proto_validators.FieldError("DeviceConfig", fmt.Errorf("message must exist"))
	}
	if this.DeviceConfig != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DeviceConfig); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DeviceConfig", err)
		}
	}
	return nil
}
func (this *AddResponse) Validate() error {
	if nil == this.DeviceInfo {
		return github_com_mwitkow_go_proto_validators.FieldError("DeviceInfo", fmt.Errorf("message must exist"))
	}
	if this.DeviceInfo != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DeviceInfo); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DeviceInfo", err)
		}
	}
	return nil
}

var _regex_DeleteRequest_DeviceID = regexp.MustCompile(`.+`)

func (this *DeleteRequest) Validate() error {
	if !_regex_DeleteRequest_DeviceID.MatchString(this.DeviceID) {
		return github_com_mwitkow_go_proto_validators.FieldError("DeviceID", fmt.Errorf(`value '%v' must be a string conforming to regex ".+"`, this.DeviceID))
	}
	return nil
}
func (this *DeleteResponse) Validate() error {
	return nil
}

var _regex_GetRequest_DeviceID = regexp.MustCompile(`.+`)

func (this *GetRequest) Validate() error {
	if !_regex_GetRequest_DeviceID.MatchString(this.DeviceID) {
		return github_com_mwitkow_go_proto_validators.FieldError("DeviceID", fmt.Errorf(`value '%v' must be a string conforming to regex ".+"`, this.DeviceID))
	}
	return nil
}
func (this *GetResponse) Validate() error {
	if nil == this.DeviceInfo {
		return github_com_mwitkow_go_proto_validators.FieldError("DeviceInfo", fmt.Errorf("message must exist"))
	}
	if this.DeviceInfo != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DeviceInfo); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DeviceInfo", err)
		}
	}
	return nil
}
func (this *ListRequest) Validate() error {
	return nil
}
func (this *ListResponse) Validate() error {
	for _, item := range this.DeviceInfos {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DeviceInfos", err)
			}
		}
	}
	return nil
}
