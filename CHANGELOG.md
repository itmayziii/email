# Changelog

## [1.1.1](https://github.com/itmayziii/email/compare/v1.1.0...v1.1.1) (2023-10-06)


### Bug Fixes

* **no angle-addr error:** fixed issue where to, cc, and bcc were not … ([#10](https://github.com/itmayziii/email/issues/10)) ([6ba5251](https://github.com/itmayziii/email/commit/6ba525139557e6bd0cbcdef8159c0c8e3e03cae6))

## [1.1.0](https://github.com/itmayziii/email/compare/v1.0.0...v1.1.0) (2023-10-05)


### Features

* **cc and bcc:** added the ability to (blind) carbon copy by proving "cc" and "bcc" ([0d4e18a](https://github.com/itmayziii/email/commit/0d4e18a6d09cf71da5109aed22471275b9491d7a))

## 1.0.0 (2023-10-05)


### Features

* **domain parsing:** now parsing domain from sender + using domain to match appropriate sender ([dd67d9e](https://github.com/itmayziii/email/commit/dd67d9eac47b9d63079a03942217196ec10e450b))
* **email templates:** option to provide email body or specify a template at a file path ([3422816](https://github.com/itmayziii/email/commit/34228162ea51a9cdb330918a23f4791ae5f7545f))
* **google cloud build + docs:** inital docs added with Astro and starlight ([86d6705](https://github.com/itmayziii/email/commit/86d67051f5e8e82e637e6eea7c3279b6cf16bc68))
* **initial commit :rocket::** cloud function that sends email based on CloudEvent ([8b52157](https://github.com/itmayziii/email/commit/8b52157d75f90a547a26ccbbc5460ed9a438cf87))
* **interfaces:** using interfaces to make the main parts of this function swappable ([cc2d0ab](https://github.com/itmayziii/email/commit/cc2d0ab2127bfa6181333d6e25500124b758df03))
* **pub/sub:** auto extracting pub sub data when type matches pub sub ([a25719a](https://github.com/itmayziii/email/commit/a25719afe0fc1e200cfb33b6b524db03babf5bf8))
