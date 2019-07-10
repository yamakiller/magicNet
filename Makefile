BUILD_VERSION   := v1.0.0
BUILD_TIME      := $(shell date "+%F %T")
BUILD_NAME      := app_$(shell date "+%Y%m%d%H" )
SOURCE          := ./src/main/main.go
TARGET_DIR      := ./build/
COMMIT_SHA1     := $(shell git rev-parse HEAD ) 
# svn info last rev
#$(shell svn info |grep "Last Changed Rev: " |sed -e "s/Last Changed Rev: "//g)
#$(shell svn info |grep Revision|awk '{print $$2}')

all:
    # CGO_ENABLED=0 GOOS=linux GOARCH=amd64
    go build -ldflags                           \
    "                                           \
    -X 'version.BuildVersion=${BUILD_VERSION}'  \
    -X 'version.BuildTime=${BUILD_TIME}'        \
    -X 'version.BuildName=${BUILD_NAME}'        \
    -X 'version.CommitID=${COMMIT_SHA1}'        \
    "                                           \
    -o ${BUILD_NAME} ${SOURCE}

clean:
	rm ${BUILD_NAME} -f 

install:
	mkdir -p ${TARGET_DIR}
	cp ${BUILD_NAME} ${TARGET_DIR} -f

.PHONY : all clean install ${BUILD_NAME}