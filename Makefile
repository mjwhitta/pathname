-include gomk/main.mk
-include local/Makefile

clean: clean-default
ifeq ($(unameS),windows)
ifneq ($(wildcard testdata),)
	@powershell -c Remove-Item -Force -Recurse ./testdata
endif
else
	@rm -f -r ./testdata
endif
