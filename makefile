md:
	brew tap ahuigo/homebrew-arun https://github.com/ahuigo/homebrew-arun.git
	cd $(brew --repo ahuigo/tap)
build:
	brew create https://github.com/ahuigo/arun/archive/refs/tags/master.tar.gz --tap ahuigo/homebrew-arun
	
install:
	brew install ahuigo/tap/arun
	
