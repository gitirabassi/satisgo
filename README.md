## Satisgo SDK

This is an SDK for the Satispay Business API.

It's not production ready but all the functions are there.

Obviously since every call has go to the the Satispay endpoint, making an https request, this cannot be done syncronously with you page loading, but it can be done asyncrounous and nothing bad can happen.

## Roadmap

- [ ] Strengthen security with `http.Transport` & `tls.Config` structs. The fundation work has been done, PR welcomed
- [ ] Shorten methods name

## Installation

`go get github.com/drymonsoon/satisgo`

## Usage

- examples are coming

## Documentation

https://s3-eu-west-1.amazonaws.com/docs.online.satispay.com/index.html

## Warning

This library is licenced MIT: everyone can use, improve it, change it and make a lot of money on it.
You cannot sue me if something goes wrong, you take all the responsability using this.
