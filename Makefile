skeleton: skeleton/bin skeleton/run

skeleton/bin: wsh/*.c
	mkdir -p skeleton/bin
	cd wsh && make
	cp wsh/wsh wsh/wshd skeleton/bin

skeleton/run:
	mkdir -p skeleton/run
