{ pkgs }:
with pkgs;
let
  go-tools_1_18 = callPackage ./go-tools.nix {};
in buildEnv {
  name = "resto-dev-env";
  paths = [
    # Go dev env:
    # - Compiler
    # - VS Code tools (https://github.com/golang/vscode-go#tools)
    # - Nix build tool (https://www.tweag.io/blog/2021-03-04-gomod2nix/)
    go_1_18 gopls delve go-tools_1_18 go-outline
    # NOTE. delve includes dlv-dap; go-tools includes staticcheck
    gomod2nix
  ];
}
