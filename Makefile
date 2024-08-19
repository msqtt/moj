IMAGE_TAG ?= dev

test:
	go test -v moj/...
	
img-web-bff:
	docker build . -t msqt/moj-web-bff:$(IMAGE_TAG) -f builds/web-bff.Dockerfile 
img-user:
	docker build . -t msqt/moj-user:$(IMAGE_TAG) -f builds/user.Dockerfile 
img-question:
	docker build . -t msqt/moj-question:$(IMAGE_TAG) -f builds/question.Dockerfile 
img-game:
	docker build . -t msqt/moj-game:$(IMAGE_TAG) -f builds/game.Dockerfile 
img-record:
	docker build . -t msqt/moj-record:$(IMAGE_TAG) -f builds/record.Dockerfile 
img-judgement:
	docker build . -t msqt/moj-judgement:$(IMAGE_TAG) -f builds/judgement.Dockerfile 

.PHONY: test img-web-bff img-user img-question img-game img-record img-judgement
