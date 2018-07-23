phony_targets = all

.PHONY: $(phony_targets)

$(phony_targets) :
	@$(MAKE) subdir="$(CURDIR)" -C "$(TOP_SRCDIR)" $@
