package types

import (
	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
)

type DataSourceDefinitionInternal struct {
	SourceType dataSourceType
}

func (i *DataSourceDefinitionInternal) isDataSourceType() {}

func (i *DataSourceDefinitionInternal) oneOfProto() interface{} {
	return i.IntoProto()
}

func (i *DataSourceDefinitionInternal) IntoProto() interface{} {
	ds := &zetapb.DataSourceDefinitionInternal{}

	if i.SourceType != nil {
		switch dsn := i.SourceType.oneOfProto().(type) {
		case *zetapb.DataSourceDefinitionInternal_Time:
			ds = &zetapb.DataSourceDefinitionInternal{
				SourceType: dsn,
			}
		}
	}

	return ds
}

func (i *DataSourceDefinitionInternal) String() string {
	if i.SourceType != nil {
		return i.SourceType.String()
	}

	return ""
}

func (i *DataSourceDefinitionInternal) DeepClone() dataSourceType {
	if i.SourceType != nil {
		return i.SourceType.DeepClone()
	}

	return nil
}

// //
// DataSourceDefinitionInternalFromProto tries to build the DataSourceDefinitionInternal object
// from the given proto configuration.
func DataSourceDefinitionInternalFromProto(protoConfig *zetapb.DataSourceDefinitionInternal) *DataSourceDefinitionInternal {
	ds := &DataSourceDefinitionInternal{
		SourceType: &DataSourceDefinitionInternalTime{},
	}

	if protoConfig != nil {
		if protoConfig.SourceType != nil {
			switch tp := protoConfig.SourceType.(type) {
			case *zetapb.DataSourceDefinitionInternal_Time:
				ds.SourceType = DataSourceDefinitionInternalTimeFromProto(tp)
			}
		}
	}

	return ds
}
