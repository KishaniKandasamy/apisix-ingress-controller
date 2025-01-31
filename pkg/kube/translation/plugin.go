// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package translation

import (
	configv2alpha1 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2alpha1"
	apisixv1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
)

func (t *translator) translateTrafficSplitPlugin(ar *configv2alpha1.ApisixRoute, defaultBackendWeight int,
	backends []*configv2alpha1.ApisixRouteHTTPBackend) ([]*apisixv1.Upstream, *apisixv1.TrafficSplitConfig, error) {
	var (
		upstreams []*apisixv1.Upstream
		wups      []apisixv1.TrafficSplitConfigRuleWeightedUpstream
	)

	for _, backend := range backends {
		svcClusterIP, svcPort, err := t.getServiceClusterIPAndPort(backend, ar)
		if err != nil {
			return nil, nil, err
		}
		ups, err := t.translateUpstream(ar.Namespace, backend.ServiceName, backend.ResolveGranularity, svcClusterIP, svcPort)
		if err != nil {
			return nil, nil, err
		}
		upstreams = append(upstreams, ups)

		weight := _defaultWeight
		if backend.Weight != 0 {
			weight = backend.Weight
		}
		wups = append(wups, apisixv1.TrafficSplitConfigRuleWeightedUpstream{
			UpstreamID: ups.ID,
			Weight:     weight,
		})
	}

	// Finally append the default upstream in the route.
	wups = append(wups, apisixv1.TrafficSplitConfigRuleWeightedUpstream{
		Weight: defaultBackendWeight,
	})

	tsCfg := &apisixv1.TrafficSplitConfig{
		Rules: []apisixv1.TrafficSplitConfigRule{
			{
				WeightedUpstreams: wups,
			},
		},
	}
	return upstreams, tsCfg, nil
}
