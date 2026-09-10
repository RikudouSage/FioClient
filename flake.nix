{
  description = "Dev shell for CGO cross-compilation (386, armv7, arm64)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nixpkgs-legacy.url = "github:NixOS/nixpkgs/5d61a47f31092319ab82b44213420974b91c9bc1";
  };

  outputs = { self, nixpkgs, nixpkgs-legacy }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      forAllSystems = f:
        builtins.listToAttrs (map (system: {
          name = system;
          value = f system;
        }) systems);
    in {
      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
          legacyPkgs = import nixpkgs-legacy { inherit system; };

          cross386   = legacyPkgs.pkgsCross.gnu32;
          crossArmv7 = legacyPkgs.pkgsCross.armv7l-hf-multiplatform;
          crossArm64 = legacyPkgs.pkgsCross.aarch64-multiplatform;

          cc386   = cross386.stdenv.cc;
          ccArmv7 = crossArmv7.stdenv.cc;
          ccArm64 = crossArm64.stdenv.cc;
        in {
          default = pkgs.mkShell {
            buildInputs = [
              pkgs.bash
              pkgs.go_1_27
              pkgs.golangci-lint
              pkgs.patchelf
              pkgs.revive
              pkgs.gnumake

              pkgs.gcc
              pkgs.pkg-config
              pkgs.openssl
              pkgs.tcl

              cc386
              ccArmv7
              ccArm64

              cross386.openssl
              crossArmv7.openssl
              crossArm64.openssl
            ];

            shellHook = ''
              unset GOROOT

              export CC_386=${cc386.targetPrefix}cc
              export CC_ARMV7=${ccArmv7.targetPrefix}cc
              export CC_ARM64=${ccArm64.targetPrefix}cc

              export OPENSSL_386_INCLUDE=${cross386.openssl.dev}/include
              export OPENSSL_386_LIB=${cross386.openssl.out}/lib

              export OPENSSL_ARMV7_INCLUDE=${crossArmv7.openssl.dev}/include
              export OPENSSL_ARMV7_LIB=${crossArmv7.openssl.out}/lib

              export OPENSSL_ARM64_INCLUDE=${crossArm64.openssl.dev}/include
              export OPENSSL_ARM64_LIB=${crossArm64.openssl.out}/lib

              export OPENSSL_CURRENT_INCLUDE=${pkgs.openssl.dev}/include
              export OPENSSL_CURRENT_LIB=${pkgs.openssl.out}/lib

              export SQLCIPHER_TCLSH=${pkgs.tcl}/bin/tclsh
              export SQLCIPHER_TCL_CONFIG_DIR=${pkgs.tcl}/lib

              echo "Using cross compilers:"
              echo "  CC_386   = $CC_386"
              echo "  CC_ARMV7 = $CC_ARMV7"
              echo "  CC_ARM64 = $CC_ARM64"

              echo ""

              echo "Using go config:"
              echo "  GOROOT   = $(go env GOROOT)"
              echo "  GOCACHE  = $(go env GOCACHE)"
              echo "  GOPATH   = $(go env GOPATH)"
            '';
          };
        }
      );
    };
}
