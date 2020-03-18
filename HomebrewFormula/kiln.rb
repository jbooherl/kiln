# This file was generated by GoReleaser. DO NOT EDIT.
class Kiln < Formula
  desc ""
  homepage ""
  version "0.41.0"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/pivotal-cf/kiln/releases/download/0.41.0/kiln-darwin-0.41.0.tar.gz"
    sha256 "c767ecd099f46c75aefb9dfd738fb45c21832dfe859c4e4a8ee136ae02192def"
  elsif OS.linux?
    if Hardware::CPU.intel?
      url "https://github.com/pivotal-cf/kiln/releases/download/0.41.0/kiln-linux-0.41.0.tar.gz"
      sha256 "87cc9e5e7e7ada47dc21d730eea36f62db66dc684956ac4c7aea3f3185d3bed0"
    end
  end

  def install
    bin.install "kiln"
  end

  test do
    system "#{bin}/kiln --version"
  end
end