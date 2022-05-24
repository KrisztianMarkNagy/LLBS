.PHONY : install clean reinstall uninstall

PROGRAM_NAME = llbs
BINDIR = /usr/bin/

install:
	@echo "Installing..."
	go build .
	sudo mv ./$(PROGRAM_NAME) $(BINDIR)
	@echo "Done!"

clean:
	@echo "Cleaning..."
	rm -f ./$(PROGRAM_NAME)
	sudo rm $(BINDIR)$(PROGRAM_NAME)
# rm go.sum
	@echo "Done!"

reinstall: clean
	@echo "Tidying up the go module..."
	go mod tidy
	@echo "Done!"
#
	@echo "Installing..."
	go build .
	sudo mv ./$(PROGRAM_NAME) $(BINDIR)
	@echo "Done!"

uninstall: clean
	@echo "Removing executable from '/usr/bin/'..."
	sudo rm $(BINDIR)$(PROGRAM_NAME)
	@echo "Done!"

