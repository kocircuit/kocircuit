//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package symbol

import (
	"reflect"

	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// Type implementations:
// *OptionalType
//	BasicType, EmptyType,
// *SeriesType, *StructType
// NamedType, *OpaqueType, *MapType
// VarietyType
// BlobType
type Type interface {
	// String returns a string representation of the type
	String() string
	tree.Splayer

	// IsType is only used to enforce the implementation of Type
	IsType()

	// GoType returns the Go equivalent of the type.
	GoType() reflect.Type
}

// Types is a list of Type's
type Types []Type
