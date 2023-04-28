// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package entities

import (
	"time"

	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
)

type OracleSpec struct {
	ExternalDataSourceSpec *ExternalDataSourceSpec
}

func OracleSpecFromProto(spec *zetapb.OracleSpec, txHash TxHash, zetaTime time.Time) (*OracleSpec, error) {
	if spec.ExternalDataSourceSpec != nil {
		ds, err := ExternalDataSourceSpecFromProto(spec.ExternalDataSourceSpec, txHash, zetaTime)
		if err != nil {
			return nil, err
		}

		return &OracleSpec{
			ExternalDataSourceSpec: ds,
		}, nil
	}

	return &OracleSpec{
		ExternalDataSourceSpec: &ExternalDataSourceSpec{},
	}, nil
}

func (os OracleSpec) ToProto() *zetapb.OracleSpec {
	return &zetapb.OracleSpec{
		ExternalDataSourceSpec: os.ExternalDataSourceSpec.ToProto(),
	}
}

func (os OracleSpec) Cursor() *Cursor {
	return NewCursor(DataSourceSpecCursor{os.ExternalDataSourceSpec.Spec.ZetaTime, os.ExternalDataSourceSpec.Spec.ID}.String())
}

func (os OracleSpec) ToProtoEdge(_ ...any) (*v2.OracleSpecEdge, error) {
	return &v2.OracleSpecEdge{
		Node:   os.ToProto(),
		Cursor: os.Cursor().Encode(),
	}, nil
}
