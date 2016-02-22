
all clean:
	make -C Thrift $@
	make -C backEndProcessor $@
	make -C webService $@
