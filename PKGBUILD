# Maintainer: Your Name <your.email@example.com>
pkgname=godotctl-git
pkgver=r7.baa0be3
pkgrel=1
pkgdesc="A fast, interactive dotfiles installer with symlink management and automatic backups (git version)"
arch=('x86_64' 'aarch64' 'armv7h')
url="https://github.com/grainedlotus515/godotctl"
license=('MIT')
depends=('git')
makedepends=('go')
provides=('godotctl')
conflicts=('godotctl')
source=("${pkgname}::git+${url}.git")
sha256sums=('SKIP')

pkgver() {
    cd "${srcdir}/${pkgname}"
    printf "r%s.%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short HEAD)"
}

# Ensure Go builds reproducible binaries
export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"

build() {
    cd "${srcdir}/${pkgname}"

    go build \
        -ldflags="-s -w -X main.version=${pkgver}" \
        -o godotctl \
        ./cmd
}

check() {
    cd "${srcdir}/${pkgname}"
    go test ./...
}

package() {
    cd "${srcdir}/${pkgname}"

    # Install binary
    install -Dm755 godotctl "${pkgdir}/usr/bin/godotctl"

    # Install license
    install -Dm644 LICENSE "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"

    # Install documentation
    install -Dm644 README.md "${pkgdir}/usr/share/doc/${pkgname}/README.md"
}
