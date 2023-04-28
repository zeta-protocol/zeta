// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.ZETA file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package core

// Noop ...
// this is just here so this package contains
// actual go code (not only test files) and can
// and lsp go does not fail with:
// go build github.com/zeta-protocol/zeta/core/integration: no non-test
// Go files in .../zeta/integration.
func Noop() {}