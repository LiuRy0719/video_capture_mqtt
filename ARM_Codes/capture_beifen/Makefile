STAGING_DIR := /home/meetyoo/t113/Tina-Linux/out/t113-bingpi_m2/staging_dir/toolchain/bin
export STAGING_DIR
CC = /home/meetyoo/t113/Tina-Linux/prebuilt/gcc/linux-x86/arm/toolchain-sunxi-musl/toolchain/bin/arm-openwrt-linux-gcc
CFLAGS = -g -I.
LDFLAGS = -L../../ARM_mosq/mosquitto-2.0.9/lib libmosquitto.so.1 \
          -L../../ARM_mosq/openssl-1.0.2d -lssl -lcrypto \
          -L/usr/local/lib/jpeg/lib -ljpeg \
          -lpthread

all: capture_mqtt

capture_mqtt: capture_mqtt.o
	$(CC) $^ $(LDFLAGS) -o $@

capture_mqtt.o: capture_mqtt.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f *.o capture

