commands := commands/backup commands/backup-key-agent
tests := lib/crypto/identity

all:
	@for d in $(commands); do \
	  echo $$d; \
	  (cd $$d; go install) \
	done

test:
	@for d in $(tests); do \
	  echo $$d; \
	  (cd $$d; go test) \
	done
