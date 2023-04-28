package types

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/zeta-protocol/zeta/libs/crypto"
	"github.com/zeta-protocol/zeta/libs/proto"
	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
	datapb "github.com/zeta-protocol/zeta/protos/zeta/data/v1"
)

var (
	// ErrDataSourceSpecHasMultipleSameKeyNamesInFilterList is returned when filters with same key names exists inside a single list.
	ErrDataSourceSpecHasMultipleSameKeyNamesInFilterList = errors.New("multiple keys with same name found in filter list")
	// ErrDataSourceSpecHasInvalidTimeCondition is returned when timestamp value is used with 'LessThan'
	// or 'LessThanOrEqual' condition operator value.
	ErrDataSourceSpecHasInvalidTimeCondition = errors.New("data source spec time value is used with 'less than' or 'less than equal' condition")
)

type DataSourceDefinitionInternalx struct {
	Internal *DataSourceDefinitionInternal
}

func (s *DataSourceDefinitionInternalx) isDataSourceType() {}

func (s *DataSourceDefinitionInternalx) oneOfProto() interface{} {
	return s.IntoProto()
}

// IntoProto returns the proto object from DataSourceDefinitionInternalx.
// This method is not called from anywhere.
func (s *DataSourceDefinitionInternalx) IntoProto() *zetapb.DataSourceDefinition_Internal {
	ds := &zetapb.DataSourceDefinition_Internal{
		Internal: &zetapb.DataSourceDefinitionInternal{},
	}

	if s.Internal != nil {
		if s.Internal.SourceType != nil {
			switch dsn := s.Internal.SourceType.oneOfProto().(type) {
			case *zetapb.DataSourceDefinitionInternal_Time:
				ds = &zetapb.DataSourceDefinition_Internal{
					Internal: &zetapb.DataSourceDefinitionInternal{
						SourceType: dsn,
					},
				}
			}
		}
	}

	return ds
}

// DeepClone returns a clone of the DataSourceDefinitionInternalx object.
func (s *DataSourceDefinitionInternalx) DeepClone() dataSourceType {
	cpy := &DataSourceDefinitionInternalx{
		Internal: &DataSourceDefinitionInternal{
			SourceType: s.Internal.SourceType.DeepClone(),
		},
	}
	return cpy
}

// String returns the DataSourceDefinitionInternalx content as a string.
func (s *DataSourceDefinitionInternalx) String() string {
	if s.Internal != nil {
		// Does not return the type of the internal data source, becase the base object
		// definitions are located in core/zeta/protos/ and do not access the local interface
		// and accessing it will lead to cycle import.
		return fmt.Sprintf("internal(%s)", s.Internal.String())
	}

	return ""
}

type DataSourceDefinitionExternalx struct {
	External *DataSourceDefinitionExternal
}

func (s *DataSourceDefinitionExternalx) isDataSourceType() {}

func (s *DataSourceDefinitionExternalx) oneOfProto() interface{} {
	return s.IntoProto()
}

// IntoProto returns the proto object from DataSourceDefinitionInternalx
// This method is not called from anywhere.
func (s *DataSourceDefinitionExternalx) IntoProto() *zetapb.DataSourceDefinition_External {
	ds := &zetapb.DataSourceDefinition_External{
		External: &zetapb.DataSourceDefinitionExternal{},
	}

	if s.External != nil {
		if s.External.SourceType != nil {
			switch dsn := s.External.SourceType.oneOfProto().(type) {
			case *zetapb.DataSourceDefinitionExternal_Oracle:
				ds = &zetapb.DataSourceDefinition_External{
					External: &zetapb.DataSourceDefinitionExternal{
						SourceType: dsn,
					},
				}
			}
		}
	}

	return ds
}

func (s *DataSourceDefinitionExternalx) DeepClone() dataSourceType {
	cpy := &DataSourceDefinitionExternalx{
		External: &DataSourceDefinitionExternal{
			SourceType: s.External.SourceType.DeepClone(),
		},
	}
	return cpy
}

// String returns the DataSourceDefinitionExternalx content as a string.
func (s *DataSourceDefinitionExternalx) String() string {
	if s.External != nil {
		// Does not return the type of the external data source, becase the base object
		// definitions are located in core/zeta/protos/ and do not access the local intrface
		// and accessing it will lead to cycle import.
		return fmt.Sprintf("external(%s)", s.External.String())
	}

	return ""
}

