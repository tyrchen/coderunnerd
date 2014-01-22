CHECK=\033[32m✔\033[39m
DONE="\n$(CHECK) Done.\n"

SERVER=jcnrd.us
PROJECT=coderunnerd
PATH=deployment/$(PROJECT)
SUPERVISORCTL=/usr/bin/supervisorctl
SUCOPY=/bin/sucopy
SSH=/usr/bin/ssh
ECHO=/bin/echo -e
RM=/bin/rm
SUDO=/usr/bin/sudo
GO=$(shell which go)
CLOC=$(shell which cloc)
BIN=./bin
TARGETS=$(patsubst %.go,$(BIN)/%,$(wildcard *.go))

build: $(TARGETS)
	@echo $(DONE)

$(TARGETS): $(BIN)/%: %.go
	@echo "building $<..."
	@$(GO) build -o $@ $<

remote_deploy:
	@$(SSH) -t $(SERVER) "echo Deploy $(PROJECT) to the $(SERVER) server.; cd $(PATH);  git pull; make deploy;"

dependency:
	@$(ECHO) "\nInstall project dependencies..."

configuration:
	@$(ECHO) "\nUpdate configuration..."
	@$(SUDO) $(SUCOPY) -r _deploy/etc/. /etc/.

supervisor:
	@$(ECHO) "\nUpdate supervisor configuration..."
	@$(SUDO) $(SUPERVISORCTL) reread
	@$(SUDO) $(SUPERVISORCTL) update
	@$(ECHO) "\nRestart $(PROJECT)..."
	@$(SUDO) $(SUPERVISORCTL) restart $(PROJECT)

deploy: clean dependency configuration supervisor
	@$(ECHO) $(DONE)

clean:
	@$(RM) -f $(TARGETS)
	@echo $(DONE)

cloc:
	@$(CLOC) . --exclude-dir=webclient/assets
