{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  name = "dev-environment";
  buildInputs = [
    pkgs.minikube
    pkgs.kubectl
  ];
}
