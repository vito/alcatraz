skeleton: skeleton/bin skeleton/ssh skeleton/ssh_auth

skeleton/bin: iomux/*.c
	mkdir -p skeleton/bin
	cd iomux && make
	mv iomux/iomux-link iomux/iomux-spawn skeleton/bin

skeleton/ssh:
	mkdir -p skeleton/ssh

skeleton/ssh_auth:
	mkdir -p skeleton/ssh_auth
