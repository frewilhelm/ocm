// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package ociimage

import (
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/open-component-model/ocm/cmds/ocm/commands/ocmcmds/common/inputs"
	"github.com/open-component-model/ocm/cmds/ocm/commands/ocmcmds/common/inputs/cpi"
	"github.com/open-component-model/ocm/pkg/common/accessio"
	"github.com/open-component-model/ocm/pkg/contexts/oci"
	"github.com/open-component-model/ocm/pkg/contexts/oci/grammar"
	"github.com/open-component-model/ocm/pkg/contexts/oci/repositories/artifactset"
	"github.com/open-component-model/ocm/pkg/contexts/oci/repositories/docker"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/ociartifact"
)

type Spec struct {
	// PathSpec holds the repository path and tag of the image in the docker daemon
	cpi.PathSpec
	// Repository is the repository hint for the index artifact
	Repository string `json:"repository,omitempty"`
}

var _ inputs.InputSpec = (*Spec)(nil)

func New(pathtag string) *Spec {
	return &Spec{
		PathSpec: cpi.NewPathSpec(TYPE, pathtag),
	}
}

func (s *Spec) Validate(fldPath *field.Path, ctx inputs.Context, inputFilePath string) field.ErrorList {
	allErrs := s.PathSpec.Validate(fldPath, ctx, inputFilePath)
	allErrs = ValidateRepository(fldPath.Child("repository"), allErrs, s.Repository)
	if s.Path != "" {
		pathField := fldPath.Child("path")
		_, _, err := docker.ParseGenericRef(s.Path)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(pathField, s.Path, err.Error()))
		}
	}
	return allErrs
}

func (s *Spec) GetBlob(ctx inputs.Context, info inputs.InputResourceInfo) (accessio.TemporaryBlobAccess, string, error) {
	ctx.Printf("image %s\n", s.Path)
	ref, err := oci.ParseRef(s.Path)
	if err != nil {
		return nil, "", err
	}

	spec, err := ctx.OCIContext().MapUniformRepositorySpec(&ref.UniformRepositorySpec)
	if err != nil {
		return nil, "", err
	}

	repo, err := ctx.OCIContext().RepositoryForSpec(spec)
	if err != nil {
		return nil, "", err
	}
	ns, err := repo.LookupNamespace(ref.Repository)
	if err != nil {
		return nil, "", err
	}

	version := ref.Version()
	if version == "" || version == "latest" {
		version = info.ComponentVersion.GetVersion()
	}
	blob, err := artifactset.SynthesizeArtifactBlob(ns, version)
	if err != nil {
		return nil, "", err
	}
	return blob, ociartifact.Hint(info.ComponentVersion, info.ElementName, s.Repository, version), nil
}

func ValidateRepository(fldPath *field.Path, allErrs field.ErrorList, repo string) field.ErrorList {
	if repo == "" {
		return allErrs
	}
	if strings.Contains(repo, grammar.DigestSeparator) || strings.Contains(repo, grammar.TagSeparator) {
		return append(allErrs, field.Invalid(fldPath, repo, "unexpected digest or tag"))
	}

	if !grammar.AnchoredRepositoryRegexp.MatchString(strings.TrimPrefix(repo, "/")) {
		return append(allErrs, field.Invalid(fldPath, repo, "no repository name"))
	}
	return allErrs
}
