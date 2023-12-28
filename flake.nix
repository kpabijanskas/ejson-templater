{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }: 
    flake-utils.lib.eachDefaultSystem (system:
      let
        name = "ejson-templater";
        version = "1.0.1";
        pkgs = nixpkgs.legacyPackages."${system}";
      in {
        packages.default = pkgs.buildGoModule {
          name = "ejson-templater";

          src = ./.;

          vendorHash = "sha256-p5/w5uiUQ00IImQJjxClxnMX6yljhdTKoxJqA7lqeK0=";

          postInstall = ''
            mv $out/bin/templater $out/bin/${name}
          '';

          meta = with pkgs.lib; {
            description = "ejson-templater";
            homepage = "https://github.com/kpabijanskas/ejson-templater";
            license = licenses.mit;
          };
        };
      }
    );
}
