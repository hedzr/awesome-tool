# Documentation: https://docs.brew.sh/Formula-Cookbook
#                https://rubydoc.brew.sh/Formula
# PLEASE REMOVE ALL GENERATED COMMENTS BEFORE SUBMITTING YOUR PULL REQUEST!
class AwesomeTool < Formula
  desc "a command-line tool to retrieve the stars of all repos in an awesome-list"
  homepage ""
  url "https://github.com/hedzr/awesome-tool/releases/download/v1.1.11/awesome-tool-v1.1.11-darwin-amd64.tgz"
  sha256 "027594a8e005ea6b893a544a18e5c98c8c152e80c52a17182d00a20057d99ad7"
  license "MIT"

  depends_on "go" => :build

  def install
    # ENV.deparallelize  # if your formula fails when building in parallel
    
    # ENV["GOPROXY"] = "https://goproxy.io"
    # system "make"

    system "go", "build", *std_go_args
  end

  test do
    # `test do` will create, run in and delete a temporary directory.
    #
    # This test will fail and we won't accept that! For Homebrew/homebrew-core
    # this will need to be a test that verifies the functionality of the
    # software. Run the test with `brew test awesome-tool`. Options passed
    # to `brew install` such as `--HEAD` also need to be provided to `brew test`.
    #
    # The installed folder is not in the path, so use the entire path to any
    # executables being tested: `system "#{bin}/program", "do", "something"`.
    system "false"
  end
end
