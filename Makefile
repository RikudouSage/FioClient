.DEFAULT_GOAL := build-lib-current
SQLCIPHER_DIR := $(CURDIR)/external/sqlcipher
GO_BUILD_TAGS := libsqlite3 no_postgres no_mysql no_ydb no_clickhouse no_libsql no_mssql no_vertica
GO_SQLITE3_MODFILE := $(CURDIR)/build/go-sqlite3-sqlcipher.mod

SQLCIPHER_CFLAGS := -O2 -fPIC \
	-DSQLITE_HAS_CODEC \
	-DSQLCIPHER_CRYPTO_OPENSSL \
	-DSQLITE_EXTRA_INIT=sqlcipher_extra_init \
	-DSQLITE_EXTRA_SHUTDOWN=sqlcipher_extra_shutdown

SQLCIPHER_CURRENT_DIR := $(CURDIR)/build/sqlcipher/current
SQLCIPHER_386_DIR     := $(CURDIR)/build/sqlcipher/386
SQLCIPHER_ARM7_DIR    := $(CURDIR)/build/sqlcipher/arm7
SQLCIPHER_ARM64_DIR   := $(CURDIR)/build/sqlcipher/arm64

SQLCIPHER_CURRENT_INSTALL := $(SQLCIPHER_CURRENT_DIR)/install
SQLCIPHER_386_INSTALL     := $(SQLCIPHER_386_DIR)/install
SQLCIPHER_ARM7_INSTALL    := $(SQLCIPHER_ARM7_DIR)/install
SQLCIPHER_ARM64_INSTALL   := $(SQLCIPHER_ARM64_DIR)/install

SQLCIPHER_CURRENT_LIB := $(SQLCIPHER_CURRENT_INSTALL)/lib/libsqlite3.a
SQLCIPHER_386_LIB     := $(SQLCIPHER_386_INSTALL)/lib/libsqlite3.a
SQLCIPHER_ARM7_LIB    := $(SQLCIPHER_ARM7_INSTALL)/lib/libsqlite3.a
SQLCIPHER_ARM64_LIB   := $(SQLCIPHER_ARM64_INSTALL)/lib/libsqlite3.a


# ------------------------------------------------------------------------------
# SQLCipher
# ------------------------------------------------------------------------------

build-sqlcipher-current: $(SQLCIPHER_CURRENT_LIB)

$(SQLCIPHER_CURRENT_LIB): $(SQLCIPHER_DIR)/configure
	mkdir -p $(SQLCIPHER_CURRENT_DIR)
	cd $(SQLCIPHER_CURRENT_DIR) && \
		CC="gcc" \
		CPPFLAGS="-I$$OPENSSL_CURRENT_INCLUDE" \
		CFLAGS="$(SQLCIPHER_CFLAGS)" \
		LDFLAGS="-L$$OPENSSL_CURRENT_LIB -lcrypto" \
		$(SQLCIPHER_DIR)/configure \
			--with-tclsh="$$SQLCIPHER_TCLSH" \
			--with-tcl="$$SQLCIPHER_TCL_CONFIG_DIR" \
			--disable-shared \
			--enable-static \
			--with-tempstore=yes

	$(MAKE) -C $(SQLCIPHER_CURRENT_DIR) libsqlite3.a sqlite3.h

	mkdir -p $(SQLCIPHER_CURRENT_INSTALL)/lib
	mkdir -p $(SQLCIPHER_CURRENT_INSTALL)/include

	cp $(SQLCIPHER_CURRENT_DIR)/libsqlite3.a \
		$(SQLCIPHER_CURRENT_INSTALL)/lib/

	cp $(SQLCIPHER_CURRENT_DIR)/sqlite3.h \
		$(SQLCIPHER_CURRENT_INSTALL)/include/

	cp $(SQLCIPHER_DIR)/src/sqlite3ext.h \
		$(SQLCIPHER_CURRENT_INSTALL)/include/


build-sqlcipher-386: $(SQLCIPHER_386_LIB)

$(SQLCIPHER_386_LIB): $(SQLCIPHER_DIR)/configure
	mkdir -p $(SQLCIPHER_386_DIR)
	cd $(SQLCIPHER_386_DIR) && \
		CC="$$CC_386" \
		CPPFLAGS="-I$$OPENSSL_386_INCLUDE" \
		CFLAGS="$(SQLCIPHER_CFLAGS)" \
		LDFLAGS="-L$$OPENSSL_386_LIB -lcrypto" \
		$(SQLCIPHER_DIR)/configure \
			--host="$$($$CC_386 -dumpmachine)" \
			--with-tclsh="$$SQLCIPHER_TCLSH" \
			--with-tcl="$$SQLCIPHER_TCL_CONFIG_DIR" \
			--disable-shared \
			--enable-static \
			--with-tempstore=yes

	$(MAKE) -C $(SQLCIPHER_386_DIR) libsqlite3.a sqlite3.h

	mkdir -p $(SQLCIPHER_386_INSTALL)/lib
	mkdir -p $(SQLCIPHER_386_INSTALL)/include

	cp $(SQLCIPHER_386_DIR)/libsqlite3.a \
		$(SQLCIPHER_386_INSTALL)/lib/

	cp $(SQLCIPHER_386_DIR)/sqlite3.h \
		$(SQLCIPHER_386_INSTALL)/include/

	cp $(SQLCIPHER_DIR)/src/sqlite3ext.h \
		$(SQLCIPHER_386_INSTALL)/include/


