#
# Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

.PHONY: helm/install
## install helm
helm/install: $(BINDIR)/helm

$(BINDIR)/helm:
	mkdir -p $(BINDIR)
	curl "https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3" | HELM_INSTALL_DIR=$(BINDIR) bash

.PHONY: helm-docs/install
## install helm-docs
helm-docs/install: $(BINDIR)/helm-docs

ifeq ($(UNAME),Darwin)
$(BINDIR)/helm-docs:
	mkdir -p $(BINDIR)
	cd $$(mktemp -d) \
	    && curl -LO https://github.com/norwoodj/helm-docs/releases/download/v$(HELM_DOCS_VERSION)/helm-docs_$(HELM_DOCS_VERSION)_Darwin_x86_64.tar.gz \
	    && tar xzvf helm-docs_$(HELM_DOCS_VERSION)_Darwin_x86_64.tar.gz \
	    && mv helm-docs $(BINDIR)/helm-docs
else
$(BINDIR)/helm-docs:
	mkdir -p $(BINDIR)
	cd $$(mktemp -d) \
	    && curl -LO https://github.com/norwoodj/helm-docs/releases/download/v$(HELM_DOCS_VERSION)/helm-docs_$(HELM_DOCS_VERSION)_Linux_x86_64.tar.gz \
	    && tar xzvf helm-docs_$(HELM_DOCS_VERSION)_Linux_x86_64.tar.gz \
	    && mv helm-docs $(BINDIR)/helm-docs
endif

.PHONY: helm/package/vald
## packaging Helm chart for Vald
helm/package/vald:
	helm package charts/vald

.PHONY: helm/package/vald-helm-operator
## packaging Helm chart for vald-helm-operator
helm/package/vald-helm-operator:
	helm package charts/vald-helm-operator

.PHONY: helm/repo/add
## add Helm chart repository
helm/repo/add:
	helm repo add vald https://vald.vdaas.org/charts

.PHONY: helm/docs/vald
helm/docs/vald: charts/vald/README.md

# force to rebuild
.PHONY: charts/vald/README.md
charts/vald/README.md: \
	charts/vald/README.md.gotmpl \
	charts/vald/values.yaml
	helm-docs

.PHONY: helm/docs/vald-helm-operator
helm/docs/vald-helm-operator: charts/vald-helm-operator/README.md

# force to rebuild
.PHONY: charts/vald-helm-operator/README.md
charts/vald-helm-operator/README.md: \
	charts/vald-helm-operator/README.md.gotmpl \
	charts/vald-helm-operator/values.yaml
	helm-docs

.PHONY: helm/schema/vald
## generate json schema for Vald Helm Chart
helm/schema/vald: charts/vald/values.schema.json

charts/vald/values.schema.json: \
	charts/vald/values.yaml \
	hack/helm/schema/gen/main.go
	go run hack/helm/schema/gen/main.go charts/vald/values.yaml > charts/vald/values.schema.json

.PHONY: helm/schema/vald-helm-operator
## generate json schema for Vald Helm Operator Chart
helm/schema/vald-helm-operator: charts/vald-helm-operator/values.schema.json

charts/vald-helm-operator/values.schema.json: \
	charts/vald-helm-operator/values.yaml \
	hack/helm/schema/gen/main.go
	go run hack/helm/schema/gen/main.go charts/vald-helm-operator/values.yaml > charts/vald-helm-operator/values.schema.json
