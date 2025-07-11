// Code generated by BobGen mysql v0.38.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"database/sql/driver"
	"fmt"
)

// Enum values for BenchmarksResult
const (
	BenchmarksResultPassed BenchmarksResult = "passed"
	BenchmarksResultFailed BenchmarksResult = "failed"
	BenchmarksResultError  BenchmarksResult = "error"
)

func AllBenchmarksResult() []BenchmarksResult {
	return []BenchmarksResult{
		BenchmarksResultPassed,
		BenchmarksResultFailed,
		BenchmarksResultError,
	}
}

type BenchmarksResult string

func (e BenchmarksResult) String() string {
	return string(e)
}

func (e BenchmarksResult) Valid() bool {
	switch e {
	case BenchmarksResultPassed,
		BenchmarksResultFailed,
		BenchmarksResultError:
		return true
	default:
		return false
	}
}

func (e BenchmarksResult) MarshalText() ([]byte, error) {
	return []byte(e), nil
}

func (e *BenchmarksResult) UnmarshalText(text []byte) error {
	return e.Scan(text)
}

func (e BenchmarksResult) MarshalBinary() ([]byte, error) {
	return []byte(e), nil
}

func (e *BenchmarksResult) UnmarshalBinary(data []byte) error {
	return e.Scan(data)
}

func (e BenchmarksResult) Value() (driver.Value, error) {
	return string(e), nil
}

func (e *BenchmarksResult) Scan(value any) error {
	switch x := value.(type) {
	case string:
		*e = BenchmarksResult(x)
	case []byte:
		*e = BenchmarksResult(x)
	case nil:
		return fmt.Errorf("cannot nil into BenchmarksResult")
	default:
		return fmt.Errorf("cannot scan type %T: %v", value, value)
	}

	if !e.Valid() {
		return fmt.Errorf("invalid BenchmarksResult value: %s", *e)
	}

	return nil
}

// Enum values for BenchmarksStatus
const (
	BenchmarksStatusWaiting  BenchmarksStatus = "waiting"
	BenchmarksStatusRunning  BenchmarksStatus = "running"
	BenchmarksStatusFinished BenchmarksStatus = "finished"
)

func AllBenchmarksStatus() []BenchmarksStatus {
	return []BenchmarksStatus{
		BenchmarksStatusWaiting,
		BenchmarksStatusRunning,
		BenchmarksStatusFinished,
	}
}

type BenchmarksStatus string

func (e BenchmarksStatus) String() string {
	return string(e)
}

func (e BenchmarksStatus) Valid() bool {
	switch e {
	case BenchmarksStatusWaiting,
		BenchmarksStatusRunning,
		BenchmarksStatusFinished:
		return true
	default:
		return false
	}
}

func (e BenchmarksStatus) MarshalText() ([]byte, error) {
	return []byte(e), nil
}

func (e *BenchmarksStatus) UnmarshalText(text []byte) error {
	return e.Scan(text)
}

func (e BenchmarksStatus) MarshalBinary() ([]byte, error) {
	return []byte(e), nil
}

func (e *BenchmarksStatus) UnmarshalBinary(data []byte) error {
	return e.Scan(data)
}

func (e BenchmarksStatus) Value() (driver.Value, error) {
	return string(e), nil
}

func (e *BenchmarksStatus) Scan(value any) error {
	switch x := value.(type) {
	case string:
		*e = BenchmarksStatus(x)
	case []byte:
		*e = BenchmarksStatus(x)
	case nil:
		return fmt.Errorf("cannot nil into BenchmarksStatus")
	default:
		return fmt.Errorf("cannot scan type %T: %v", value, value)
	}

	if !e.Valid() {
		return fmt.Errorf("invalid BenchmarksStatus value: %s", *e)
	}

	return nil
}

// Enum values for InstancesStatus
const (
	InstancesStatusRunning  InstancesStatus = "running"
	InstancesStatusBuilding InstancesStatus = "building"
	InstancesStatusStarting InstancesStatus = "starting"
	InstancesStatusStopping InstancesStatus = "stopping"
	InstancesStatusStopped  InstancesStatus = "stopped"
	InstancesStatusDeleting InstancesStatus = "deleting"
	InstancesStatusDeleted  InstancesStatus = "deleted"
)

func AllInstancesStatus() []InstancesStatus {
	return []InstancesStatus{
		InstancesStatusRunning,
		InstancesStatusBuilding,
		InstancesStatusStarting,
		InstancesStatusStopping,
		InstancesStatusStopped,
		InstancesStatusDeleting,
		InstancesStatusDeleted,
	}
}

type InstancesStatus string

func (e InstancesStatus) String() string {
	return string(e)
}

func (e InstancesStatus) Valid() bool {
	switch e {
	case InstancesStatusRunning,
		InstancesStatusBuilding,
		InstancesStatusStarting,
		InstancesStatusStopping,
		InstancesStatusStopped,
		InstancesStatusDeleting,
		InstancesStatusDeleted:
		return true
	default:
		return false
	}
}

func (e InstancesStatus) MarshalText() ([]byte, error) {
	return []byte(e), nil
}

func (e *InstancesStatus) UnmarshalText(text []byte) error {
	return e.Scan(text)
}

func (e InstancesStatus) MarshalBinary() ([]byte, error) {
	return []byte(e), nil
}

func (e *InstancesStatus) UnmarshalBinary(data []byte) error {
	return e.Scan(data)
}

func (e InstancesStatus) Value() (driver.Value, error) {
	return string(e), nil
}

func (e *InstancesStatus) Scan(value any) error {
	switch x := value.(type) {
	case string:
		*e = InstancesStatus(x)
	case []byte:
		*e = InstancesStatus(x)
	case nil:
		return fmt.Errorf("cannot nil into InstancesStatus")
	default:
		return fmt.Errorf("cannot scan type %T: %v", value, value)
	}

	if !e.Valid() {
		return fmt.Errorf("invalid InstancesStatus value: %s", *e)
	}

	return nil
}
