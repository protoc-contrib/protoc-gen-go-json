{
  description = "protoc-gen-go-json - A protoc plugin that generates MarshalJSON and UnmarshalJSON methods backed by protojson";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    nix-release-bin = {
      url = "github:nixos-contrib/nix-release-bin";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      nix-release-bin,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version = (pkgs.lib.importJSON ./.github/config/release-please-manifest.json).".";

        source = pkgs.buildGoModule {
          pname = "protoc-gen-go-json";
          inherit version;
          src = pkgs.lib.cleanSource ./.;
          subPackages = [ "cmd/protoc-gen-go-json" ];
          vendorHash = "sha256-KLrEjd4dBG3zNrz2wTP4dBbGE5kImdCkv+NG7wD3WxM=";
          ldflags = [ "-s" "-w" ];
          meta = with pkgs.lib; {
            description = "A protoc plugin that generates MarshalJSON and UnmarshalJSON methods backed by protojson";
            license = licenses.asl20;
            mainProgram = "protoc-gen-go-json";
          };
        };
      in
      {
        packages = {
          # The latest release binary, where it has one for the system: CI pins
          # them in the manifest once the release has published them.
          default = nix-release-bin.lib.mkReleaseBin {
            inherit pkgs;
            manifest = ./.github/config/nix-release-bin-manifest.json;
            pname = "protoc-gen-go-json";
            fallback = source;
          };
          inherit source;
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
