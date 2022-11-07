PROTOC := protoc 

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including proto.mk)
endif

APIDIR := ${ROOT_DIR}/api/
PBDIR := ${ROOT_DIR}/api/singer

PROTOSDIR ?= $(filter-out %.md, $(wildcard ${APIDIR}/protocol/*))
INCLUDE ?= $(foreach proto,${PROTOSDIR},$(addprefix -I, ${proto}))
INPUTS ?= $(foreach proto,${PROTOSDIR}, $(notdir $(wildcard $(proto)/*.proto)))

.PHONY: proto.build
proto.build:	
	@mkdir -p ${PBDIR}
	@$(PROTOC) $(INCLUDE) --go_out=plugins=grpc:${PBDIR} $(INPUTS)
	