type dataSourceType interface {
	isDataSourceType()
	oneOfProto() interface{}

	String() string
	DeepClone() dataSourceType
	// ToDataSourceSpec() *DataSourceSpec
}

type DataSourceDefinition struct {
	SourceType dataSourceType
}

// /
// IntoProto returns the proto object from DataSourceDefinition
// that is - zetapb.DataSourceDefinition that may have external or internal SourceType.
// Returns the whole proto object.
func (s DataSourceDefinition) IntoProto() *zetapb.DataSourceDefinition {
	ds := &zetapb.DataSourceDefinition{}

	if s.SourceType != nil {
		switch dsn := s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_External:
			if dsn.External != nil {
				if dsn.External.SourceType != nil {
					switch dsn.External.SourceType.(type) {
					case *zetapb.DataSourceDefinitionExternal_Oracle:
						ds = &zetapb.DataSourceDefinition{
							SourceType: &zetapb.DataSourceDefinition_External{
								External: &zetapb.DataSourceDefinitionExternal{
									// This will return the external data source oracle object that satisfies the interface
									SourceType: dsn.External.GetSourceType(),
								},
							},
						}
					}
				}
			}

		case *zetapb.DataSourceDefinition_Internal:
			if dsn.Internal != nil {
				if dsn.Internal.SourceType != nil {
					switch dsn.Internal.SourceType.(type) {
					case *zetapb.DataSourceDefinitionInternal_Time:
						ds = &zetapb.DataSourceDefinition{
							SourceType: &zetapb.DataSourceDefinition_Internal{
								Internal: &zetapb.DataSourceDefinitionInternal{
									// This will return the internal data source time object that satisfies the interface
									SourceType: dsn.Internal.GetSourceType(),
								},
							},
						}
						// More types of internal sources that will come in the future - will go here.
					}
				}
			}
		}
	}

	return ds
}

// /
// String returns the data source definition content as a string.
func (s DataSourceDefinition) String() string {
	if s.SourceType != nil {
		switch dsn := s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_External:
			s := ""
			if dsn.External != nil {
				s = dsn.External.String()
			}
			return fmt.Sprintf("external(%s)", s)

		case *zetapb.DataSourceDefinition_Internal:
			s := ""
			if dsn.Internal != nil {
				s = dsn.Internal.String()
			}
			return fmt.Sprintf("internal(%s)", s)
		}
	}

	return ""
}

// DeepClone returns a clone of the DataSourceDefinition object.
func (s DataSourceDefinition) DeepClone() DataSourceDefinition {
	cpy := &DataSourceDefinition{}

	if s.SourceType != nil {
		switch dsn := s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_External:
			cpy = &DataSourceDefinition{
				SourceType: &DataSourceDefinitionExternalx{
					External: &DataSourceDefinitionExternal{
						SourceType: &DataSourceSpecConfiguration{
							Signers: SignersFromProto(dsn.External.GetOracle().Signers),
							Filters: DataSourceSpecFiltersFromProto(dsn.External.GetOracle().Filters),
						},
					},
				},
			}

		case *zetapb.DataSourceDefinition_Internal:
			cpy = &DataSourceDefinition{
				SourceType: s.SourceType.DeepClone(),
			}
		}
	}

	return *cpy
}

// /
// DataSourceDefinitionFromProto tries to build the DataSourceDfiniition object
// from the given proto object.
func DataSourceDefinitionFromProto(protoConfig *zetapb.DataSourceDefinition) *DataSourceDefinition {
	if protoConfig != nil {
		if protoConfig.SourceType != nil {
			switch tp := protoConfig.SourceType.(type) {
			case *zetapb.DataSourceDefinition_External:
				return &DataSourceDefinition{
					SourceType: &DataSourceDefinitionExternalx{
						// Checking if the tp.External is nil is made in the `DataSourceDefinitionExternalFromProto` step
						External: DataSourceDefinitionExternalFromProto(tp.External),
					},
				}

			case *zetapb.DataSourceDefinition_Internal:
				return &DataSourceDefinition{
					SourceType: &DataSourceDefinitionInternalx{
						// Checking if the tp.Internal is nil is made in the `DataSourceDefinitionInternalFromProto` step
						Internal: DataSourceDefinitionInternalFromProto(tp.Internal),
					},
				}
			}
		}
	}

	return &DataSourceDefinition{}
}

