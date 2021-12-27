.PHONY: clean plugin release tex

xpl := ATR\ 72-500/plugins/default/lin_x64/default.xpl
xplwin := ATR\ 72-500/plugins/default/win_x64/default.xpl
xplmac := ATR\ 72-500/plugins/default/mac_x64/default.xpl

release: 
	cp COPYING ATR\ 72-500/
	cp README-release.txt ATR\ 72-500/
	7z a release/atr72_500-`git describe`.7z \
		ATR\ 72-500/*acf \
		ATR\ 72-500/cockpit_3d/-PANELS-/Panel_Airliner.png \
		ATR\ 72-500/objects/* \
		ATR\ 72-500/plugins/* \
		ATR\ 72-500/airfoils/* 

tex:
	inkscape svg/cockpit.svg -i std -j -C --export-png="ATR 72-500/objects/cockpit.png"
	inkscape svg/cockpit.svg -i LIT -j -C --export-png="ATR 72-500/objects/cockpit_LIT.png"
	inkscape svg/switches.svg -i std -j -C --export-png="ATR 72-500/objects/switches.png"
	inkscape svg/switches.svg -i LIT -j -C --export-png="ATR 72-500/objects/switches_LIT.png"

panel:
	inkscape svg/panel.svg -j -C --export-png="ATR 72-500/cockpit_3d/-PANELS-/Panel_Airliner.png"
	inkscape svg/panel_gen_trim_elev.svg -i gen_pointer -j -C --export-png="ATR 72-500/cockpit_3d/generic/trim_elev/gen_pointer.png"
	inkscape svg/panel_gen_trim_elev.svg -i gen_pointer-1 -j -C --export-png="ATR 72-500/cockpit_3d/generic/trim_elev/gen_pointer-1.png"

plugin: $(xplwin)

$(xplwin): plugin/*.go plugin/xpl/*.go
	GOOS=linux \
	GOARCH=amd64 \
	CGO_ENABLED=1 \
	CGO_CFLAGS="-I./plugin/xpl/XPLM -DLIN=1" \
	go build -buildmode c-shared -o $(xpl) plugin/xpl/main.go 
	CGO_CFLAGS="-DIBM=1 -static" \
	CGO_LDFLAGS="-L/home/mh/src/github.com/maxhille/atr-xplane/plugin/xpl/XPLM -lXPLM_64 -static-libgcc -static-libstdc++ -Wl,--exclude-libs,ALL" \
	GOOS=windows \
	GOARCH=amd64 \
	CGO_ENABLED=1 \
	CC=x86_64-w64-mingw32-gcc \
	CXX=x86_64-w64-mingw32-g++ \
	go build -x -buildmode c-shared -o $(xplwin) plugin/xpl/main.go \
    GOOS=darwin \
    GOARCH=amd64 \
    CGO_ENABLED=1 \
    CGO_CFLAGS="-I/Users/dzou/Downloads/SDK/CHeaders -DAPL=1 -DIBM=0 -DLIN=0 -DXPLM210=1 " \
    CGO_LDFLAGS="-F/System/Library/Frameworks/ -F/Users/dzou/Downloads/SDK/Libraries/Mac -framework XPLM" \
    go build -buildmode c-shared -o $(xplmac) plugin/xpl/main.go
clean:
	rm -r ATR\ 72-500/plugins
