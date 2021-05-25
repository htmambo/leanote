# http://revel.github.io/manual/tool.html

# go get -u github.com/revel/cmd/revel
SCRIPTPATH=$(dirname "$PWD")
echo $SCRIPTPATH;
cd $SCRIPTPATH;

/Volumes/Workarea/golang/bin/revel run -a .