// /
// GetSigners tries to get the signers from the DataSourceDefinition if they exist.
func (s DataSourceDefinition) GetSigners() []*Signer {
	signers := []*Signer{}

	if s.SourceType != nil {
		switch dsn := s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_External:
			if dsn.External != nil {
				if dsn.External.SourceType != nil {
					o := dsn.External.GetOracle()
					if o != nil {
						signers = SignersFromProto(o.Signers)
					}
				}
			}
		case **zetapb.DataSourceDefinition_Internal:
			// Do not try to get signers from internal type of data source
		}
	}

	return signers
}

// /
// GetFilters tries to get the filters from the DataSourceDefinition if they exist.
func (s DataSourceDefinition) GetFilters() []*DataSourceSpecFilter {
	filters := []*DataSourceSpecFilter{}

	if s.SourceType != nil {
		switch dsn := s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_External:
			if dsn.External != nil {
				if dsn.External.SourceType != nil {
					switch dsn.External.SourceType.(type) {
					case *zetapb.DataSourceDefinitionExternal_Oracle:
						o := dsn.External.GetOracle()
						if o != nil {
							filters = DataSourceSpecFiltersFromProto(o.Filters)
						}
					}
				}
			}

		// For the case the definition source type is zetapb.DataSourceDefinition_Internal
		case *zetapb.DataSourceDefinition_Internal:
			if dsn.Internal != nil {
				if dsn.Internal.SourceType != nil {
					switch idsn := dsn.Internal.SourceType.(type) {
					// For the case the zetapb.DataSourceDefinitionInternal is not nill
					// and its embedded object is of type zetapb.DataSourceDefinitionInternal_Time
					case *zetapb.DataSourceDefinitionInternal_Time:
						// Retrieve the filters for this specific type of internal source
						// Checking if the idsn.Time object is nil is done in `DataSourceSpecConfigurationTimeFromProto`
						dst := DataSourceSpecConfigurationTimeFromProto(idsn.Time)

						// For the case the internal data source is time based
						// (as of OT https://github.com/zetaprotocol/specs/blob/master/protocol/0048-DSRI-data_source_internal.md#13-zeta-time-changed)
						// We add the filter key values manually to match a time based data source
						// Ensure only a single filter has been created, that holds the first condition
						if len(dst.Conditions) > 0 {
							filters = append(
								filters,
								&DataSourceSpecFilter{
									Key: &DataSourceSpecPropertyKey{
										Name: "zetaprotocol.builtin.timestamp",
										Type: datapb.PropertyKey_TYPE_TIMESTAMP,
									},
									Conditions: []*DataSourceSpecCondition{
										dst.Conditions[0],
									},
								},
							)
						}
					}
				}
			}
		}
	}

	return filters
}

// GetDataSourceSpecConfiguration returns the base object - DataSourceSpecConfiguration
// from the DataSourceDefinition.
func (s DataSourceDefinition) GetDataSourceSpecConfiguration() *zetapb.DataSourceSpecConfiguration {
	if s.SourceType != nil {
		switch tp := s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_External:
			if tp.External != nil {
				if tp.External.SourceType != nil {
					o := tp.External.GetOracle()
					if o != nil {
						return o
					}
				}
			}
		case *zetapb.DataSourceDefinitionInternal:
			//
		}
	}

	return &zetapb.DataSourceSpecConfiguration{}
}

// NewDataSourceDefinition creates a new EMPTY DataSourceDefinition object.
func NewDataSourceDefinition(tp int) *DataSourceDefinition {
	ds := &DataSourceDefinition{}
	switch tp {
	case zetapb.DataSourceDefinitionTypeInt:
		ds.SourceType = &DataSourceDefinitionInternalx{
			Internal: &DataSourceDefinitionInternal{
				// Create internal type definition with time for now.
				SourceType: &DataSourceDefinitionInternalTime{
					Time: &DataSourceSpecConfigurationTime{
						Conditions: []*DataSourceSpecCondition{},
					},
				},
			},
		}

	case zetapb.DataSourceDefinitionTypeExt:
		ds.SourceType = &DataSourceDefinitionExternalx{
			External: &DataSourceDefinitionExternal{
				// Create external definition for oracles for now.
				// Extened when needed.
				SourceType: &DataSourceDefinitionExternalOracle{
					Oracle: &DataSourceSpecConfiguration{
						Signers: []*Signer{},
						Filters: []*DataSourceSpecFilter{},
					},
				},
			},
		}
	}

	return ds
}

