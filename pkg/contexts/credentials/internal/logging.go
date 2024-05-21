// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	ocmlog "github.com/open-component-model/ocm/pkg/logging"
)

var (
	REALM = ocmlog.DefineSubRealm("Credentials", "credentials")
	log   = ocmlog.DynamicLogger(REALM)
)