build-sqlcipher-arm7: $(SQLCIPHER_ARM7_LIB)

$(SQLCIPHER_ARM7_LIB): $(SQLCIPHER_DIR)/configure
	mkdir -p $(SQLCIPHER_ARM7_DIR)
	cd $(SQLCIPHER_ARM7_DIR) && \
		CC="$$CC_ARMV7" \
		CPPFLAGS="-I$$OPENSSL_ARMV7_INCLUDE" \
		CFLAGS="$(SQLCIPHER_CFLAGS)" \
		LDFLAGS="-L$$OPENSSL_ARMV7_LIB -lcrypto" \
		$(SQLCIPHER_DIR)/configure \
			--host="$$($$CC_ARMV7 -dumpmachine)" \
			--with-tclsh="$$SQLCIPHER_TCLSH" \
			--with-tcl="$$SQLCIPHER_TCL_CONFIG_DIR" \
			--disable-shared \
			--enable-static \
			--with-tempstore=yes

	$(MAKE) -C $(SQLCIPHER_ARM7_DIR) libsqlite3.a sqlite3.h

	mkdir -p $(SQLCIPHER_ARM7_INSTALL)/lib
	mkdir -p $(SQLCIPHER_ARM7_INSTALL)/include

	cp $(SQLCIPHER_ARM7_DIR)/libsqlite3.a \
		$(SQLCIPHER_ARM7_INSTALL)/lib/

	cp $(SQLCIPHER_ARM7_DIR)/sqlite3.h \
		$(SQLCIPHER_ARM7_INSTALL)/include/

	cp $(SQLCIPHER_DIR)/src/sqlite3ext.h \
		$(SQLCIPHER_ARM7_INSTALL)/include/


build-sqlcipher-arm64: $(SQLCIPHER_ARM64_LIB)

$(SQLCIPHER_ARM64_LIB): $(SQLCIPHER_DIR)/configure
	mkdir -p $(SQLCIPHER_ARM64_DIR)
	cd $(SQLCIPHER_ARM64_DIR) && \
		CC="$$CC_ARM64" \
		CPPFLAGS="-I$$OPENSSL_ARM64_INCLUDE" \
		CFLAGS="$(SQLCIPHER_CFLAGS)" \
		LDFLAGS="-L$$OPENSSL_ARM64_LIB -lcrypto" \
		$(SQLCIPHER_DIR)/configure \
			--host="$$($$CC_ARM64 -dumpmachine)" \
			--with-tclsh="$$SQLCIPHER_TCLSH" \
			--with-tcl="$$SQLCIPHER_TCL_CONFIG_DIR" \
			--disable-shared \
			--enable-static \
			--with-tempstore=yes

	$(MAKE) -C $(SQLCIPHER_ARM64_DIR) libsqlite3.a sqlite3.h

	mkdir -p $(SQLCIPHER_ARM64_INSTALL)/lib
	mkdir -p $(SQLCIPHER_ARM64_INSTALL)/include

	cp $(SQLCIPHER_ARM64_DIR)/libsqlite3.a \
		$(SQLCIPHER_ARM64_INSTALL)/lib/

	cp $(SQLCIPHER_ARM64_DIR)/sqlite3.h \
		$(SQLCIPHER_ARM64_INSTALL)/include/

	cp $(SQLCIPHER_DIR)/src/sqlite3ext.h \
		$(SQLCIPHER_ARM64_INSTALL)/include/


# ------------------------------------------------------------------------------
# Go shared library
# ------------------------------------------------------------------------------

generate-go-sqlite3:
	./scripts/generate-go-sqlite3-sqlcipher.sh


build-lib-current: build-sqlcipher-current generate-go-sqlite3
	CGO_ENABLED=1 \
	CC="gcc" \
	CGO_CFLAGS="-I$(SQLCIPHER_CURRENT_INSTALL)/include" \
	CGO_LDFLAGS="-L$(SQLCIPHER_CURRENT_INSTALL)/lib -Wl,-Bstatic -lsqlite3 -Wl,-Bdynamic -L$$OPENSSL_CURRENT_LIB -lcrypto -lm" \
	go build \
		-modfile "$(GO_SQLITE3_MODFILE)" \
		-tags "$(GO_BUILD_TAGS)" \
		-buildmode=c-shared \
		-o libfioclient.so \
		./cbindings


