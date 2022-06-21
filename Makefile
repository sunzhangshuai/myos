OPERATING_SYSTEM_GO := lottery malloc mlfq paging-linear-translate relocation scheduler segmentation x86
OPERATING_SYSTEM_C := null syscall syscall-test task1 task2 task8
OPERATING_SYSTEM_GO_CLEAN := $(addprefix clean, $(OPERATING_SYSTEM_GO))
OPERATING_SYSTEM_C_CLEAN := $(addprefix clean, $(OPERATING_SYSTEM_C))

all: operating_system

# 操作系统
operating_system:  $(OPERATING_SYSTEM_GO) $(OPERATING_SYSTEM_C)
$(OPERATING_SYSTEM_GO):
	cd ./operating-system-exercises/go;go build $@.go
$(OPERATING_SYSTEM_C):
	cd ./operating-system-exercises/c;gcc $@.cpp -o $@


clean: operating_system_clean

# 操作系统
operating_system_clean: $(OPERATING_SYSTEM_GO_CLEAN) $(OPERATING_SYSTEM_C_CLEAN)
$(OPERATING_SYSTEM_GO_CLEAN):
	rm -f ./operating-system-exercises/go/$(subst clean,,$@)
$(OPERATING_SYSTEM_C_CLEAN):
	rm -f ./operating-system-exercises/c/$(subst clean,,$@)


.PHONY: all operating_system $(OPERATING_SYSTEM_GO) $(OPERATING_SYSTEM_C) clean operating_system_clean $(OPERATING_SYSTEM_GO_CLEAN) $(OPERATING_SYSTEM_C_CLEAN)