// /
// UpdateFilters updates the DataSourceDefinition Filters.
func (s *DataSourceDefinition) UpdateFilters(filters []*DataSourceSpecFilter) error {
	fTypeCheck := map[*DataSourceSpecFilter]struct{}{}
	fNameCheck := map[string]struct{}{}
	for _, f := range filters {
		if _, ok := fTypeCheck[f]; ok {
			return ErrDataSourceSpecHasMultipleSameKeyNamesInFilterList
		}
		if f.Key != nil {
			if _, ok := fNameCheck[f.Key.Name]; ok {
				return ErrDataSourceSpecHasMultipleSameKeyNamesInFilterList
			}
			fNameCheck[f.Key.Name] = struct{}{}
		}
		fTypeCheck[f] = struct{}{}
	}

	if s.SourceType != nil {
		switch dsn := s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_External:
			if dsn.External != nil {
				if dsn.External.SourceType != nil {
					switch dsn.External.SourceType.(type) {
					case *zetapb.DataSourceDefinitionExternal_Oracle:
						o := dsn.External.GetOracle()
						signers := []*datapb.Signer{}
						if o != nil {
							signers = o.Signers
						}

						ds := &zetapb.DataSourceDefinition{
							SourceType: &zetapb.DataSourceDefinition_External{
								External: &zetapb.DataSourceDefinitionExternal{
									SourceType: &zetapb.DataSourceDefinitionExternal_Oracle{
										Oracle: &zetapb.DataSourceSpecConfiguration{
											// We do not care if return empty lists for signers and filters here
											Filters: DataSourceSpecFilters(filters).IntoProto(),
											Signers: signers,
										},
									},
								},
							},
						}

						dsd := DataSourceDefinitionFromProto(ds)
						if dsd.SourceType != nil {
							*s = *dsd
						}
					}
				}
			}

		case *zetapb.DataSourceDefinition_Internal:
			if dsn.Internal != nil {
				if dsn.Internal.SourceType != nil {
					switch dsn.Internal.SourceType.(type) {
					case *zetapb.DataSourceDefinitionInternal_Time:
						// The data source definition is an internal time based source
						// For this case we take only the first item from the list of filters
						// https://github.com/zetaprotocol/specs/blob/master/protocol/0048-DSRI-data_source_internal.md#13-zeta-time-changed
						c := []*datapb.Condition{}
						if len(filters) > 0 {
							if len(filters[0].Conditions) > 0 {
								c = append(c, filters[0].IntoProto().Conditions[0])
							}
						}
						ds := &zetapb.DataSourceDefinition{
							SourceType: &zetapb.DataSourceDefinition_Internal{
								Internal: &zetapb.DataSourceDefinitionInternal{
									SourceType: &zetapb.DataSourceDefinitionInternal_Time{
										Time: &zetapb.DataSourceSpecConfigurationTime{
											Conditions: c,
										},
									},
								},
							},
						}

						dsd := DataSourceDefinitionFromProto(ds)
						if dsd.SourceType != nil {
							*s = *dsd
						}
					}
				}
			}
		}
	}

	return nil
}

func (s *DataSourceDefinition) SetFilterDecimals(d uint64) *DataSourceDefinition {
	if s.SourceType != nil {
		switch dsn := s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_External:
			filters := dsn.External.GetOracle().Filters
			for i := range filters {
				filters[i].Key.NumberDecimalPlaces = &d
			}

			ds := &zetapb.DataSourceDefinition{
				SourceType: &zetapb.DataSourceDefinition_External{
					External: &zetapb.DataSourceDefinitionExternal{
						SourceType: &zetapb.DataSourceDefinitionExternal_Oracle{
							Oracle: &zetapb.DataSourceSpecConfiguration{
								Filters: filters,
								Signers: dsn.External.GetOracle().Signers,
							},
						},
					},
				},
			}

			dsd := DataSourceDefinitionFromProto(ds)
			if dsd.SourceType != nil {
				*s = *dsd
			}
		}
	}

	return s
}

