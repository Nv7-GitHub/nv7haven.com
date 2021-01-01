THING=$(cat docs/index.html | pcregrep -o1 '\?v=([0-9]+)')
THING=$((THING + 1))
export THING
sed -i '' "s/\?v=1/?v=${THING}/g" docs/index.html
