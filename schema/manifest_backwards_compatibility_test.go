// Copyright 2016 The Linux Foundation
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

package schema_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/opencontainers/image-spec/schema"
)

var compatMap = map[string]string{
	"application/vnd.docker.distribution.manifest.list.v2+json": "application/vnd.oci.image.manifest.list.v1+json",
	"application/vnd.docker.distribution.manifest.v2+json":      "application/vnd.oci.image.manifest.v1+json",
	"application/vnd.docker.image.rootfs.diff.tar.gzip":         "application/vnd.oci.image.rootfs.tar.gzip",
	"application/vnd.docker.container.image.v1+json":            "application/vnd.oci.image.serialization.config.v1+json",
}

// convertFormats converts Docker v2.2 image format JSON documents to OCI
// format by simply replacing instances of the strings found in the compatMap
// found in the input string.
func convertFormats(input string) string {
	out := input
	for k, v := range compatMap {
		out = strings.Replace(out, v, k, -1)
	}
	return out
}

func TestBackwardsCompatibilityManifestList(t *testing.T) {
	for i, tt := range []struct {
		manifest string
		digest   string
		fail     bool
	}{
		{
			digest: "sha256:e588eb8123f2031a41f2e60bc27f30a4388e181e07410aff392f7dc96b585969",
			manifest: `{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
   "manifests": [
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v1+json",
         "size": 2094,
         "digest": "sha256:7820f9a86d4ad15a2c4f0c0e5479298df2aa7c2f6871288e2ef8546f3e7b6783",
         "platform": {
            "architecture": "ppc64le",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v1+json",
         "size": 1922,
         "digest": "sha256:ae1b0e06e8ade3a11267564a26e750585ba2259c0ecab59ab165ad1af41d1bdd",
         "platform": {
            "architecture": "amd64",
            "os": "linux",
            "features": [
               "sse"
            ]
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v1+json",
         "size": 2084,
         "digest": "sha256:e4c0df75810b953d6717b8f8f28298d73870e8aa2a0d5e77b8391f16fdfbbbe2",
         "platform": {
            "architecture": "s390x",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v1+json",
         "size": 2084,
         "digest": "sha256:07ebe243465ef4a667b78154ae6c3ea46fdb1582936aac3ac899ea311a701b40",
         "platform": {
            "architecture": "arm",
            "os": "linux",
            "variant": "armv7"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v1+json",
         "size": 2090,
         "digest": "sha256:fb2fc0707b86dafa9959fe3d29e66af8787aee4d9a23581714be65db4265ad8a",
         "platform": {
            "architecture": "arm64",
            "os": "linux",
            "variant": "armv8"
         }
      }
   ]
}`,
			fail: false,
		},
	} {
		sum := sha256.Sum256([]byte(tt.manifest))
		got := fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
		if tt.digest != got {
			t.Errorf("test %d: expected digest %s but got %s", i, tt.digest, got)
		}

		manifest := convertFormats(tt.manifest)
		r := strings.NewReader(manifest)
		err := schema.MediaTypeManifestList.Validate(r)

		if got := err != nil; tt.fail != got {
			t.Errorf("test %d: expected validation failure %t but got %t, err %v", i, tt.fail, got, err)
		}
	}
}

func TestBackwardsCompatibilityManifest(t *testing.T) {
	for i, tt := range []struct {
		manifest string
		digest   string
		fail     bool
	}{
		// manifest pulled from docker hub using hash value
		//
		// curl -L -H "Authorization: Bearer ..." -H \
		// "Accept: application/vnd.docker.distribution.manifest.v2+json" \
		// https://registry-1.docker.io/v2/library/docker/manifests/sha256:888206c77cd2811ec47e752ba291e5b7734e3ef137dfd222daadaca39a9f17bc
		{
			digest: "sha256:888206c77cd2811ec47e752ba291e5b7734e3ef137dfd222daadaca39a9f17bc",
			manifest: `{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
   "config": {
      "mediaType": "application/octet-stream",
      "size": 3210,
      "digest": "sha256:5359a4f250650c20227055957e353e8f8a74152f35fe36f00b6b1f9fc19c8861"
   },
   "layers": [
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 2310272,
         "digest": "sha256:fae91920dcd4542f97c9350b3157139a5d901362c2abec284de5ebd1b45b4957"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 913022,
         "digest": "sha256:f384f6ab36adad485192f09379c0b58dc612a3cde82c551e082a7c29a87c95da"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 9861668,
         "digest": "sha256:ed0d2dd5e1a0e5e650a330a864c8a122e9aa91fa6ba9ac6f0bd1882e59df55e7"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 465,
         "digest": "sha256:ec4d00b58417c45f7ddcfde7bcad8c9d62a7d6d5d17cdc1f7d79bcb2e22c1491"
      }
   ]
}`,
			fail: false,
		},
	} {
		sum := sha256.Sum256([]byte(tt.manifest))
		got := fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
		if tt.digest != got {
			t.Errorf("test %d: expected digest %s but got %s", i, tt.digest, got)
		}

		manifest := convertFormats(tt.manifest)
		r := strings.NewReader(manifest)
		err := schema.MediaTypeManifest.Validate(r)

		if got := err != nil; tt.fail != got {
			t.Errorf("test %d: expected validation failure %t but got %t, err %v", i, tt.fail, got, err)
		}
	}
}