func (s DataSourceDefinition) ToDataSourceSpec() *DataSourceSpec {
	bytes, _ := proto.Marshal(s.IntoProto())
	specID := hex.EncodeToString(crypto.Hash(bytes))
	return &DataSourceSpec{
		ID:   specID,
		Data: &s,
	}
}

func (s *DataSourceDefinition) ToExternalDataSourceSpec() *ExternalDataSourceSpec {
	return &ExternalDataSourceSpec{
		Spec: s.ToDataSourceSpec(),
	}
}

// SetOracleConfig sets a given oracle config in the receiver.
// If the receiver is not external oracle type data source - it is not changed.
// This method does not care about object previous contents.
func (s *DataSourceDefinition) SetOracleConfig(oc *DataSourceSpecConfiguration) *DataSourceDefinition {
	if s.SourceType != nil {
		switch def := s.SourceType.oneOfProto().(type) {
		// For the case the definition source type is zetapb.DataSourceDefinition_External
		case *zetapb.DataSourceDefinition_External:
			if def.External != nil {
				if def.External.SourceType != nil {
					switch def.External.SourceType.(type) {
					// For the case the zetapb.DataSourceDefinitionExternal is not nill
					// and its embedded object is of type zetapb.DataSourceDefinitionExternal_Oracle
					case *zetapb.DataSourceDefinitionExternal_Oracle:
						// Set the new config only in this case
						ds := &zetapb.DataSourceDefinition{
							SourceType: &zetapb.DataSourceDefinition_External{
								External: &zetapb.DataSourceDefinitionExternal{
									SourceType: &zetapb.DataSourceDefinitionExternal_Oracle{
										Oracle: oc.IntoProto(),
									},
								},
							},
						}

						dsd := DataSourceDefinitionFromProto(ds)
						if dsd.SourceType != nil {
							*s = *dsd
						}
					}
				}
			}
		}
	}

	return s
}

// SetTimeTriggerConditionConfig sets a given conditions config in the receiver.
// If the receiver is not a time triggered data source - it does not set anything to it.
// This method does not care about object previous contents.
func (s *DataSourceDefinition) SetTimeTriggerConditionConfig(c []*DataSourceSpecCondition) *DataSourceDefinition {
	if s.SourceType != nil {
		switch def := s.SourceType.oneOfProto().(type) {
		// For the case the definition source type is zetapb.DataSourceDefinition_Internal
		case *zetapb.DataSourceDefinition_Internal:
			if def.Internal != nil {
				if def.Internal.SourceType != nil {
					switch def.Internal.SourceType.(type) {
					// For the case the zetapb.DataSourceDefinitionInternal is not nill
					// and its embedded object is of type zetapb.DataSourceDefinitionInternal_Time
					case *zetapb.DataSourceDefinitionInternal_Time:
						// Set the new first condition only in this case
						cond := []*datapb.Condition{}
						if len(c) > 0 {
							cond = append(cond, c[0].IntoProto())
						}

						ds := &zetapb.DataSourceDefinition{
							SourceType: &zetapb.DataSourceDefinition_Internal{
								Internal: &zetapb.DataSourceDefinitionInternal{
									SourceType: &zetapb.DataSourceDefinitionInternal_Time{
										Time: &zetapb.DataSourceSpecConfigurationTime{
											// We do not care if we return an empty list of conditions in this place
											Conditions: cond,
										},
									},
								},
							},
						}

						dsd := DataSourceDefinitionFromProto(ds)
						if dsd.SourceType != nil {
							*s = *dsd
						}
					}
				}
			}
		}
	}

	return s
}

func (s *DataSourceDefinition) IsExternal() (bool, error) {
	if s.SourceType != nil {
		switch s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_External:
			return true, nil
		}

		return false, nil
	}

	return false, errors.New("unknown type of data source provided")
}

func (s *DataSourceDefinition) GetDataSourceSpecConfigurationTime() *zetapb.DataSourceSpecConfigurationTime {
	if s.SourceType != nil {
		switch tp := s.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinition_Internal:
			if tp.Internal != nil {
				if tp.Internal.SourceType != nil {
					return tp.Internal.GetTime()
				}
			}
		}
	}

	return &zetapb.DataSourceSpecConfigurationTime{}
}
