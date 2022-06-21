PROJECT_GO := lottery malloc mlfq paging-linear-translate relocation scheduler segmentation x86
PROJECT_C := null syscall syscall-test task1 task2 task8
PROJECT_GO_CLEAN := $(addprefix clean, $(PROJECT_GO))
PROJECT_C_CLEAN := $(addprefix clean, $(PROJECT_C))

all: $(PROJECT_GO) $(PROJECT_C)

# 操作系统导论
$(PROJECT_GO):
	cd ./go;go build $@.go
$(PROJECT_C):
	cd ./c;gcc $@.cpp -o $@


clean: $(PROJECTGO_CLEAN) $(PROJECTC_CLEAN)

# 操作系统导论
$(PROJECT_GO_CLEAN):
	rm -f ./go/$(subst clean,,$@)
$(PROJECT_C_CLEAN):
	rm -f ./c/$(subst clean,,$@)


.PHONY: all operating_system $(PROJECTGO) $(PROJECTC) clean PROJECTclean $(PROJECTGO_CLEAN) $(PROJECTC_CLEAN)