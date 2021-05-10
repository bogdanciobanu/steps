OUTPUTDIR  :=  out
MANIFESTDIR := $(OUTPUTDIR)/manifests
MANIFEST_PATH ?= gs://stackpulse-steps/
VENDORS_PATH ?= gs://stackpulse-public/
MANIFEST_PARSER ?= gcr.io/stackpulse/step-manifest-parser:prd-21.04.1

.DELETE_ON_ERROR:
# declare-shortcut(source, dest) // declare a phony rule shortcut from src to dst.
define declare-shortcut
	$(eval .PHONY: $(1))                                                                                                       \
	$(eval                                                                                                                     \
		$(1): $(2)                                                                                                             \
	 )
endef

fmt:
	@gofmt -l -w steps

apps:
	@./scripts/baur_apps.py

gomod:
	@./scripts/go_mod.py

pg:
	@./scripts/local_postgres.sh up

local: pg
	@baur run

clean:
	@./scripts/local_postgres.sh down
	rm -rf $(OUTPUTDIR)

all:

ifeq ("$(FORCE_REBUILD)","true")
	@echo "Forcing rebuild"
	$(eval BUILD_ARGS=--force)
endif


ifeq ("$(CIRCLE_BRANCH)","master")
	@echo "Building master branch"
	baur run ${BUILD_ARGS}
else
	@echo "Building side branch"
	baur run --skip-upload ${BUILD_ARGS}
endif

.PHONY: local clean apps gomod all pg publish-manifests-no-deps fmt

packall:
	./scripts/prepare-manifests.py ./steps "$(MANIFESTDIR)"
	docker run --rm -w /root -v $(CURDIR):/root $(MANIFEST_PARSER) validate "/root/$(MANIFESTDIR)/*"


$(call declare-shortcut,indexfile,$(MANIFESTDIR)/indexfile.yml)
$(MANIFESTDIR)/indexfile.yml: packall
	./generate-index-file.py $(MANIFESTDIR) $@

.PHONY: publish-manifests
publish-manifests: packall $(MANIFESTDIR)/indexfile.yml
	gsutil -m cp -r $(MANIFESTDIR)/*.yml $(MANIFEST_PATH)

.PHONY: publish-manifests
publish-manifests-no-deps:
	gsutil -m cp -r $(MANIFESTDIR)/*.yml $(MANIFEST_PATH)

.PHONY: validate-vendors
validate-vendors:
	./validate-vendors.py

.PHONY: publish-vendors
publish-vendors:
	gsutil -m cp -r vendors $(VENDORS_PATH)