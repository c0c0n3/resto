#
# Tweaked from the original to upgrade version to `2022.1.1` and build
# w/ Go `1.18`:
#
# - https://github.com/NixOS/nixpkgs/blob/nixos-unstable/pkgs/development/tools/go-tools/
#
# Reason: The current Nixpkgs version `2021.1.2` doesn't work with VS
# Code and Go `1.8`. VS Code moans about:
#
#   Tools (staticcheck) need recompiling to work with
#   go version go1.18.2 darwin/amd64
#
# NOTE. SHA256 hashes.
# To get the right values, I first commented out `src.sha256` and then
# built. The Nix error message contained the right hash to use. For the
# vendor hash, I followed this procedure:
# - https://stackoverflow.com/questions/71927395
#
# TODO get rid of this package as soon as `2022.1.1` built w/ Go `1.18`
# gets released on Nixpkgs.
#
{ buildGo118Module
, lib
, fetchFromGitHub
}:

buildGo118Module rec {
  pname = "go-tools";
  version = "2022.1.1";

  src = fetchFromGitHub {
    owner = "dominikh";
    repo = "go-tools";
    rev = version;
    sha256 = "sha256-SWbpn3IKSCuVnu2bHKkcLwLyqx4P9vkCZG17Fxq6cA4=";
  };

  vendorSha256 = "sha256-aOtNjWHQUN2iD26PvJEKpOCog72L9mXFXcsJiusGm20=";

  subPackages = [
    "cmd/keyify"
    "cmd/staticcheck"
    "cmd/structlayout-optimize"
    "cmd/structlayout-pretty"
    "cmd/structlayout"
  ];

  doCheck = false;

  meta = with lib; {
    description = "A collection of tools and libraries for working with Go code, including linters and static analysis";
    homepage = "https://staticcheck.io";
    license = licenses.mit;
    maintainers = with maintainers; [ rvolosatovs kalbasit ];
  };
}