build-lib-386: build-sqlcipher-386 generate-go-sqlite3
	GOOS=linux \
	GOARCH=386 \
	CGO_ENABLED=1 \
	CC="$$CC_386" \
	CGO_CFLAGS="-I$(SQLCIPHER_386_INSTALL)/include" \
	CGO_LDFLAGS="-L$(SQLCIPHER_386_INSTALL)/lib -Wl,-Bstatic -lsqlite3 -Wl,-Bdynamic -L$$OPENSSL_386_LIB -lcrypto -lm" \
	go build \
		-modfile "$(GO_SQLITE3_MODFILE)" \
		-tags "$(GO_BUILD_TAGS)" \
		-buildmode=c-shared \
		-o libfioclient.so \
		./cbindings
	patchelf --remove-rpath libfioclient.so


build-lib-arm7: build-sqlcipher-arm7 generate-go-sqlite3
	GOOS=linux \
	GOARCH=arm \
	GOARM=7 \
	CGO_ENABLED=1 \
	CC="$$CC_ARMV7" \
	CGO_CFLAGS="-I$(SQLCIPHER_ARM7_INSTALL)/include" \
	CGO_LDFLAGS="-L$(SQLCIPHER_ARM7_INSTALL)/lib -Wl,-Bstatic -lsqlite3 -Wl,-Bdynamic -L$$OPENSSL_ARMV7_LIB -lcrypto -lm" \
	go build \
		-modfile "$(GO_SQLITE3_MODFILE)" \
		-tags "$(GO_BUILD_TAGS)" \
		-buildmode=c-shared \
		-o libfioclient.so \
		./cbindings
	patchelf --remove-rpath libfioclient.so


build-lib-arm64: build-sqlcipher-arm64 generate-go-sqlite3
	GOOS=linux \
	GOARCH=arm64 \
	CGO_ENABLED=1 \
	CC="$$CC_ARM64" \
	CGO_CFLAGS="-I$(SQLCIPHER_ARM64_INSTALL)/include" \
	CGO_LDFLAGS="-L$(SQLCIPHER_ARM64_INSTALL)/lib -Wl,-Bstatic -lsqlite3 -Wl,-Bdynamic -L$$OPENSSL_ARM64_LIB -lcrypto -lm" \
	go build \
		-modfile "$(GO_SQLITE3_MODFILE)" \
		-tags "$(GO_BUILD_TAGS)" \
		-buildmode=c-shared \
		-o libfioclient.so \
		./cbindings
	patchelf --remove-rpath libfioclient.so


# ------------------------------------------------------------------------------
# Release
# ------------------------------------------------------------------------------

release-lib-current: build-lib-current
	mkdir -p out/current-os
	mv libfioclient.so out/current-os/
	mv libfioclient.h out/current-os/


release-lib-386: build-lib-386
	mkdir -p out/386
	mv libfioclient.so out/386/
	mv libfioclient.h out/386/


release-lib-arm7: build-lib-arm7
	mkdir -p out/arm7
	mv libfioclient.so out/arm7/
	mv libfioclient.h out/arm7/


release-lib-arm64: build-lib-arm64
	mkdir -p out/arm64
	mv libfioclient.so out/arm64/
	mv libfioclient.h out/arm64/


release-all: \
	release-lib-current \
	release-lib-arm7 \
	release-lib-386 \
	release-lib-arm64


# ------------------------------------------------------------------------------
# Tests
# ------------------------------------------------------------------------------

test: test-c
	go test ./...


test-c: build-lib-current
	$(CC) -std=c11 -Wall -Wextra -Werror -Wno-unused-function \
		-I. cbindings/tests/api_test.c -L. -lfioclient \
		-o /tmp/fioclient-c-api-test
	LD_LIBRARY_PATH=. /tmp/fioclient-c-api-test


# ------------------------------------------------------------------------------
# Cleanup
# ------------------------------------------------------------------------------

clean-sqlcipher:
	rm -rf build/sqlcipher


clean:
	rm -rf build out
	rm -f libfioclient.so libfioclient.h


# ------------------------------------------------------------------------------
# Sailfish architecture aliases
# ------------------------------------------------------------------------------

build-lib-armv7hl: build-lib-arm7
build-lib-aarch64: build-lib-arm64
build-lib-i486: build-lib-386


.PHONY: \
	build-sqlcipher-current \
	build-sqlcipher-386 \
	build-sqlcipher-arm7 \
	build-sqlcipher-arm64 \
	generate-go-sqlite3 \
	build-lib-current \
	build-lib-386 \
	build-lib-arm7 \
	build-lib-arm64 \
	release-lib-current \
	release-lib-386 \
	release-lib-arm7 \
	release-lib-arm64 \
	release-all \
	test \
	test-c \
	clean-sqlcipher \
	clean \
	build-lib-armv7hl \
	build-lib-aarch64 \
	build-lib-i486
