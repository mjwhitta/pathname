-include gomk/main.mk
-include local/Makefile

clean: clean-default
ifeq ($(unameS),windows)
ifneq ($(wildcard testdata),)
	@remove-item -force -recurse ./testdata
endif
else
	@rm -f -r ./testdata
endif
