PKGS=$(shell go list ./...)
LAMBDAS = \
	cmd/extract \
	cmd/transform \
	cmd/load

clean:
	git clean -fx

test:
	go vet $(PKGS)
	go test -v -coverprofile coverage.out $(PKGS)
	go tool cover -func=coverage.out

build:
	for lambda in $(LAMBDAS); do \
		GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o $$lambda/main ./$$lambda; \
	done

zip:
	go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
	for lambda in $(LAMBDAS); do \
		build-lambda-zip -o $$lambda/main.zip $$lambda/main; \
	done

deploy: zip
	aws s3 mb s3://${DEPLOYMENT_NAME} || true
	aws s3 rm --recursive s3://${IMPORT_SOURCE_BUCKET_NAME} || true
	aws s3 rm --recursive s3://${IMPORT_RESULT_BUCKET_NAME} || true
	sam package --template-file template.yml --output-template-file packaged.yml --s3-bucket ${DEPLOYMENT_NAME}
	sam deploy --template-file packaged.yml --stack-name ${DEPLOYMENT_NAME} \
		--capabilities CAPABILITY_IAM --parameter-overrides \
		SourceBucketName=${IMPORT_SOURCE_BUCKET_NAME} \
		ResultBucketName=${IMPORT_RESULT_BUCKET_NAME} \
		LoadUrl=${LOAD_URL} LoadUser=${LOAD_USER} LoadPswd=${LOAD_PSWD} \
		|| true
