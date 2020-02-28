/*
    Copyright SecureKey Technologies Inc. All Rights Reserved.

    SPDX-License-Identifier: Apache-2.0
*/

var {When} = require('cucumber');
var myStepDefinitionsWrapper = function () {
    this.When(/^client sends request to "([^"]*)" to create DID document "([^"]*)" in namespace "([^"]*)"$/, function (arg1, arg2, arg3, callback) {
        callback.pending();
    });
    this.When(/^client sends request to "([^"]*)" to resolve DID document with initial value$/, function (arg1, callback) {
        callback.pending();
    });
    this.Then(/^check success response contains "([^"]*)"$/, function (arg1, callback) {
        callback.pending();
    });
    this.Then(/^check error response contains "([^"]*)"$/, function (arg1, callback) {
        callback.pending();
    });
    this.When(/^client sends request to "([^"]*)" to resolve DID document$/, function (callback) {
        callback.pending();
    });
    this.When(/^client writes content "([^"]*)" using "([^"]*)" on the "([^"]*)" channel$/, function (arg1, arg2, arg3, callback) {
        callback.pending();
    });
    this.Then(/^client verifies that written content at the returned address from "([^"]*)" matches original content on the "([^"]*)" channel$/, function (arg1, arg2, callback) {
        callback.pending();
    });
    this.When(/^client creates document with ID "([^"]*)" using "([^"]*)" on the "([^"]*)" channel$/, function (arg1, arg2, arg3, callback) {
        callback.pending();
    });
    this.Then(/^client verifies that query by index ID "([^"]*)" from "([^"]*)" will return "([^"]*)" versions of the document on the "([^"]*)" channel$/, function (arg1, arg2, arg3, arg4, callback) {
        callback.pending();
    });
    this.Then(/^client verifies that query by index ID "([^"]*)" from "([^"]*)" will return "([^"]*)" versions of the document on the "([^"]*)" channel on peers "([^"]*)"$/, function (arg1, arg2, arg3, arg4, arg5, callback) {
        callback.pending();
    });
    this.When(/^client writes operations batch file and anchor file for ID "([^"]*)" using "([^"]*)" on the "([^"]*)" channel$/, function (arg1, arg2, arg3, callback) {
        callback.pending();
    });
    this.When(/^client sends request to "([^"]*)" to delete DID document$/, function (arg1, callback) {
        callback.pending();
    });
};
module.exports = myStepDefinitionsWrapper;
