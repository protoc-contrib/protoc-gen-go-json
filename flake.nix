{
  description = "protoc-gen-go-json - A protoc plugin that generates MarshalJSON and UnmarshalJSON methods backed by protojson";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version = (pkgs.lib.importJSON ./.github/config/release-please-manifest.json).".";
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "protoc-gen-go-json";
          inherit version;
          src = pkgs.lib.cleanSource ./.;
          subPackages = [ "cmd/protoc-gen-go-json" ];
          vendorHash = null;
          ldflags = [ "-s" "-w" ];
          meta = with pkgs.lib; {
            description = "A protoc plugin that generates MarshalJSON and UnmarshalJSON methods backed by protojson";
            license = licenses.asl20;
            mainProgram = "protoc-gen-go-json";
          };
        };

        devShells.default = pkgs.mkShell {
          name = "protoc-gen-go-json";
          packages = [
            pkgs.go
            pkgs.protobuf
            pkgs.buf
          ];
        };
      }
    );
}
