package types

import (
	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
)

type DataSourceDefinitionExternal struct {
	SourceType dataSourceType
}

func (e *DataSourceDefinitionExternal) isDataSourceType() {}

func (e *DataSourceDefinitionExternal) oneOfProto() interface{} {
	return e.IntoProto()
}

// /
// IntoProto tries to return the base proto object from DataSourceDefinitionExternal.
func (e *DataSourceDefinitionExternal) IntoProto() *zetapb.DataSourceDefinitionExternal {
	ds := &zetapb.DataSourceDefinitionExternal{}

	if e.SourceType != nil {
		switch dsn := e.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinitionExternal_Oracle:
			ds = &zetapb.DataSourceDefinitionExternal{
				SourceType: dsn,
			}
		}
	}

	return ds
}

func (e *DataSourceDefinitionExternal) String() string {
	if e.SourceType != nil {
		return e.SourceType.String()
	}

	return ""
}

func (e *DataSourceDefinitionExternal) DeepClone() dataSourceType {
	if e.SourceType != nil {
		return e.SourceType.DeepClone()
	}

	return nil
}

// /
// DataSourceDefinitionExternalFromProto tries to build the DataSourceDefinitionExternal object
// from the given proto object..
func DataSourceDefinitionExternalFromProto(protoConfig *zetapb.DataSourceDefinitionExternal) *DataSourceDefinitionExternal {
	ds := &DataSourceDefinitionExternal{
		SourceType: &DataSourceDefinitionExternalOracle{},
	}

	if protoConfig != nil {
		if protoConfig.SourceType != nil {
			switch tp := protoConfig.SourceType.(type) {
			case *zetapb.DataSourceDefinitionExternal_Oracle:
				ds.SourceType = DataSourceDefinitionExternalOracleFromProto(tp)
			}
		}
	}

	return ds
}
