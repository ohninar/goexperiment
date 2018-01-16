GO18='/usr/local/go181/go/bin/go'
GCCGO7='gccgo-7'

build-go18-randmat-par:
	$(GO18) build -o target/randmat/expertpar/maingo18bin target/randmat/expertpar/main.go

build-go18-randmat-seq:
	$(GO18) build -o target/randmat/expertseq/maingo18bin target/randmat/expertseq/main.go

build-go18-outer-par:
	$(GO18) build -o target/outer/expertpar/maingo18bin target/outer/expertpar/main.go

build-go18-outer-seq:
	$(GO18) build -o target/outer/expertseq/maingo18bin target/outer/expertseq/main.go

build-go18-product-par:
	$(GO18) build -o target/product/expertpar/maingo18bin target/product/expertpar/main.go

build-go18-product-seq:
	$(GO18) build -o target/product/expertseq/maingo18bin target/product/expertseq/main.go

build-go18-thresh-par:
	$(GO18) build -o target/thresh/expertpar/maingo18bin target/thresh/expertpar/main.go

build-go18-thresh-seq:
	$(GO18) build -o target/thresh/expertseq/maingo18bin target/thresh/expertseq/main.go

build-go18-winnow-par:
	$(GO18) build -o target/winnow/expertpar/maingo18bin target/winnow/expertpar/main.go

build-go18-winnow-seq:
	$(GO18) build -o target/winnow/expertseq/maingo18bin target/winnow/expertseq/main.go

build-go18:
    build-go18-randmat-par
    build-go18-randmat-seq
    build-go18-outer-par
    build-go18-outer-seq
    build-go18-product-par
    build-go18-product-seq
    build-go18-winnow-par
    build-go18-winnow-seq

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

build-gccgo7:
    build-gccgo7-randmat-par
    build-gccgo7-randmat-seq
    build-gccgo7-outer-par
    build-gccgo7-outer-seq
    build-gccgo7-product-par
    build-gccgo7-product-seq
    build-gccgo7-winnow-par
    build-gccgo7-winnow-seq

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

build-gccgo7opt:
    build-gccgo7opt-randmat-par
    build-gccgo7opt-randmat-seq
    build-gccgo7opt-outer-par
    build-gccgo7opt-outer-seq
    build-gccgo7opt-product-par
    build-gccgo7opt-product-seq
    build-gccgo7opt-winnow-par
    build-gccgo7opt-winnow-seq