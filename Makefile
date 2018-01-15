GO18='/usr/local/go181/go/bin/go'
GCCGO7='gccgo-7'

build-go18-randmat-par:
	$(GO18) build -o target/randmat/expertpar/maingo18bin target/randmat/expertpar/main.go

build-go18-randmat-seq:
	$(GO18) build -o target/randmat/expertseq/maingo18bin target/randmat/expertseq/main.go

build-go18-outer-par:
	$GO18 target/outer/expertpar/main.go -o maingo18bin

build-go18-outer-seq:
	$GO18 target/outer/expertseq/main.go -o maingo18bin

build-go18-product-par:
	$GO18 target/product/expertpar/main.go -o maingo18bin

build-go18-product-seq:
	$GO18 target/product/expertseq/main.go -o maingo18bin

build-go18-thresh-par:
	$GO18 target/thresh/expertpar/main.go -o maingo18bin

build-go18-thresh-seq:
	$GO18 target/thresh/expertseq/main.go -o maingo18bin

build-go18-winnow-par:
	$GO18 target/winnow/expertpar/main.go -o maingo18bin

build-go18-winnow-seq:
	$GO18 target/winnow/expertseq/main.go -o maingo18bin


build-gccgo7-randmat-par:
	$GCCGO7 target/randmat/expertpar/main.go -o maingccgo7bin

build-gccgo7-randmat-seq:
	$GCCGO7 target/randmat/expertseq/main.go -o maingccgo7bin

build-gccgo7-outer-par:
	$GCCGO7 target/outer/expertpar/main.go -o maingccgo7bin

build-gccgo7-outer-seq:
	$GCCGO7 target/outer/expertseq/main.go -o maingccgo7bin

build-gccgo7-product-par:
	$GCCGO7 target/product/expertpar/main.go -o maingccgo7bin

build-gccgo7-product-seq:
	$GCCGO7 target/product/expertseq/main.go -o maingccgo7bin

build-gccgo7-thresh-par:
	$GCCGO7 target/thresh/expertpar/main.go -o maingccgo7bin

build-gccgo7-thresh-seq:
	$GCCGO7 target/thresh/expertseq/main.go -o maingccgo7bin

build-gccgo7-winnow-par:
	$GCCGO7 target/winnow/expertpar/main.go -o maingccgo7bin

build-gccgo7-winnow-seq:
	$GCCGO7 target/winnow/expertseq/main.go -o maingccgo7bin


build-gccgo7opt-randmat-par:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/randmat/expertpar/main.go -o maingccgo7optbin

build-gccgo7opt-randmat-seq:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/randmat/expertseq/main.go -o maingccgo7optbin

build-gccgo7opt-outer-par:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/outer/expertpar/main.go -o maingccgo7optbin

build-gccgo7opt-outer-seq:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/outer/expertseq/main.go -o maingccgo7optbin

build-gccgo7opt-product-par:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/product/expertpar/main.go -o maingccgo7optbin

build-gccgo7opt-product-seq:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/product/expertseq/main.go -o maingccgo7optbin

build-gccgo7opt-thresh-par:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/thresh/expertpar/main.go -o maingccgo7optbin

build-gccgo7opt-thresh-seq:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/thresh/expertseq/main.go -o maingccgo7optbin

build-gccgo7opt-winnow-par:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/winnow/expertpar/main.go -o maingccgo7optbin

build-gccgo7opt-winnow-seq:
	$GCCGO7 -O2 -O3 -fgo-optimize-allocs target/winnow/expertseq/main.go -o maingccgo7optbin
